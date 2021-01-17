package usecase

import (
	"bytes"
	"context"
	"errors"
	"github.com/saeed-kamyabi/dns-proxy/proxy/domain"
	"log"
	"net"
	"sync"
	"time"
)

var errHttpServerClosed = errors.New("http: Server closed")

type httpProxy struct {
	dns      domain.Dns
	mu       sync.Mutex
	doneChan chan struct{}
}

func newHttpServer(dns domain.Dns) *httpProxy {
	s := &httpProxy{
		dns: dns,
	}
	return s
}

func (s *httpProxy) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getDoneChanLocked()
}

func (s *httpProxy) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func (s *httpProxy) Serve(serverConfig domain.ServerConfig, l net.Listener) error {
	baseCtx := context.Background()

	var tempDelay time.Duration // how long to sleep on accept failure

	ctx := context.WithValue(baseCtx, "http-server", s)
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
				log.Printf("http: Accept error: %v; retrying in %v\n",
					err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		connCtx := ctx
		tempDelay = 0
		c := s.newConn(serverConfig, rw)
		go c.serve(connCtx)
	}
}

// Create new connection from rwc.
func (s *httpProxy) newConn(serverConfig domain.ServerConfig, rwc net.Conn) *httpConn {
	c := &httpConn{
		conn: conn{
			serverConfig: serverConfig,
			dns:          s.dns,
			rwc:          rwc,
		},
	}
	return c
}

type httpConn struct {
	conn
}

func (c *httpConn) serve(ctx context.Context) {
	var data []byte
	hostFound := false
	hostStartInd := 0
	// host index, index for check in data
	hi := 0
	var host string

	c.handler = func(n int, buf []byte) {
		data = append(data, buf[:n]...)

		if !hostFound {
			if hostStartInd == 0 {
				if i := bytes.Index(data[hi:], []byte("Host: ")); i >= 0 {
					hostStartInd = i + 6
				} else {
					hi = len(data) - 7
				}
			}
			if hostStartInd > 0 {
				for i := hostStartInd; i < len(data); i++ {
					if data[i] == '\n' {
						host = string(data[hostStartInd:i])
						hostFound = true
						go c.beginSendRequest(host)

						c.handler = nil
						return
					}
				}
			}
		}
	}
	c.conn.serve(ctx)
}
