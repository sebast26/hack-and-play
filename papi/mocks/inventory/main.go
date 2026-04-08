package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/inventory/check", handleInventoryCheck)

	port := "8083"
	log.Printf("Inventory mock service starting on port %s (always returns 400)", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleInventoryCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "inventory_unavailable",
		"message": "Inventory service rejected the request: insufficient stock levels",
	})
}
