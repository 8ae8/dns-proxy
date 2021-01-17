package usecase

import (
	"bytes"
	"context"
	"github.com/saeed-kamyabi/dns-proxy/proxy/domain"
	"github.com/saeed-kamyabi/dns-proxy/proxy/usecase/tls"
	"log"
	"net"
	"sync"
	"time"
)

type httpsProxy struct {
	dns      domain.Dns
	mu       sync.Mutex
	doneChan chan struct{}
}

func newHttpsServer(dns domain.Dns) *httpsProxy {
	s := &httpsProxy{
		dns: dns,
	}
	return s
}

func (s *httpsProxy) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getDoneChanLocked()
}

func (s *httpsProxy) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func (s *httpsProxy) Serve(serverConfig domain.ServerConfig, l net.Listener) error {
	baseCtx := context.Background()

	//l = tls.NewListener(l, nil)

	var tempDelay time.Duration // how long to sleep on accept failure

	ctx := context.WithValue(baseCtx, "https-server", s)
	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-s.getDoneChan():
				return errHttpServerClosed
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("https: Accept error: %v; retrying in %v\n",
					err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		httpsConnCtx := ctx
		tempDelay = 0
		c := s.newConn(serverConfig, rw)
		go c.serve(httpsConnCtx)
	}
}

// Create new httpsConnection from rwc.
func (s *httpsProxy) newConn(serverConfig domain.ServerConfig, rwc net.Conn) *httpsConn {
	c := &httpsConn{
		conn: conn{
			serverConfig: serverConfig,
			dns:          s.dns,
			rwc:          rwc,
		},
	}
	return c
}

type httpsConn struct {
	conn
}

// Serve a new httpsConnection.
func (c *httpsConn) serve(ctx context.Context) {
	c.handler = func(n int, buf []byte) {
		data := buf[5:n]
		hand := *bytes.NewBuffer(data[:])
		l := int(data[1])<<16 | int(data[2])<<8 | int(data[3])
		if l > tls.MaxHandshake {
			c.rwc.Close()
			return
		}
		data = hand.Next(4 + l)
		if data[0] == 1 {
			// client hello
			m := new(tls.ClientHelloMsg)
			if m.Unmarshal(data) {
				c.handler = nil
				go c.beginSendRequest(m.ServerName)
				return
			} else {
				c.rwc.Close()
				return
			}
		} else {
			c.rwc.Close()
			return
		}
	}
	c.conn.serve(ctx)
}
