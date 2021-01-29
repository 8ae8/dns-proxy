package domain

import (
	cd "github.com/saeed-kamyabi/dns-proxy/config/domain"
)

type ConfigModel interface {
	// name server contains information of exposing name server
	GetNameServer() cd.NameServer
	// name servers would be list of servers containing
	// information of secondary name servers
	GetSecondaryNameServers() []cd.NameServer
	GetIP(address string) (ip string)
}
