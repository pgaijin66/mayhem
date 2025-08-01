package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pgaijin66/phailure/internal/chaos"
	"github.com/pgaijin66/phailure/internal/server"
	"github.com/pgaijin66/phailure/pkg/usage"
	"github.com/pgaijin66/phailure/pkg/version"
)

func main() {
	flag.Usage = usage.CustomUsage
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
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *showVersion {
		version.Print()
		os.Exit(0)
	}

	if *target == "" {
		fmt.Println("❌ Target service URL is required")
		flag.Usage()
		os.Exit(1)
	}

	targetURL, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("❌ Invalid target URL: %v", err)
	}

	config, err := chaos.NewConfigFromFlags(*delayMin, *delayMax, *delayProb, *errorProb, *errorCodes, *errorMsg, *timeoutDur, *timeoutProb)
	if err != nil {
		log.Fatalf("❌ Invalid configuration: %v", err)
	}

	if *configFile != "" {
		if err := config.LoadFromFile(*configFile); err != nil {
			log.Printf("⚠️  Failed to load config file: %v", err)
		} else {
			log.Printf("📄 Loaded configuration from %s", *configFile)
			log.Printf("Config: DelayMin=%v, DelayMax=%v, TimeoutDuration=%v",
				config.DelayMin.Duration, config.DelayMax.Duration, config.TimeoutDuration.Duration)
		}
	} else {
		log.Printf("🚀 Using command line configuration")
		log.Printf("Config: DelayMin=%v, DelayMax=%v, TimeoutDuration=%v",
			config.DelayMin.Duration, config.DelayMax.Duration, config.TimeoutDuration.Duration)
	}

	srv := server.New(*port, config, targetURL)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("🛑 Shutting down chaos proxy...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("❌ Server shutdown error: %v", err)
		}
	}()

	// Start server (this blocks until shutdown)
	srv.Start()
}
