package geoip

import (
	"encoding/json"
)

const nekudoURL = "http://geoip.nekudo.com/api/"

type nekudoCountry struct {
	Name string `json:"name"`
}

type nekudoIPInfo struct {
	Country *nekudoCountry `json:"country"`
}

type nekudoProvider struct {
	client Client
}

func (provider *nekudoProvider) GetIpInfo(ipAddress string) (*IPInfo, error) {
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
