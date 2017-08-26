package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/BaxZzZz/country-resolver/geoip"
)

type TcpServerConfig struct {
	Address string `json:"address"`
}

type DBCache struct {
	MongoDBURL string `json:"mongo_db_url"`
	ItemsLimit uint   `json:"items_limit"`
}

type GeoIPProviderConfig struct {
	Providers       []string `json:"providers"`
	RequestsLimit   uint     `json:"requests_limit"`
	TimeIntervalMin uint     `json:"time_interval_min"`
}

type AppConfig struct {
	TcpServer     TcpServerConfig     `json:"tcp_server"`
	GeoIPProvider GeoIPProviderConfig `json:"geo_ip_provider"`
	Cache         DBCache             `json:"cache"`
}

func (config *AppConfig) WriteToFile(filename string) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		return err
	}

	return nil
}

func (config *AppConfig) ReadFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}

func (config *AppConfig) Exists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (config *AppConfig) SetDefault() {
	config.TcpServer.Address = "0.0.0.0:9999"
	config.GeoIPProvider.Providers = []string{geoip.FREE_GEO_IP_NAME, geoip.NEKUDO_NAME}
	config.GeoIPProvider.RequestsLimit = 10
	config.GeoIPProvider.TimeIntervalMin = 1
	config.Cache.MongoDBURL = "mongodb://localhost/country_resolver"
	config.Cache.ItemsLimit = 100000
}
