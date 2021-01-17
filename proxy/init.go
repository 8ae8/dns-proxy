package proxy

import (
	"github.com/saeed-kamyabi/dns-proxy/proxy/delivery"
	"github.com/saeed-kamyabi/dns-proxy/proxy/domain"
	"github.com/saeed-kamyabi/dns-proxy/proxy/usecase"
)

func NewProxyUseCase(dns domain.Dns) domain.UseCase {
	return usecase.NewProxy(dns)
}

func SetupServers(config domain.ConfigModel,
	uc domain.UseCase) []domain.Server {
	return delivery.SetupServers(config, uc)
}
