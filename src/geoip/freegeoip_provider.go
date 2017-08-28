package geoip

import (
	"encoding/json"
)

const freeGeoIPURL = "http://freegeoip.net/json/"

// Information from Freegeoip.net about IP address
type freeGeoIPInfo struct {
	CountryName string `json:"country_name"`
}

// Provider freegeoip.net resource
type freeGeoIPProvider struct {
	client Client
}

// Get information by IP address
func (provider *freeGeoIPProvider) GetIPInfo(ipAddress string) (*IPInfo, error) {
	data, err := provider.client.Request(freeGeoIPURL + ipAddress)

	if err != nil {
		return nil, err
	}

	info := &freeGeoIPInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	ipInfo := &IPInfo{}
	ipInfo.CountryName = info.CountryName

	return ipInfo, nil
}
