package domain

type Config struct {
	NameServers []map[string]interface{}
	Domains     map[string]string
	Servers     Servers
}
