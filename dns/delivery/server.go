package delivery

import (
	"github.com/miekg/dns"
	"github.com/saeed-kamyabi/dns-proxy/dns/domain"
	"log"
	"regexp"
)

type server struct {
	dns.Server
	config domain.ConfigModel
	uc     domain.UseCase
}

func NewServer(config domain.ConfigModel, uc domain.UseCase) domain.ServerModel {
	ns := config.GetNameServer()
	s := &server{
		Server: dns.Server{
			Addr: ns["host"].(string),
			Net:  ns["net"].(string),
		},
		config: config,
		uc:     uc,
	}
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		switch r.Opcode {
		case dns.OpcodeQuery:
			m, err := uc.GetResponse(r)
			if err != nil {
				log.Printf("Failed lookup for %s with error: %s\n",
					r, err.Error())
				if m != nil {
					m.SetReply(r)
					_ = w.WriteMsg(m)
				}
				return
			}
			if len(m.Answer) > 0 {
				pattern := regexp.MustCompile(`\b(?:(?:2(?:[0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9])\.){3}(?:(?:2([0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9]))\b`)
				ipAddress := pattern.FindAllString(m.Answer[0].String(), -1)

				if len(ipAddress) > 0 {
					log.Printf("Lookup for %s with ip %s\n",
						m.Answer[0].Header().Name, ipAddress[0])
				} else {
					log.Printf("Lookup for %s with response %s\n",
						m.Answer[0].Header().Name, m.Answer[0])
				}
			}
			m.SetReply(r)
			_ = w.WriteMsg(m)
		}
	})
	return s
}
