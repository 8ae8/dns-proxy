package delivery

import (
	"fmt"
	"net"
)

func (s *server) serveHTTP() {
	fmt.Printf("listen %s://%s\n", s.Scheme(), s.config.GetAddress())

	l, err := net.Listen("tcp", s.config.GetAddress())
	if err != nil {
		fmt.Printf("error on listen http, err: %v\n", err)
	}
	defer func() {
		_ = l.Close()
	}()

	err = s.uc.Serve(s.config, l)
	if err != nil {
		panic(err)
	}
}
