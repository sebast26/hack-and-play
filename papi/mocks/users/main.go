package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type User struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	CreatedAt string `json:"created_at"`
}

var users = map[string]User{
	"456": {
		UserID:    "456",
		Email:     "partner@acme.com",
		Name:      "John Doe",
		Company:   "Acme Corp",
		CreatedAt: "2024-01-10T08:00:00Z",
	},
	"789": {
		UserID:    "789",
		Email:     "admin@widget.io",
		Name:      "Jane Smith",
		Company:   "Widget Inc",
		CreatedAt: "2024-01-12T11:30:00Z",
	},
	"101": {
		UserID:    "101",
		Email:     "dev@example.com",
		Name:      "Bob Johnson",
		Company:   "Example LLC",
		CreatedAt: "2024-01-14T16:45:00Z",
	},
}

func main() {
	http.HandleFunc("/api/v1/users", handleUsers)
	http.HandleFunc("/api/v1/users/", handleUserByID)
	http.HandleFunc("/openapi.yaml", handleOpenAPISpec)

	port := "8082"
	log.Printf("Users mock service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": userList,
		"total": len(userList),
	})
}

func handleUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path: /api/v1/users/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	id := strings.TrimSuffix(path, "/")

	user, exists := users[id]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "not_found",
			"message": "User not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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
