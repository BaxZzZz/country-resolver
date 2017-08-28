package main

import (
	"encoding/json"
	"geoip"
	"io/ioutil"
	"os"
)

type TcpServerConfig struct {
	Address string `json:"address"`
}

type MongoDBCache struct {
	Address    string `json:"address"`
	DBName     string `json:"db_name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Collection string `json:"collection"`
	ItemsLimit uint   `json:"items_limit"`
}

type GeoIPProviderConfig struct {
	Providers       []string `json:"providers"`
	RequestsLimit   uint     `json:"requests_limit"`
	TimeIntervalMin uint     `json:"time_interval_min"`
}

type ResolverConfig struct {
	TcpServer     TcpServerConfig     `json:"tcp_server"`
	GeoIPProvider GeoIPProviderConfig `json:"geo_ip_provider"`
	Cache         MongoDBCache        `json:"cache"`
}

func (config *ResolverConfig) WriteToFile(filename string) error {
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

func (config *ResolverConfig) ReadFromFile(filename string) error {
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

func (config *ResolverConfig) Exists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (config *ResolverConfig) SetDefault() {
	config.TcpServer.Address = "0.0.0.0:9999"
	config.GeoIPProvider.Providers = []string{geoip.FREE_GEO_IP_NAME, geoip.NEKUDO_NAME}
	config.GeoIPProvider.RequestsLimit = 100
	config.GeoIPProvider.TimeIntervalMin = 1
	config.Cache.Address = "localhost"
	config.Cache.Username = "test"
	config.Cache.Password = "test"
	config.Cache.DBName = "resolver"
	config.Cache.Collection = "cache"
	config.Cache.ItemsLimit = 100000
}
