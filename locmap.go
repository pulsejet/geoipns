package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func initializeLocationMap(config *Config) map[string]string {
	lmap := map[string]string{}
	// Open the file
	csvfile, err := os.Open(config.LocationFile)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", config.LocationFile, err)
	}
	defer csvfile.Close()

	// Open a new reader
	r := csv.NewReader(csvfile)

	// Get header
	header, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Get indices
	keyIndex := -1
	locIndex := -1
	for i, f := range header {
		switch f {
		case config.LocationFileKey:
			keyIndex = i
		case config.LocationFileField:
			locIndex = i
		}
	}

	for {
		// Read record
		record, err := r.Read()

		// Check if file ended
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		lmap[record[keyIndex]] = record[locIndex]
	}

	return lmap
}
