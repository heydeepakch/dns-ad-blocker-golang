package main

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

type CacheEntry struct {
	Response *dns.Msg
	Expiry   time.Time
}

type DNSCache struct {
	data map[string]*CacheEntry
	mu   sync.RWMutex
}

func NewDNSCache() *DNSCache {
	return &DNSCache{
		data: make(map[string]*CacheEntry),
	}
}

func (c *DNSCache) Get(domain string) *dns.Msg {
	c.mu.RLock()
	entry, exists := c.data[domain]
	c.mu.RUnlock()

	if !exists || time.Now().After(entry.Expiry) {
		return nil
	}

	return entry.Response.Copy()
}

func (c *DNSCache) Set(domain string, msg *dns.Msg) {
	var ttl uint32 = 60

	if len(msg.Answer) > 0 {
		ttl = msg.Answer[0].Header().Ttl
	}

	c.mu.Lock()
	c.data[domain] = &CacheEntry{
		Response: msg.Copy(),
		Expiry:   time.Now().Add(time.Duration(ttl) * time.Second),
	}
	c.mu.Unlock()
}
