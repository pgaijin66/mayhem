package chaos

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// ChaosMiddleware represents the chaos engineering middleware
type ChaosMiddleware struct {
	config     *ChaosConfig
	next       http.Handler
	proxy      *httputil.ReverseProxy
	targetURL  *url.URL
	statsDelay int64
	statsError int64
	statsTotal int64
	startTime  time.Time
}

// NewChaosMiddleware creates a new chaos middleware
func NewChaosMiddleware(config *ChaosConfig, targetURL *url.URL) *ChaosMiddleware {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Customize the proxy to handle chaos injection
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
	}

	return &ChaosMiddleware{
		config:    config,
		proxy:     proxy,
		targetURL: targetURL,
		startTime: time.Now(),
	}
}

// ServeHTTP implements the http.Handler interface
func (cm *ChaosMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cm.statsTotal++

	// Skip chaos for management endpoints
	if strings.HasPrefix(r.URL.Path, "/_chaos") {
		cm.handleManagement(w, r)
		return
	}

	// Apply chaos effects
	if cm.shouldApplyChaos() {
		if cm.config.TimeoutEnabled && cm.shouldApplyTimeout() {
			cm.applyTimeout(w, r)
			return
		}

		if cm.config.DelayEnabled && cm.shouldApplyDelay() {
			cm.applyDelay()
		}

		if cm.config.ErrorEnabled && cm.shouldApplyError() {
			cm.applyError(w, r)
			return
		}
	}

	// Add chaos headers
	w.Header().Set("X-Chaos-Applied", "true")
	w.Header().Set("X-Chaos-Timestamp", time.Now().Format(time.RFC3339))

	// Forward to target service
	cm.proxy.ServeHTTP(w, r)
}

func (cm *ChaosMiddleware) shouldApplyChaos() bool {
	return true // Always consider chaos, individual methods check probabilities
}

func (cm *ChaosMiddleware) shouldApplyDelay() bool {
	return rand.Float64() < cm.config.DelayProbability
}

func (cm *ChaosMiddleware) shouldApplyError() bool {
	return rand.Float64() < cm.config.ErrorProbability
}

func (cm *ChaosMiddleware) shouldApplyTimeout() bool {
	return rand.Float64() < cm.config.TimeoutProbability
}

func (cm *ChaosMiddleware) applyDelay() {
	// Access the underlying time.Duration from the Duration wrapper
	minDelay := cm.config.DelayMin.Duration
	maxDelay := cm.config.DelayMax.Duration

	// Calculate random delay between min and max
	delayRange := maxDelay - minDelay
	delay := minDelay + time.Duration(rand.Int63n(int64(delayRange)))

	cm.statsDelay++
	log.Printf("💥 Injecting delay: %v", delay)
	time.Sleep(delay)
}

func (cm *ChaosMiddleware) applyError(w http.ResponseWriter, r *http.Request) {
	statusCode := cm.config.ErrorCodes[rand.Intn(len(cm.config.ErrorCodes))]
	cm.statsError++

	log.Printf("💥 Injecting error: HTTP %d", statusCode)

	w.Header().Set("X-Chaos-Injected-Error", fmt.Sprintf("%d", statusCode))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]interface{}{
		"error":     cm.config.ErrorMessage,
		"code":      statusCode,
		"chaos":     true,
		"timestamp": time.Now().Format(time.RFC3339),
		"path":      r.URL.Path,
	}

	json.NewEncoder(w).Encode(errorResponse)
}

func (cm *ChaosMiddleware) applyTimeout(w http.ResponseWriter, r *http.Request) {
	log.Printf("💥 Injecting timeout: %v", cm.config.TimeoutDuration)

	time.Sleep(cm.config.TimeoutDuration.Duration)

	w.Header().Set("X-Chaos-Injected-Timeout", cm.config.TimeoutDuration.String())
	w.WriteHeader(http.StatusGatewayTimeout)

	errorResponse := map[string]interface{}{
		"error":     "Request timeout due to chaos engineering",
		"code":      504,
		"chaos":     true,
		"timeout":   cm.config.TimeoutDuration.String(),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(errorResponse)
}
