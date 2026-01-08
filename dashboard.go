package main

import (
	"encoding/json"
	"net/http"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"total_queries":   totalQueries,
		"blocked_queries": blockedQueries,
	}

	json.NewEncoder(w).Encode(response)
}

func blockedHandler(w http.ResponseWriter, r *http.Request) {
	blockedMu.Lock()
	defer blockedMu.Unlock()

	json.NewEncoder(w).Encode(blockedCount)
}

func startDashboard() {
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/blocked", blockedHandler)

	http.ListenAndServe(":8080", nil)
}
