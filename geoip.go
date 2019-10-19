package main

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// DatabaseConfig is the format of configuration for geoip db
type DatabaseConfig struct {
	File   string               `json:"file"`
	Fields DatabaseConfigFields `json:"fields"`
}

// DatabaseConfigFields is the format of fields for geoip db
type DatabaseConfigFields struct {
	CIDR string `json:"cidr"`
}

type dbFieldIndex struct {
	CIDR int
}

// SetupDatabase caches the databse in memory
func SetupDatabase(dbc *DatabaseConfig) {
	// Indices of fields
	indices := dbFieldIndex{CIDR: -1}

	// Open the file
	csvfile, err := os.Open(dbc.File)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", dbc.File, err)
	}
	defer csvfile.Close()

	// Open a new reader
	r := csv.NewReader(csvfile)

	// Get header
	header, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	for i, f := range header {
		switch f {
		case dbc.Fields.CIDR:
			indices.CIDR = i
		}
	}

	k := 0
	for {
		// Read record
		record, err := r.Read()

		// Check i file ended
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Start and end IP addresses
		var lowIP string
		var highIP string

		// Check if CIDR is to be parsed
		if indices.CIDR == -1 {
			// TODO: NO CIDR
		} else {
			// Get CIDR
			cidr := record[indices.CIDR]

			// Parse CIDR
			_, n, err := net.ParseCIDR(cidr)
			if err != nil {
				log.Println(err)
				continue
			}

			// Get the lower IP
			for i := range n.IP {
				n.IP[i] &= n.Mask[i]
			}
			lowIP = hex.EncodeToString(n.IP.To16())

			// Get the upper IP
			for i := range n.IP {
				n.IP[i] |= ^n.Mask[i]
			}
			highIP = hex.EncodeToString(n.IP.To16())

			fmt.Printf("%s to %s\n", lowIP, highIP)
		}

		k++
		if k > 4 {
			break
		}
	}

	fmt.Printf("Read it all!\n")
}

// GeoHandle returns ip data
func GeoHandle(ipstr string) string {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return ipstr
	}
	return hex.EncodeToString(ip.To16())
}
