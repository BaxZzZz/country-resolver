package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BaxZzZz/country-resolver/geoip"
	"github.com/BaxZzZz/country-resolver/tcp"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "33333"
)

func main() {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Start country-resolver service")

	providers, err := geoip.NewProviders([]string{geoip.FREE_GEO_IP_NAME, geoip.NEKUDO_NAME})

	if err != nil {
		log.Fatalf("GeoIP providers failed to create, error: %v", err)
	}

	ipInfoRequest, err := geoip.NewRequest(providers, 1, time.Duration(1)*time.Minute)

	if err != nil {
		log.Fatalf("IP info provider manager failed to create, error: %v", err)
	}

	server := tcp.NewTcpServer()
	server.OnNewClient(func(client *tcp.Client) {
		ipAddr, err := client.GetRemoteIpAddress()
		if err != nil {
			log.Printf("Can't get client IP address, error: %v", err)
			client.Close()
			return
		}

		ipInfo, err := ipInfoRequest.GetIPInfo(ipAddr)
		if err != nil {
			log.Printf("Can't get IP information, error: %v", err)
			client.SendMessage(err.Error() + "\n")
			client.Close()
			return
		}

		client.SendMessage(ipInfo.CountryName + "\n")
		client.Close()
	})

	err = server.Start(CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		log.Fatalf("TCP server failed to start, error: %v", err)
	}

	<-signals

	log.Println("Stop country-resolver service")
	server.Stop()
}
