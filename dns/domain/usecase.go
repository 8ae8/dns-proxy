package domain

import "github.com/miekg/dns"

type UseCase interface {
	GetResponse(requestMsg *dns.Msg) (*dns.Msg, error)
	LookupHost(host string) (ip string, err error)
}
