package geoip

import (
	"encoding/json"
)

const freeGeoIPURL = "http://freegeoip.net/json/"

type freeGeoIPInfo struct {
	CountryName string `json:"country_name"`
}

type freeGeoIPProvider struct {
	client Client
}

func (provider *freeGeoIPProvider) GetIpInfo(ipAddress string) (*IPInfo, error) {
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
