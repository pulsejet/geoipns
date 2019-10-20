package main

import (
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
)

// Main database
var dbs []Database
var locmap map[string]string

// DatabaseRow epresents a single row in the databse
type DatabaseRow struct {
	IP       string
	IsHigh   bool
	Location string
}

func (r DatabaseRow) getIP() (string, error) {
	bin, err := hex.DecodeString(r.IP)
	if err != nil {
		return r.IP, err
	}
	return net.IP(bin).String(), nil
}

func (r DatabaseRow) getResponse(db *Database) string {
	if db.UseLocMap {
		return locmap[r.Location]
	}
	return r.Location
}

// Database a database of GeoIP
type Database struct {
	Rows      []DatabaseRow
	UseLocMap bool
}

func (db Database) Len() int           { return len(db.Rows) }
func (db Database) Less(i, j int) bool { return db.Rows[i].IP < db.Rows[j].IP }
func (db Database) Swap(i, j int)      { db.Rows[i], db.Rows[j] = db.Rows[j], db.Rows[i] }

// Lookup the database
func (db Database) Lookup(hexIP string) (DatabaseRow, error) {
	// Get the index to be inserted at
	i := sort.Search(db.Len(), func(i int) bool {
		return db.Rows[i].IP >= hexIP
	})

	// Tracker for HighIPs encountered
	numHigh := 0

	// Check if index matches
	if i < db.Len() {
		// Check for immediate match
		if db.Rows[i].IP == hexIP {
			return db.Rows[i], nil
		}

		// Go back five paces at most
		for j := 1; j <= 5; j++ {
			// Look out for invalid calls
			if i-j < 0 {
				break
			}

			// Get the row
			row := db.Rows[i-j]

			// Check if IP matches or unbalanced LowIP
			if row.IP == hexIP || (!row.IsHigh && numHigh <= 0) {
				return row, nil
			}

			// Increment counter if high IP
			if row.IsHigh {
				numHigh++
			} else {
				numHigh--
			}
		}
	}

	return DatabaseRow{}, errors.New("Not Found")
}

// DatabaseConfig is the format of configuration for geoip db
type DatabaseConfig struct {
	File      string               `json:"file"`
	Fields    DatabaseConfigFields `json:"fields"`
	UseLocMap bool                 `json:"use_loc_map"`
}

// DatabaseConfigFields is the format of fields for geoip db
type DatabaseConfigFields struct {
	CIDR     string `json:"cidr"`
	Location string `json:"location"`
}

type dbFieldIndex struct {
	CIDR     int
	Location int
}

// SetupEngine initializes the engine
func SetupEngine(config *Config) {
	dbs = make([]Database, 0)
	locmap = initializeLocationMap(config)
}

// SetupDatabase caches the databse in memory
func SetupDatabase(dbc *DatabaseConfig) {
	// Initialize
	mdb := Database{make([]DatabaseRow, 0), dbc.UseLocMap}

	// Indices of fields
	indices := dbFieldIndex{CIDR: -1, Location: -1}

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
		case dbc.Fields.Location:
			indices.Location = i
		}
	}

	k := 0
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
		}

		// Get low row
		lowRow := DatabaseRow{
			IP:       lowIP,
			IsHigh:   false,
			Location: record[indices.Location],
		}
		mdb.Rows = append(mdb.Rows, lowRow)

		// Get high Row
		highRow := DatabaseRow(lowRow)
		highRow.IP = highIP
		highRow.IsHigh = true
		mdb.Rows = append(mdb.Rows, highRow)

		k++
		if k > 4 {
			break
		}
	}

	// Sort the database
	sort.Sort(mdb)

	// Add database to databases
	dbs = append(dbs, mdb)

	// Print database
	for _, x := range mdb.Rows {
		ipx, _ := x.getIP()
		fmt.Println(ipx, x.IsHigh, x.Location)
	}

	fmt.Println("Read it all!")
}

func unknownResponse() string {
	return "UnknownLocation"
}

// GeoHandle returns ip data
func GeoHandle(ipstr string) string {
	// Parse the IP to bytes
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return unknownResponse()
	}

	// Get hexadecimal for lookup
	hexIP := hex.EncodeToString(ip.To16())

	// Lookup all databases
	for _, db := range dbs {
		row, err := db.Lookup(hexIP)
		if err == nil {
			return row.getResponse(&db)
		}
	}

	// Fallback
	return unknownResponse()
}
