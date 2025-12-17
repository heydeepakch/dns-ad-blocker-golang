package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

var blockedDomains map[string]bool
var dnsCache = NewDNSCache()

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		return
	}

	q := r.Question[0]
	domain := q.Name[:len(q.Name)-1] // remove trailing dot

	// 🔴 BLOCK CHECK
	if blockedDomains[domain] {
		msg := new(dns.Msg)
		msg.SetReply(r)

		rr := &dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    300,
			},
			A: net.ParseIP("0.0.0.0"),
		}

		msg.Answer = append(msg.Answer, rr)
		w.WriteMsg(msg)

		log.Println("BLOCKED:", domain)
		return
	}
	// Check cache first
	if cached := dnsCache.Get(domain); cached != nil {
		w.WriteMsg(cached)
		return
	}

	// ✅ Forward allowed queries
	client := new(dns.Client)
	client.Net = "udp"

	resp, _, err := client.Exchange(r, "8.8.8.8:53")
	if err != nil {
		log.Println("Forward error:", err)
		return
	}

	// Cache the response
	dnsCache.Set(domain, resp)

	w.WriteMsg(resp)
}

func main() {
	var err error
	blockedDomains, err = loadBlocklist("blocklist.txt")
	if err != nil {
		log.Fatal("Failed to load blocklist:", err)
	}

	fmt.Println("Loaded", len(blockedDomains), "blocked domains")

	dns.HandleFunc(".", handleDNSRequest)

	server := &dns.Server{
		Addr: ":53",
		Net:  "udp",
	}

	fmt.Println("Ad-blocking DNS server running on port 53")
	log.Fatal(server.ListenAndServe())
}
