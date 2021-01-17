package delivery

import (
	"fmt"
	"github.com/saeed-kamyabi/dns-proxy/proxy/domain"
)

type server struct {
	config domain.ServerConfig

	uc domain.UseCase
}

func SetupServers(config domain.ConfigModel,
	uc domain.UseCase) []domain.Server {
	var result []domain.Server
	for {
		if n := config.NextServer(); n != nil {
			conf := n.(domain.ServerConfig)
			s := &server{
				config: conf,
				uc:     uc,
			}
			result = append(result, s)
		} else {
			break
		}
	}
	return result
}

func (s *server) Scheme() string {
	return s.config.GetScheme()
}

func (s *server) ListenAndServe() {
	switch s.Scheme() {
	//case "udp":
	//	s.serveUDP()
	//case "tcp":
	//	s.serveTCP()
	//case "tls":
	//	s.serveTLS()
	case "http":
		s.serveHTTP()
	case "https":
		s.serveHTTP()
	default:
		panic(fmt.Sprintf("unsupported scheme %s", s.Scheme()))
	}
}
