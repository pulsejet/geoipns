package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/miekg/dns"
)

// Config is the configuration format
type Config struct {
	Databases []DatabaseConfig `json:"databases"`
}

func main() {
	// Open config file
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Read config
	var config Config
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	// Get the database into memory
	for _, dbc := range config.Databases {
		setupDatabase(&dbc)
	}

	// Attach handler function
	dns.HandleFunc("location.", handleDNSRequest)

	// Set up server
	port := 5312
	server := &dns.Server{Addr: "127.0.0.1:" + strconv.Itoa(port), Net: "udp"}

	// Log that we are starting server
	log.Printf("Starting Geolocation DNS server at %d\n", port)

	// Start listening
	err = server.ListenAndServe()

	// Shutdown when done
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
