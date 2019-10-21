package geoip

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func initializeHashMap(config *Config) map[string]string {
	hm := map[string]string{}

	// Check if configuration present
	if config.LocationFile == "" {
		return hm
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
	indices := make([]int, len(config.LocationFileField))
	for i, f := range header {
		if f == config.LocationFileKey {
			keyIndex = i
		} else {
			for j, x := range config.LocationFileField {
				if x == f {
					indices[j] = i
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
		val := ""
		for _, i := range indices {
			if i != -1 && record[i] != "" {
				if val == "" {
					val += record[i]
				} else {
					val += ", " + record[i]
				}
			}
		}
		hm[record[keyIndex]] = val
	}

	return hm
}
