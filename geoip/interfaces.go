package geoip

import "net"

// Config is the configuration format
type Config struct {
	Debug bool `json:"debug"`

	Databases [][]DatabaseConfig `json:"databases"`

	LocationFile      string   `json:"location_file"`
	LocationFileField []string `json:"location_file_field"`
	LocationFileKey   string   `json:"location_file_key"`
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

// Database a database of GeoIP
type Database struct {
	Rows      []*DatabaseRow
	UseLocMap bool
}

// DatabaseRow represents a single row in the databse
type DatabaseRow struct {
	IP         *net.IP
	Complement *net.IP
	IsHigh     bool
	Location   *string
	Parent     *DatabaseRow
}

type dbFieldIndex struct {
	CIDR     int
	LowIP    int
	HighIP   int
	Location int
}
