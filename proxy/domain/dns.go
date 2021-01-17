package domain

type Dns interface {
	LookupHost(host string) (ip string, err error)
}
