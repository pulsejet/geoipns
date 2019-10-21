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
	if config.HashMapFile == "" {
		return hm
	}

	// Open the file
	csvfile, err := os.Open(config.HashMapFile)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", config.HashMapFile, err)
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
	indices := make([]int, len(config.HashMapFileField))
	for i, f := range header {
		if f == config.HashMapFileKey {
			keyIndex = i
		} else {
			for j, x := range config.HashMapFileField {
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
