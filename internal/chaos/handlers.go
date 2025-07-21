package chaos

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (cm *ChaosMiddleware) handleManagement(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/_chaos/config":
		cm.handleConfigEndpoint(w, r)
	case "/_chaos/stats":
		cm.handleStatsEndpoint(w, r)
	case "/_chaos/health":
		cm.handleHealthEndpoint(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (cm *ChaosMiddleware) handleConfigEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cm.config)
	case http.MethodPost, http.MethodPut:
		var newConfig ChaosConfig
		if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		cm.config = &newConfig
		log.Printf("🔧 Configuration updated")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (cm *ChaosMiddleware) handleStatsEndpoint(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"total_requests":   cm.statsTotal,
		"delays_injected":  cm.statsDelay,
		"errors_injected":  cm.statsError,
		"delay_percentage": float64(cm.statsDelay) / float64(cm.statsTotal) * 100,
		"error_percentage": float64(cm.statsError) / float64(cm.statsTotal) * 100,
		"uptime":           time.Since(cm.startTime).String(),
		"config":           cm.config,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (cm *ChaosMiddleware) handleHealthEndpoint(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"chaos":     "enabled",
		"timestamp": time.Now().Format(time.RFC3339),
		"target":    cm.targetURL.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
