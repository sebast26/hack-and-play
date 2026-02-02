package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Order struct {
	OrderID      string  `json:"order_id"`
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	Total        float64 `json:"total"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

var orders = map[string]Order{
	"123": {
		OrderID:      "123",
		CustomerID:   "c456",
		CustomerName: "Acme Corp",
		Total:        299.99,
		Currency:     "USD",
		Status:       "completed",
		CreatedAt:    "2024-01-15T10:30:00Z",
	},
	"456": {
		OrderID:      "456",
		CustomerID:   "c789",
		CustomerName: "Widget Inc",
		Total:        149.50,
		Currency:     "USD",
		Status:       "pending",
		CreatedAt:    "2024-01-16T14:20:00Z",
	},
	"789": {
		OrderID:      "789",
		CustomerID:   "c456",
		CustomerName: "Acme Corp",
		Total:        599.00,
		Currency:     "EUR",
		Status:       "shipped",
		CreatedAt:    "2024-01-17T09:15:00Z",
	},
}

func main() {
	http.HandleFunc("/api/v3/orders", handleOrders)
	http.HandleFunc("/api/v3/orders/", handleOrderByID)
	http.HandleFunc("/openapi.yaml", handleOpenAPISpec)

	port := "8081"
	log.Printf("Orders mock service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderList := make([]Order, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orders": orderList,
		"total":  len(orderList),
	})
}

func handleOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path: /api/v3/orders/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/v3/orders/")
	id := strings.TrimSuffix(path, "/")

	order, exists := orders[id]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "not_found",
			"message": "Order not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	spec, err := os.ReadFile("openapi.yaml")
	if err != nil {
		http.Error(w, "Spec not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(spec)
}
