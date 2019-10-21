package geoip

import "net"

// Config is the configuration format
type Config struct {
	Debug bool `json:"debug"`

	DatabaseSets []DatabaseConfigSet `json:"database_sets"`

	HashMapFile      string   `json:"hashmap_file"`
	HashMapFileField []string `json:"hashmap_file_field"`
	HashMapFileKey   string   `json:"hashmap_file_key"`
}

// DatabaseConfigSet is a set of databases that provide one answer record
type DatabaseConfigSet struct {
	AttributeName string           `json:"attribute_name"`
	Databases     []DatabaseConfig `json:"databases"`
}

// DatabaseConfig is the format of configuration for geoip db
type DatabaseConfig struct {
	File       string               `json:"file"`
	Fields     DatabaseConfigFields `json:"fields"`
	UseHashMap bool                 `json:"use_hashmap"`
}

// DatabaseConfigFields is the format of fields for geoip db
type DatabaseConfigFields struct {
	CIDR   string `json:"cidr"`
	LowIP  string `json:"low_ip"`
	HighIP string `json:"high_ip"`
	Data   string `json:"data"`
}

// DatabaseSet a set of databases of GeoIP
type DatabaseSet struct {
	AttributeName string
	Databases     []*Database
}

// Database a database of GeoIP
type Database struct {
	Rows       []*DatabaseRow
}

// DatabaseRow represents a single row in the databse
type DatabaseRow struct {
	IP         *net.IP
	Complement *net.IP
	IsHigh     bool
	Data       *string
	Parent     *DatabaseRow
}

type dbFieldIndex struct {
	CIDR   int
	LowIP  int
	HighIP int
	Data   int
}
