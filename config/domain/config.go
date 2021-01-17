package domain

type Config struct {
	DnsServerHost string
	DnsServers    []string
	Domains       map[string]string
	Servers       Servers
}
