package config

import (
	"github.com/saeed-kamyabi/dns-proxy/config/domain"
	"github.com/saeed-kamyabi/dns-proxy/config/usecase"
)

func NewConfigUseCase(path string) domain.UseCaseModel {
	uc := usecase.NewConfig(path)
	return uc
}
