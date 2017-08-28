package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const configFilename = "config.json"

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	log.Println("Start resolver service")

	config, err := GetResolverConfig(configFilename)
	if err != nil {
		log.Fatal(err)
	}

	resolver, err := NewResolver(config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		sig := <-signals
		log.Println()
		log.Printf("Catching signal: %s", sig)
		resolver.Close()
		done <- true
	}()

	go resolver.Run()

	<-done
	log.Println("Stop resolver service")
}
