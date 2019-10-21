package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
	g "github.com/pulsejet/geoipns/geoip"
)

// parseQuery parses and responds to the message
func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeTXT:
			// Debug log
			if mConfig.Debug {
				log.Printf("TXT Query for %s\n", q.Name)
			}

			// Get IP
			replacer := strings.NewReplacer(
				"."+mConfig.Suffix+".", "",
				"x", ":",
				"z", ".")
			ip := replacer.Replace(q.Name)

			// Send response
			for _, response := range g.GeoHandle(ip) {
				rr, err := dns.NewRR(fmt.Sprintf("%s 1 TXT \"%s\"", q.Name, response))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

// handleDNSRequest handles the connection from client
func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	defer w.Close()

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}
