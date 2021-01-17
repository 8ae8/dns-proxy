package domain

type UseCaseModel interface {
	// GetSecondaryDnsServers get list of other dns servers rather than
	// current server for resolving non-registered domains
	GetSecondaryDnsServers() []string
	// GetDnsServerHost determines listening host address of dns server
	GetDnsServerHost() string
	GetIP(domain string) (ip string)
	NextServer() interface{}
}
