package main

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

// GeoIPNSConfig is the GeoIPNS configuration
type GeoIPNSConfig struct {
	Suffix string
	Server string
}

// GeoIPNS performs a geolocation query
func GeoIPNS(ip string, config *GeoIPNSConfig) map[string]string {
	// Construct and make query
	target := ip + "." + config.Suffix + "."
	server := config.Server
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target, dns.TypeTXT)
	r, _, _ := c.Exchange(&m, server)

	// Parse response
	res := make(map[string]string, len(r.Answer))
	for _, ans := range r.Answer {
		record := ans.(*dns.TXT)
		spl := strings.Split(record.Txt[0], "=")
		res[spl[0]] = spl[1]
	}

	return res
}

func main() {
	config := &GeoIPNSConfig{
		Suffix: "geoipns.iitb.ac.in",
		Server: "10.105.177.129:53",
	}
	fmt.Println(GeoIPNS("3.105.177.8", config))
}
