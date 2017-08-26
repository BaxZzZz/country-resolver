package main

import (
	"log"
	"os"
	"time"

	"github.com/BaxZzZz/country-resolver/geoip"
	"github.com/BaxZzZz/country-resolver/tcp"
)

const configFilename = "config.json"

func runApplication(config AppConfig) {
	log.Println("Start country-resolver service")

	providers, err := geoip.NewProviders(config.GeoIPProvider.Providers)

	if err != nil {
		log.Fatalf("GeoIP providers failed to create, error: %v", err)
	}

	timeIntervalMin := time.Duration(config.GeoIPProvider.TimeIntervalMin) * time.Minute
	ipInfoRequest, err := geoip.NewRequest(providers, 1, timeIntervalMin)

	if err != nil {
		log.Fatalf("IP info provider manager failed to create, error: %v", err)
	}

	server, err := tcp.NewServer(config.TcpServer.Address)
	if err != nil {
		log.Fatalf("TCP server failed to start, error: %v", err)
	}

	server.OnNewClient(func(client *tcp.Client) {
		ipAddr, err := client.GetRemoteIpAddress()
		if err != nil {
			log.Printf("Can't get client IP address, error: %v", err)
			client.Close()
			return
		}

		ipInfo, err := ipInfoRequest.GetIPInfo(ipAddr)
		if err != nil {
			log.Printf("Can't get IP information %s, error: %v", ipAddr, err)
			client.SendMessage(err.Error() + "\n")
			client.Close()
			return
		}

		client.SendMessage(ipInfo.CountryName + "\n")
		client.Close()
	})

	server.Listen()

	log.Println("Stop country-resolver service")
}

func main() {
	config := AppConfig{}

	if !config.Exists(configFilename) {
		log.Println("Can't find " + configFilename + " file, write default config.")
		config.SetDefault()
		config.WriteToFile(configFilename)
		os.Exit(1)
	}

	err := config.ReadFromFile(configFilename)
	if err != nil {
		log.Fatalf("Failed to read " + configFilename + " file")
	}

	runApplication(config)

}
