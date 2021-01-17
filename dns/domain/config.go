package domain

type ConfigModel interface {
	GetSecondaryDnsServers() []string
	GetDnsServerHost() string
	GetIP(address string) (ip string)
}
