package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type DatabaseConfig struct {
	File string `json:"file"`
}

func setupDatabase(dbc *DatabaseConfig) {
	// Open the file
	csvfile, err := os.Open(dbc.File)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", dbc.File, err)
	}
	defer csvfile.Close()

	// Open a new reader
	r := csv.NewReader(csvfile)

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

		fmt.Println(record[0], record[1])

		break
	}

	fmt.Printf("Read it all!\n")
}
