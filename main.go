package main

import "log"

const configFilename = "config.json"

func main() {
	log.Println("Start country-resolver service")

	config, err := GetResolverConfig(configFilename)
	if err != nil {
		log.Fatal(err)
	}

	resolver, err := NewResolver(config)
	if err != nil {
		log.Fatal(err)
	}

	resolver.Run()
}
