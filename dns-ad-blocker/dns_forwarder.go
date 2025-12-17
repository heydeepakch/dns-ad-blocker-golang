package main

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

// handleDNSRequest handles every incoming DNS query
func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	// r = incoming DNS request
	// w = where we write the response

	// Create a new DNS client to forward request
	client := new(dns.Client)
	client.Net = "udp"

	// Forward the request to Google DNS
	resp, _, err := client.Exchange(r, "8.8.8.8:53")
	if err != nil {
		log.Println("DNS forward error:", err)
		return
	}

	// Write the response back to the original client
	err = w.WriteMsg(resp)
	if err != nil {
		log.Println("DNS write error:", err)
	}
}

func main() {
	// Attach handler to all DNS queries
	dns.HandleFunc(".", handleDNSRequest)

	server := &dns.Server{
		Addr: ":53", // DNS port
		Net:  "udp",
	}

	fmt.Println("DNS forwarder running on port 53")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start DNS server:", err)
	}
}
