package main

import (
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
				r := new(dns.TXT)
				r.Hdr = dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    1,
				}
				r.Txt = []string{response}
				m.Answer = append(m.Answer, r)
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
