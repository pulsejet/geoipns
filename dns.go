package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
)

// parseQuery parses and responds to the message
func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeTXT:
			// Debug log
			log.Printf("TXT Query for %s\n", q.Name)

			// Get IP
			replacer := strings.NewReplacer(
				".location.", "",
				"x", ":",
				"z", ".")
			ip := replacer.Replace(q.Name)

			// Send response
			rr, err := dns.NewRR(fmt.Sprintf("%s TXT %s", q.Name, ip))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

// handleDNSRequest handles the connection from client
func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}
