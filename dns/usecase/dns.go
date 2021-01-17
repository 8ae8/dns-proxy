package usecase

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"github.com/saeed-kamyabi/dns-proxy/dns/domain"
	"log"
)

type dnsUc struct {
	config domain.ConfigModel
}

func NewDns(config domain.ConfigModel) domain.UseCase {
	uc := &dnsUc{
		config: config,
	}
	return uc
}

func (uc *dnsUc) lookup(m *dns.Msg) (*dns.Msg, error) {
	dnsClient := new(dns.Client)
	dnsClient.Net = "udp"
	for _, server := range uc.config.GetSecondaryDnsServers() {
		response, _, err := dnsClient.Exchange(m, server)
		if err != nil {
			log.Printf("error on lookup Exchange for server %s, err: %v",
				server, err)
			continue
		}
		return response, nil
	}
	return nil, errors.New(fmt.Sprintf("could not resolve address %s",
		m.Question[0].Name))
}

func (uc *dnsUc) LookupHost(host string) (ip string, err error) {
	m := new(dns.Msg)
	m.SetQuestion(host, dns.TypeA)
	resp, err := uc.lookup(m)
	if err != nil {
		return "", err
	}
	if len(resp.Answer) == 0 {
		return
	}
	switch ans := resp.Answer[0].(type) {
	case *dns.A:
		return ans.A.String(), nil
	case *dns.CNAME:
		return ans.String(), nil
	}
	return "", nil
}

func (uc *dnsUc) GetResponse(requestMsg *dns.Msg) (*dns.Msg, error) {
	responseMsg := new(dns.Msg)
	if len(requestMsg.Question) > 0 {
		question := requestMsg.Question[0]
		switch question.Qtype {
		case dns.TypeA:
			answer, err := uc.processTypeA(&question, requestMsg)
			if err != nil {
				return responseMsg, err
			}
			responseMsg.Answer = append(responseMsg.Answer, *answer)
		default:
			answer, err := uc.processOtherTypes(&question, requestMsg)
			if err != nil {
				return responseMsg, err
			}
			responseMsg.Answer = append(responseMsg.Answer, *answer)
		}
	}

	return responseMsg, nil
}

func (uc *dnsUc) processOtherTypes(q *dns.Question,
	requestMsg *dns.Msg) (*dns.RR, error) {
	queryMsg := new(dns.Msg)
	requestMsg.CopyTo(queryMsg)
	queryMsg.Question = []dns.Question{*q}

	msg, err := uc.lookup(queryMsg)
	if err != nil {
		return nil, err
	}

	if len(msg.Answer) > 0 {
		return &msg.Answer[0], nil
	}
	return nil, fmt.Errorf("not found")
}

func (uc *dnsUc) processTypeA(q *dns.Question,
	requestMsg *dns.Msg) (*dns.RR, error) {
	if ip := uc.config.GetIP(q.Name); ip != "" {
		answer, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
		if err != nil {
			return nil, err
		}
		return &answer, nil
	} else {
		queryMsg := new(dns.Msg)
		requestMsg.CopyTo(queryMsg)
		queryMsg.Question = []dns.Question{*q}
		msg, err := uc.lookup(queryMsg)
		if err != nil {
			return nil, err
		}
		if len(msg.Answer) > 0 {
			return &msg.Answer[len(msg.Answer)-1], nil
		}
	}
	return nil, fmt.Errorf("not found")
}
