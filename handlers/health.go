package handlers

import (
	"http-server/utils"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, map[string]interface{}{
		"status":  "healthy",
		"server":  "Go-Chi HTTP Server",
		"version": "1.0.0",
	})
}
