package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"sort"
	"strings"
)

// Main database
var dbs [][]*Database
var locmap map[string]string

// DatabaseRow epresents a single row in the databse
type DatabaseRow struct {
	IP         *net.IP
	Complement *net.IP
	IsHigh     bool
	Location   string
}

func (r DatabaseRow) getResponse(db *Database) string {
	// Lookup location
	response := ""
	if db.UseLocMap {
		response = locmap[r.Location]
	} else {
		response = r.Location
	}

	return strings.ReplaceAll(response, " ", "_")
}

// Database a database of GeoIP
type Database struct {
	Rows      []*DatabaseRow
	UseLocMap bool
}

func (db Database) Len() int { return len(db.Rows) }
func (db Database) Less(i, j int) bool {
	ri, rj := db.Rows[i], db.Rows[j]
	c := bytes.Compare(*ri.IP, *rj.IP)
	if c == 0 {
		if ri.IsHigh != rj.IsHigh {
			return ri.IsHigh
		}
		return bytes.Compare(*ri.Complement, *rj.Complement) > 0
	}
	return c < 0
}
func (db Database) Swap(i, j int) { db.Rows[i], db.Rows[j] = db.Rows[j], db.Rows[i] }

// Lookup the database
func (db Database) Lookup(lookupIP net.IP) (*DatabaseRow, error) {
	// Get the index to be inserted at
	i := sort.Search(db.Len(), func(i int) bool {
		return bytes.Compare(*db.Rows[i].IP, lookupIP) >= 0
	})

	// Tracker for HighIPs encountered
	numHigh := 0

	// Check if index matches
	if i < db.Len() {
		// Check for immediate match
		if bytes.Compare(*db.Rows[i].IP, lookupIP) == 0 {
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
			if bytes.Compare(*row.IP, lookupIP) == 0 || (!row.IsHigh && numHigh <= 0) {
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

	return &DatabaseRow{}, errors.New("Not Found")
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
	LowIP    string `json:"low_ip"`
	HighIP   string `json:"high_ip"`
	Location string `json:"location"`
}

type dbFieldIndex struct {
	CIDR     int
	LowIP    int
	HighIP   int
	Location int
}

// SetupEngine initializes the engine
func SetupEngine(config *Config) {
	dbs = make([][]*Database, 0)
	locmap = initializeLocationMap(config)
}

// SetupDatabase caches the databse in memory
func SetupDatabase(dbc *DatabaseConfig, index int) {
	// Initialize
	mdb := Database{make([]*DatabaseRow, 0), dbc.UseLocMap}

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
		case dbc.Fields.HighIP:
			indices.HighIP = i
		case dbc.Fields.Location:
			indices.LowIP = i
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

		// Start and end IP addresses
		var lowIP net.IP
		var highIP net.IP

		// Check if CIDR is to be parsed
		if indices.CIDR == -1 {
			plowIP := net.ParseIP(record[indices.LowIP])
			phighIP := net.ParseIP(record[indices.HighIP])
			if plowIP == nil || phighIP == nil {
				log.Panicln("Failed to parse", record[indices.LowIP], record[indices.HighIP])
				continue
			}
			lowIP = plowIP.To16()
			highIP = phighIP.To16()
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
			lowIP = n.IP.To16()

			// Get the upper IP
			for i := range n.IP {
				n.IP[i] |= ^n.Mask[i]
			}
			highIP = n.IP.To16()
		}

		// Get low row
		lowRow := DatabaseRow{
			IP:         &lowIP,
			Complement: &highIP,
			IsHigh:     false,
			Location:   record[indices.Location],
		}
		mdb.Rows = append(mdb.Rows, &lowRow)

		// Get high Row
		highRow := DatabaseRow{
			IP:         &highIP,
			Complement: &lowIP,
			IsHigh:     true,
			Location:   record[indices.Location],
		}
		mdb.Rows = append(mdb.Rows, &highRow)
	}

	// Sort the database
	sort.Sort(mdb)

	// Add database to databases
	if len(dbs) <= index {
		dbs = append(dbs, make([]*Database, 0))
	}
	dbs[index] = append(dbs[index], &mdb)

	// Run garbage collection
	runtime.GC()

	fmt.Println("Read database", dbc.File)
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
	lookupIP := ip.To16()

	// Lookup all databases
	response := ""
	for _, dbl := range dbs {
		for _, db := range dbl {
			row, err := db.Lookup(lookupIP)
			if err == nil {
				response += row.getResponse(db)
				break
			}
		}
		response += "|"
	}

	return response + "S"
}
