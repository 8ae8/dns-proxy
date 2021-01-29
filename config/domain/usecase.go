package domain

type UseCaseModel interface {
	// name servers would be map containing information of
	// addr: 1.2.3.4
	// net: tcp
	// net: tcp-tls
	// net: udp
	GetNameServer() NameServer
	GetSecondaryNameServers() []NameServer
	GetIP(domain string) (ip string)
	NextServer() interface{}
}
