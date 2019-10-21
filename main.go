package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/miekg/dns"
	g "github.com/pulsejet/geoipns/geoip"
)

type config struct {
	g.Config
	Address string `json:"address"`
	Suffix  string `json:"suffix"`
}

var mConfig config

// SetupEnvironment loads conig and initializes databases
func SetupEnvironment(configFile string) {
	// Open config file
	jsonFile, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Read config
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &mConfig)

	// Setup Engine
	g.SetupEngine(&mConfig.Config)
}

func main() {
	// Set everything up
	SetupEnvironment("config.json")

	// Attach handler function
	dns.HandleFunc(mConfig.Suffix+".", handleDNSRequest)

	// Set up server
	server := &dns.Server{Addr: mConfig.Address, Net: "udp"}

	// Log that we are starting server
	log.Println("Starting Geolocation DNS server at", mConfig.Address)

	// Start listening
	err := server.ListenAndServe()

	// Shutdown when done
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
