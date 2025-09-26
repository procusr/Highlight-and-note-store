package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var serverStartTime = time.Now()

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

// HealthCheckHandler provides a simple health check endpoint
// The ID changes each time the server restarts, allowing client-side detection
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Status: "ok",
		ID:     fmt.Sprintf("%d", serverStartTime.Unix()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}