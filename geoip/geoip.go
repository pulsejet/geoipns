package geoip

import (
	"bytes"
	"encoding/csv"
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
var databaseSets []*DatabaseSet
var hashMap map[string]string

func (r DatabaseRow) getResponse(db *Database) string {
	// Lookup data
	if db.UseHashMap {
		return hashMap[*r.Data]
	}
	return *r.Data
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
	databaseSets = make([]*DatabaseSet, 0)
	hashMap = initializeHashMap(config)

	// Get the database into memory
	for _, dbcs := range config.DatabaseSets {
		// Create new database set
		dbset := &DatabaseSet{AttributeName: dbcs.AttributeName}
		databaseSets = append(databaseSets, dbset)

		// Add each database to the set
		for _, dbc := range dbcs.Databases {
			mdb := setupDatabase(&dbc)
			dbset.Databases = append(dbset.Databases, mdb)
		}
	}
}

// setupDatabase caches the databse in memory
func setupDatabase(dbc *DatabaseConfig) *Database {
	// Initialize
	mdb := &Database{make([]*DatabaseRow, 0), dbc.UseHashMap}

	// Indices of fields
	indices := dbFieldIndex{CIDR: -1, Data: -1}

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
		case dbc.Fields.Data:
			indices.Data = i
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
			if bytes.Compare(plowIP, phighIP) > 0 {
				log.Println("Ignoring incorrect record", record[indices.LowIP], record[indices.HighIP])
				continue
			}
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

		// Pointer to data
		data := record[indices.Data]

		// Get low row
		lowRow := DatabaseRow{
			IP:         &lowIP,
			Complement: &highIP,
			IsHigh:     false,
			Data:       &data,
		}
		mdb.Rows = append(mdb.Rows, &lowRow)

		// Get high Row
		highRow := DatabaseRow{
			IP:         &highIP,
			Complement: &lowIP,
			IsHigh:     true,
			Data:       &data,
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

	// Run garbage collection
	runtime.GC()

	return mdb
}

// GeoHandle returns ip data
func GeoHandle(ipstr string) []string {
	// Parse the IP to bytes
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return nil
	}

	// Get hexadecimal for lookup
	lookupIP := ip.To16()

	// Lookup all databases
	answer := make([]string, 0)
	for _, dbSet := range databaseSets {
		for _, db := range dbSet.Databases {
			row := db.Lookup(lookupIP)
			if row != nil {
				response := fmt.Sprintf("%s=%s", dbSet.AttributeName, row.getResponse(db))
				answer = append(answer, response)
				break
			}
		}
	}

	return answer
}
