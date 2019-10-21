package geoip

import (
	"bytes"
	"encoding/csv"
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

func (r DatabaseRow) getResponse(db *Database) string {
	// Lookup location
	response := ""
	if db.UseLocMap {
		response = locmap[*r.Location]
	} else {
		response = *r.Location
	}

	return strings.ReplaceAll(response, " ", "_")
}

func (db Database) Len() int { return len(db.Rows) }
func (db Database) Less(i, j int) bool {
	ri, rj := db.Rows[i], db.Rows[j]
	c := bytes.Compare(*ri.IP, *rj.IP)
	if c == 0 {
		if ri.IsHigh != rj.IsHigh {
			return (ri.IP == rj.Complement) != ri.IsHigh
		}
		return bytes.Compare(*ri.Complement, *rj.Complement) > 0
	}
	return c < 0
}
func (db Database) Swap(i, j int) { db.Rows[i], db.Rows[j] = db.Rows[j], db.Rows[i] }

// Lookup the database
func (db Database) Lookup(lookupIP net.IP) *DatabaseRow {
	// Get the index to be inserted at
	i := sort.Search(db.Len(), func(i int) bool {
		c := bytes.Compare(*db.Rows[i].IP, lookupIP)
		return c > 0 || (c == 0 && db.Rows[i].IsHigh)
	})

	// Check if index matches
	if i > 0 && i < db.Len() {
		// Get the row
		row := db.Rows[i-1]

		// Check if lowIP
		if !row.IsHigh {
			return row
		}

		// Return parent if highIP
		return row.Parent
	}

	return nil
}

// SetupEngine initializes the engine
func SetupEngine(config *Config) {
	dbs = make([][]*Database, 0)
	locmap = initializeLocationMap(config)

	// Get the database into memory
	for i, dbcs := range config.Databases {
		for _, dbc := range dbcs {
			setupDatabase(&dbc, i)
		}
	}
}

// setupDatabase caches the databse in memory
func setupDatabase(dbc *DatabaseConfig, index int) {
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
		case dbc.Fields.LowIP:
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

		// Trim all space
		for i, x := range record {
			record[i] = strings.TrimSpace(x)
		}

		// Start and end IP addresses
		var lowIP net.IP
		var highIP net.IP

		// Check if CIDR is to be parsed
		if indices.CIDR == -1 || record[indices.CIDR] == "" {
			plowIP := net.ParseIP(record[indices.LowIP])
			phighIP := net.ParseIP(record[indices.HighIP])
			if plowIP == nil || phighIP == nil {
				log.Println("Failed to parse", record[indices.LowIP], record[indices.HighIP])
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

		// Pointer to location
		location := record[indices.Location]

		// Get low row
		lowRow := DatabaseRow{
			IP:         &lowIP,
			Complement: &highIP,
			IsHigh:     false,
			Location:   &location,
		}
		mdb.Rows = append(mdb.Rows, &lowRow)

		// Get high Row
		highRow := DatabaseRow{
			IP:         &highIP,
			Complement: &lowIP,
			IsHigh:     true,
			Location:   &location,
		}
		mdb.Rows = append(mdb.Rows, &highRow)
	}

	// Sort the database
	sort.Sort(mdb)

	// Set parent values
	parents := rowStack{}
	for _, row := range mdb.Rows {
		// Pop parent if highIP
		if row.IsHigh {
			parents, _ = parents.Pop()
		}

		// Set value if has a parent
		if p := parents.Peek(); p != nil {
			row.Parent = p
		}

		// Push parent if lowIP
		if !row.IsHigh {
			parents = parents.Push(row)
		}
	}

	// Add database to databases
	if len(dbs) <= index {
		dbs = append(dbs, make([]*Database, 0))
	}
	dbs[index] = append(dbs[index], &mdb)

	// Run garbage collection
	runtime.GC()
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
			row := db.Lookup(lookupIP)
			if row != nil {
				response += row.getResponse(db)
				break
			}
		}
		response += "|"
	}

	return response + "S"
}