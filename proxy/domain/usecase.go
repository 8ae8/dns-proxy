package domain

import "net"

type UseCase interface {
	Serve(serverConfig ServerConfig, l net.Listener) error
}
