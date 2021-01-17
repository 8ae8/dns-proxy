package usecase

import (
	"context"
	"github.com/saeed-kamyabi/dns-proxy/proxy/domain"
	"log"
	"net"
	"strings"
	"time"
)

type conn struct {
	serverConfig domain.ServerConfig
	dns          domain.Dns
	// rwc is the underlying network connection.
	// This is never wrapped by other types and is the value given out
	// to CloseNotifier callers. It is usually of type *net.TCPConn or
	// *tls.Conn.
	rwc net.Conn

	handler func(n int, buf []byte)
	//
	input chan []byte
}

func (c *conn) beginSendRequest(host string) {
	host = strings.Trim(host, "\r")
	ip, err := c.dns.LookupHost(host + ".")
	if err != nil {
		log.Printf("error on LookupHost, err: %v\n", err)
		return
	}
	if ip == "" {
		log.Printf("no IP were fount for host %s", host)
		return
	}
	conn, err := net.Dial("tcp", ip+":"+c.serverConfig.GetPort())
	if err != nil {
		log.Printf("error on beginSendRequest, err: %v\n", err)
		return
	}
	go func() {
		defer func() {
			conn.Close()
			c.rwc.Close()
		}()
		buf := make([]byte, 1024)
		for {
			err = conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			if err != nil {
				log.Printf("error on conn.SetReadDeadline, err: %v\n", err)
				return
			}
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			if n == 0 {
				return
			}
			c.rwc.Write(buf[:n])
		}
	}()
	for {
		select {
		case data := <-c.input:
			conn.Write(data)
			break
		}
	}
}

// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	c.input = make(chan []byte, 10)

	// read request
	buf := make([]byte, 1024)

	for {
		n, err := c.rwc.Read(buf)
		if err != nil {
			return
		}
		if n == 0 {
			break
		}
		if c.handler != nil {
			c.handler(n, buf)
		}

		c.input <- buf[:n]
	}
}
