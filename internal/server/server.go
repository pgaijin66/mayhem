package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/pgaijin66/phailure/internal/chaos"
)

// Server represents the HTTP server
type Server struct {
	port            string
	config          *chaos.ChaosConfig
	targetURL       *url.URL
	chaosMiddleware *chaos.ChaosMiddleware
	httpServer      *http.Server
}

// New creates a new server instance
func New(port string, config *chaos.ChaosConfig, targetURL *url.URL) *Server {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	chaosMiddleware := chaos.NewChaosMiddleware(config, targetURL)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: chaosMiddleware,
	}

	return &Server{
		port:            port,
		config:          config,
		targetURL:       targetURL,
		chaosMiddleware: chaosMiddleware,
		httpServer:      httpServer,
	}
}

// Start starts the HTTP server and prints startup information
func (s *Server) Start() {
	s.printStartupInfo()

	log.Printf("🚀 Starting chaos proxy on port %s", s.port)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server failed to start: %v", err)
	}

	log.Println("👋 Chaos proxy stopped")
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) printStartupInfo() {
	delayMinMs := s.config.DelayMin.Duration.Seconds() * 1000
	delayMaxMs := s.config.DelayMax.Duration.Seconds() * 1000

	fmt.Printf(`

	███╗   ███╗ █████╗ ██╗   ██╗██╗  ██╗███████╗███╗   ███╗
	████╗ ████║██╔══██╗╚██╗ ██╔╝██║  ██║██╔════╝████╗ ████║
	██╔████╔██║███████║ ╚████╔╝ ███████║█████╗  ██╔████╔██║
	██║╚██╔╝██║██╔══██║  ╚██╔╝  ██╔══██║██╔══╝  ██║╚██╔╝██║
	██║ ╚═╝ ██║██║  ██║   ██║   ██║  ██║███████╗██║ ╚═╝ ██║
	╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝

    🔥 API Chaos Engineering Tool 🔥

=======================================
📡 Proxy listening on: http://localhost:%s
🎯 Target service: %s
⚡ Delay injection: %.1f%% (%.0fms - %.0fms)
💥 Error injection: %.1f%% (codes: %v)
⏱️ Timeout injection: %.1f%% (%v)

Management endpoints:
📊 Stats: http://localhost:%s/_chaos/stats
⚙️ Config: http://localhost:%s/_chaos/config
❤️ Health: http://localhost:%s/_chaos/health

Press Ctrl+C to stop
`, s.port, s.targetURL.String(),
		s.config.DelayProbability*100, delayMinMs, delayMaxMs,
		s.config.ErrorProbability*100, s.config.ErrorCodes,
		s.config.TimeoutProbability*100, s.config.TimeoutDuration.Duration,
		s.port, s.port, s.port)
}
