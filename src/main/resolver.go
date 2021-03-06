package main

import (
	"cache"
	"errors"
	"geoip"
	"log"
	"tcp"
	"time"
)

// GeoIP Resolver
type Resolver struct {
	tcpServer    *tcp.AsyncServer
	cache        *cache.LRUCache
	geoIPRequest *geoip.Request
}

// Get resolver config from configuration file
func GetResolverConfig(filename string) (*ResolverConfig, error) {
	config := &ResolverConfig{}

	if !config.Exists(filename) {
		config.SetDefault()
		config.WriteToFile(filename)
		return nil, errors.New("Can't find " + filename + " file, write default config")
	}

	err := config.ReadFromFile(filename)
	if err != nil {
		return nil, errors.New("Failed to read " + filename + " file")
	}

	return config, nil
}

// Creates new resolver instance
func NewResolver(config *ResolverConfig) (*Resolver, error) {
	providers, err := geoip.NewProviders(config.GeoIPProvider.Providers)
	if err != nil {
		return nil, err
	}

	geoIPRequest, err := geoip.NewRequest(providers, config.GeoIPProvider.RequestsLimit,
		time.Duration(config.GeoIPProvider.TimeIntervalMin)*time.Minute)
	if err != nil {
		return nil, err
	}

	store, err := cache.NewMongoDBStore(
		config.Cache.Address,
		config.Cache.DBName,
		config.Cache.Username,
		config.Cache.Password,
		config.Cache.Collection)
	if err != nil {
		return nil, err
	}

	lruCache, err := cache.NewLRUCache(config.Cache.ItemsLimit, store)
	if err != nil {
		return nil, err
	}

	server, err := tcp.NewServer(config.TcpServer.Address)
	if err != nil {
		return nil, err
	}

	resolver := &Resolver{
		tcpServer:    server,
		cache:        lruCache,
		geoIPRequest: geoIPRequest,
	}

	return resolver, nil
}

// Handler new connection with client
func (resolver *Resolver) handleNewClient(client *tcp.Client) {
	ipAddr, err := client.GetRemoteIpAddress()
	if err != nil {
		log.Printf("Can't get client IP address, error: %v", err)
		client.Close()
		return
	}

	log.Println("Client: " + ipAddr + " was connected.")

	ipInfo := &geoip.IPInfo{}
	err = resolver.cache.Get(ipAddr, ipInfo)
	if err == nil {
		client.SendMessage(ipInfo.CountryName + "\n")
		client.Close()
		return
	}

	ipInfo, err = resolver.geoIPRequest.GetIPInfo(ipAddr)
	if err != nil {
		log.Printf("Can't get IP information %s, error: %v", ipAddr, err)
		client.SendMessage(err.Error() + "\n")
		client.Close()
		return
	}

	err = resolver.cache.Set(ipAddr, ipInfo)
	if err != nil {
		log.Printf("Failed write info for IP - %s to local cache, error: %v", ipAddr, err)
	}

	client.SendMessage(ipInfo.CountryName + "\n")
	client.Close()
}

// Close resolver
func (resolver *Resolver) Close() {
	resolver.cache.Close()
	resolver.tcpServer.Shutdown()
}

// Run listen new connection
func (resolver *Resolver) Run() {
	resolver.tcpServer.OnNewClient(resolver.handleNewClient)
	resolver.tcpServer.Listen()
}
