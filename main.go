package main

import (
	"github.com/saeed-kamyabi/dns-proxy/config"
	"github.com/saeed-kamyabi/dns-proxy/dns"
	"github.com/saeed-kamyabi/dns-proxy/proxy"
	pd "github.com/saeed-kamyabi/dns-proxy/proxy/domain"
	"log"
	"os"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.json"
	}
	cfg := config.NewConfigUseCase(configPath)
	// dns server
	dnsUc := dns.NewDnsUseCase(cfg)
	dnsServer := dns.NewServer(cfg, dnsUc)
	// proxy server
	proxyUc := proxy.NewProxyUseCase(dnsUc)
	servers := proxy.SetupServers(cfg, proxyUc)
	// run
	for _, server := range servers {
		go func(s pd.Server) {
			log.Printf("Starting proxy server %s\n", s.Scheme())
			s.ListenAndServe()
		}(server)
	}
	log.Printf("Starting dns server at %s\n", cfg.GetNameServer().Addr)
	err := dnsServer.ListenAndServe()
	if err != nil {
		log.Printf("Failed to start server: %s\n ", err.Error())
	}
}
