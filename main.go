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
	Databases [][]DatabaseConfig `json:"databases"`

	LocationFile      string   `json:"location_file"`
	LocationFileField []string `json:"location_file_field"`
	LocationFileKey   string   `json:"location_file_key"`
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

	// Setup Engine
	SetupEngine(&config)

	// Get the database into memory
	for i, dbcs := range config.Databases {
		for _, dbc := range dbcs {
			SetupDatabase(&dbc, i)
		}
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
