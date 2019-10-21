package geoip

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func initializeLocationMap(config *Config) map[string]string {
	lmap := map[string]string{}

	// Check if configuration present
	if config.LocationFile == "" {
		return lmap
	}

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
	locIndices := make([]int, len(config.LocationFileField))
	for i, f := range header {
		if f == config.LocationFileKey {
			keyIndex = i
		} else {
			for j, x := range config.LocationFileField {
				if x == f {
					locIndices[j] = i
				}
			}
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

		// Get location
		loc := ""
		for _, i := range locIndices {
			if i != -1 && record[i] != "" {
				if loc == "" {
					loc += record[i]
				} else {
					loc += ", " + record[i]
				}
			}
		}
		lmap[record[keyIndex]] = loc
	}

	return lmap
}
