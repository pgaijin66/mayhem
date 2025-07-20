package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// ChaosConfig represents the configuration for chaos injection
type ChaosConfig struct {
	DelayEnabled       bool          `json:"delay_enabled"`
	DelayMin           time.Duration `json:"delay_min"`
	DelayMax           time.Duration `json:"delay_max"`
	DelayProbability   float64       `json:"delay_probability"`
	ErrorEnabled       bool          `json:"error_enabled"`
	ErrorCodes         []int         `json:"error_codes"`
	ErrorProbability   float64       `json:"error_probability"`
	ErrorMessage       string        `json:"error_message"`
	TimeoutEnabled     bool          `json:"timeout_enabled"`
	TimeoutDuration    time.Duration `json:"timeout_duration"`
	TimeoutProbability float64       `json:"timeout_probability"`
}

// ChaosMiddleware represents the chaos engineering middleware
type ChaosMiddleware struct {
	config     *ChaosConfig
	next       http.Handler
	proxy      *httputil.ReverseProxy
	targetURL  *url.URL
	statsDelay int64
	statsError int64
	statsTotal int64
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
	delay := cm.config.DelayMin + time.Duration(rand.Int63n(int64(cm.config.DelayMax-cm.config.DelayMin)))
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

	time.Sleep(cm.config.TimeoutDuration)

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
		"uptime":           time.Since(startTime).String(),
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

var startTime = time.Now()

func main() {
	var (
		port        = flag.String("port", "8080", "Port to run the chaos proxy on")
		target      = flag.String("target", "", "Target service URL (required)")
		delayMin    = flag.Duration("delay-min", 100*time.Millisecond, "Minimum delay duration")
		delayMax    = flag.Duration("delay-max", 2*time.Second, "Maximum delay duration")
		delayProb   = flag.Float64("delay-prob", 0.1, "Probability of delay injection (0.0-1.0)")
		errorProb   = flag.Float64("error-prob", 0.05, "Probability of error injection (0.0-1.0)")
		errorCodes  = flag.String("error-codes", "500,502,503,504", "Comma-separated list of error codes to inject")
		errorMsg    = flag.String("error-msg", "Chaos engineering fault injection", "Error message for injected errors")
		timeoutDur  = flag.Duration("timeout-dur", 30*time.Second, "Timeout duration")
		timeoutProb = flag.Float64("timeout-prob", 0.02, "Probability of timeout injection (0.0-1.0)")
		configFile  = flag.String("config", "", "JSON configuration file path")
	)
	flag.Parse()

	if *target == "" {
		fmt.Println("❌ Target service URL is required")
		flag.Usage()
		os.Exit(1)
	}

	targetURL, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("❌ Invalid target URL: %v", err)
	}

	// Parse error codes
	var codes []int
	for _, code := range strings.Split(*errorCodes, ",") {
		if c, err := strconv.Atoi(strings.TrimSpace(code)); err == nil {
			codes = append(codes, c)
		}
	}

	// Create default configuration
	config := &ChaosConfig{
		DelayEnabled:       true,
		DelayMin:           *delayMin,
		DelayMax:           *delayMax,
		DelayProbability:   *delayProb,
		ErrorEnabled:       true,
		ErrorCodes:         codes,
		ErrorProbability:   *errorProb,
		ErrorMessage:       *errorMsg,
		TimeoutEnabled:     true,
		TimeoutDuration:    *timeoutDur,
		TimeoutProbability: *timeoutProb,
	}

	// Load configuration from file if provided
	if *configFile != "" {
		if data, err := os.ReadFile(*configFile); err == nil {
			if err := json.Unmarshal(data, config); err != nil {
				log.Printf("⚠️  Failed to parse config file: %v", err)
			} else {
				log.Printf("📄 Loaded configuration from %s", *configFile)
			}
		} else {
			log.Printf("⚠️  Failed to read config file: %v", err)
		}
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	chaosMiddleware := NewChaosMiddleware(config, targetURL)

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: chaosMiddleware,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("🛑 Shutting down chaos proxy...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("❌ Server shutdown error: %v", err)
		}
	}()

	// Print startup information
	fmt.Printf(`
🔥 ChaosKit - API Chaos Engineering Tool
=======================================
📡 Proxy listening on: http://localhost:%s
🎯 Target service: %s
⚡ Delay injection: %.1f%% (%.0fms - %.0fms)
💥 Error injection: %.1f%% (codes: %v)
⏱️  Timeout injection: %.1f%% (%v)

Management endpoints:
📊 Stats: http://localhost:%s/_chaos/stats
⚙️  Config: http://localhost:%s/_chaos/config
❤️  Health: http://localhost:%s/_chaos/health

Press Ctrl+C to stop
`, *port, *target,
		*delayProb*100, delayMin.Seconds()*1000, delayMax.Seconds()*1000,
		*errorProb*100, codes,
		*timeoutProb*100, *timeoutDur,
		*port, *port, *port)

	// Start server
	log.Printf("🚀 Starting chaos proxy on port %s", *port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server failed to start: %v", err)
	}

	log.Println("👋 Chaos proxy stopped")
}
