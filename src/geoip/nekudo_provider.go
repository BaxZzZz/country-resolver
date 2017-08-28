package geoip

import (
	"encoding/json"
)

const nekudoURL = "http://geoip.nekudo.com/api/"

// Information from eoip.nekudo.com about IP address country
type nekudoCountry struct {
	Name string `json:"name"`
}

// Information from eoip.nekudo.com about IP address
type nekudoIPInfo struct {
	Country *nekudoCountry `json:"country"`
}

// Provider eoip.nekudo.com resource
type nekudoProvider struct {
	client Client
}

// Get information by IP address
func (provider *nekudoProvider) GetIPInfo(ipAddress string) (*IPInfo, error) {
	data, err := provider.client.Request(nekudoURL + ipAddress)
	if err != nil {
		return nil, err
	}

	info := &nekudoIPInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	ipInfo := &IPInfo{}

	if info.Country != nil {
		ipInfo.CountryName = info.Country.Name
	}

	return ipInfo, nil
}
