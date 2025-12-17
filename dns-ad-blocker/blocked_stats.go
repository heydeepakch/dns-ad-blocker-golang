package main

import "sync"

var (
	blockedCount = make(map[string]int)
	blockedMu    sync.Mutex
)

func recordBlocked(domain string) {
	blockedMu.Lock()
	blockedCount[domain]++
	blockedMu.Unlock()
}
