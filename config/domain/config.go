package domain

type Config struct {
	DnsHost     string
	NameServers []string
	Domains     map[string]string
	Servers     Servers
}
