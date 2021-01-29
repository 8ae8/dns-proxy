package domain

type ConfigModel interface {
	// name servers would be map containing information of
	// host: 1.2.3.4
	// net: tcp
	// net: tcp-tls
	// net: udp
	GetNameServer() map[string]interface{}
	GetSecondaryNameServers() []map[string]interface{}
	GetIP(address string) (ip string)

	GetNsHost(ns map[string]interface{}) string
	GetNsNet(ns map[string]interface{}) string
}
