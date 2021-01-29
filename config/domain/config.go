package domain

import (
	"crypto/tls"
	"github.com/miekg/dns"
)

type NameServer struct {
	*dns.Server
	TLSConfig *struct {
		Certificates []struct {
			CertFile string
			KeyFile  string
		}
	}
}

type Config struct {
	NameServers []NameServer
	Domains     map[string]string
	Servers     Servers
}

func (c *Config) Initialize() {
	for _, ns := range c.NameServers {
		if ns.TLSConfig != nil {
			if ns.Server.TLSConfig == nil {
				ns.Server.TLSConfig = new(tls.Config)
			}
			for _, cert := range ns.TLSConfig.Certificates {
				cert, err := tls.LoadX509KeyPair(cert.CertFile, cert.KeyFile)
				if err != nil {
					panic(err)
				}
				ns.Server.TLSConfig.Certificates = append(ns.Server.
					TLSConfig.Certificates, cert)
			}
		}
	}
}
