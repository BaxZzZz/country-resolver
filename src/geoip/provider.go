package geoip

import (
	"errors"
)

const (
	FREE_GEO_IP_NAME = "freegeoip.net"
	NEKUDO_NAME      = "geoip.nekudo.com"
)

// Information about IP address
type IPInfo struct {
	CountryName string
}

// GeoIP Provider interface
type Provider interface {
	GetIPInfo(ipAddress string) (*IPInfo, error)
}

// Creates new provider instance
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
