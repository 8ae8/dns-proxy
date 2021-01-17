package dns

import (
	"github.com/saeed-kamyabi/dns-proxy/dns/delivery"
	"github.com/saeed-kamyabi/dns-proxy/dns/domain"
	"github.com/saeed-kamyabi/dns-proxy/dns/usecase"
)

func NewDnsUseCase(config domain.ConfigModel) domain.UseCase {
	uc := usecase.NewDns(config)
	return uc
}

func NewServer(config domain.ConfigModel, uc domain.UseCase) domain.ServerModel {
	server := delivery.NewServer(config, uc)
	return server
}
