package geoip

import (
	"errors"
)

const (
	FREE_GEO_IP_NAME = "freegeoip.net"
	NEKUDO_NAME      = "geoip.nekudo.com"
)

type IPInfo struct {
	CountryName string
}

type Provider interface {
	GetIpInfo(ipAddress string) (*IPInfo, error)
}

func NewProviders(providerNames []string) ([]Provider, error) {
	var providers []Provider
	for _, providerName := range providerNames {
		switch providerName {
		case FREE_GEO_IP_NAME:
			providers = append(providers, &freeGeoIPProvider{
				client: &httpClient{},
			})
		case NEKUDO_NAME:
			providers = append(providers, &nekudoProvider{
				client: &httpClient{},
			})
		default:
			return nil, errors.New("Unknown provider name: " + providerName)
		}
	}

	return providers, nil
}
