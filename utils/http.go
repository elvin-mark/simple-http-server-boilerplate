package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func WriteJSONStatus(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func PrintStartupInfo() {
	fmt.Println("\n📚 Endpoints:")
	fmt.Println("  GET  /               - Home page")
	fmt.Println("  GET  /health         - Health check")
	fmt.Println("  GET  /users          - List users")
	fmt.Println("  POST /users          - Create user")
	fmt.Println("  GET  /users/{id}     - Get user")
	fmt.Println("  DELETE /users/{id}   - Delete user")
}
