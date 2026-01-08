package main

import "sync/atomic"

var (
	totalQueries  uint64
	blockedQueries uint64
)

func incTotal() {
	atomic.AddUint64(&totalQueries, 1)
}

func incBlocked() {
	atomic.AddUint64(&blockedQueries, 1)
}
