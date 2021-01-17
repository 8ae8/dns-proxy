package usecase

import (
	"errors"
	"fmt"
	"github.com/saeed-kamyabi/dns-proxy/proxy/domain"
	"net"
)

type proxy struct {
	dns domain.Dns
}

func NewProxy(dns domain.Dns) domain.UseCase {
	uc := &proxy{
		dns: dns,
	}
	return uc
}

func (uc *proxy) Serve(serverConfig domain.ServerConfig, l net.Listener) error {
	scheme := serverConfig.GetScheme()
	switch scheme {
	case "http":
		s := newHttpServer(uc.dns)
		return s.Serve(serverConfig, l)
	case "https":
		s := newHttpsServer(uc.dns)
		return s.Serve(serverConfig, l)
	}
	return errors.New(fmt.Sprintf("scheme %s not supported", scheme))
}
