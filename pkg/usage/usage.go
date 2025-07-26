package usage

import (
	"flag"
	"fmt"
	"os"
)

// Custom usage function in man page style
func CustomUsage() {
	fmt.Fprintf(os.Stderr, "phailure(1)                    Chaos Engineering Tool                    phailure(1)\n\n")

	fmt.Fprintf(os.Stderr, "NAME\n")
	fmt.Fprintf(os.Stderr, "       phailure - chaos engineering proxy for testing service resilience\n\n")

	fmt.Fprintf(os.Stderr, "SYNOPSIS\n")
	fmt.Fprintf(os.Stderr, "       phailure -target=URL [OPTIONS]\n\n")

	fmt.Fprintf(os.Stderr, "DESCRIPTION\n")
	fmt.Fprintf(os.Stderr, "       phailure is a chaos engineering tool that acts as a proxy between clients\n")
	fmt.Fprintf(os.Stderr, "       and your target service, injecting controlled failures to test resilience.\n")
	fmt.Fprintf(os.Stderr, "       It can simulate network delays, HTTP errors, and timeouts with configurable\n")
	fmt.Fprintf(os.Stderr, "       probabilities.\n\n")

	fmt.Fprintf(os.Stderr, "EXAMPLES\n")
	fmt.Fprintf(os.Stderr, "       Basic API resilience testing:\n")
	fmt.Fprintf(os.Stderr, "           # Terminal 1: Start your API\n")
	fmt.Fprintf(os.Stderr, "           node server.js\n\n")
	fmt.Fprintf(os.Stderr, "           # Terminal 2: Start phailure\n")
	fmt.Fprintf(os.Stderr, "           phailure -target=http://localhost:3000 -delay-prob=0.3 -error-prob=0.1\n\n")
	fmt.Fprintf(os.Stderr, "           # Terminal 3: Send test requests through phailure\n")
	fmt.Fprintf(os.Stderr, "           curl http://localhost:8080/api/users\n\n")

	fmt.Fprintf(os.Stderr, "       Load testing with chaos:\n")
	fmt.Fprintf(os.Stderr, "           # Start phailure with moderate chaos\n")
	fmt.Fprintf(os.Stderr, "           phailure -target=http://localhost:3000 -delay-prob=0.2 -error-prob=0.05\n\n")
	fmt.Fprintf(os.Stderr, "           # Run load test through phailure\n")
	fmt.Fprintf(os.Stderr, "           hey -n 1000 -c 10 http://localhost:8080/api/endpoint\n\n")

	fmt.Fprintf(os.Stderr, "       Testing specific failure scenarios:\n")
	fmt.Fprintf(os.Stderr, "           # Test timeout handling only\n")
	fmt.Fprintf(os.Stderr, "           phailure -target=http://localhost:3000 \\\n")
	fmt.Fprintf(os.Stderr, "                  -delay-prob=0 -error-prob=0 -timeout-prob=0.5 -timeout-dur=5s\n\n")
	fmt.Fprintf(os.Stderr, "           # Test specific error codes\n")
	fmt.Fprintf(os.Stderr, "           phailure -target=http://localhost:3000 \\\n")
	fmt.Fprintf(os.Stderr, "                  -error-codes=503,504 -error-prob=0.3\n\n")

	fmt.Fprintf(os.Stderr, "       Gradual chaos increase:\n")
	fmt.Fprintf(os.Stderr, "           # Start with low chaos\n")
	fmt.Fprintf(os.Stderr, "           phailure -target=http://localhost:3000 -delay-prob=0.1 -error-prob=0.02\n\n")
	fmt.Fprintf(os.Stderr, "           # Increase error rate during testing\n")
	fmt.Fprintf(os.Stderr, "           curl -X POST http://localhost:8080/_chaos/config \\\n")
	fmt.Fprintf(os.Stderr, "                -H \"Content-Type: application/json\" \\\n")
	fmt.Fprintf(os.Stderr, "                -d '{\"error_probability\": 0.1}'\n\n")

	fmt.Fprintf(os.Stderr, "       Using configuration file:\n")
	fmt.Fprintf(os.Stderr, "           phailure -config=chaos-config.json\n\n")

	fmt.Fprintf(os.Stderr, "OPTIONS\n")
	flag.PrintDefaults()

	fmt.Fprintf(os.Stderr, "\nCHAOS MANAGEMENT ENDPOINTS\n")
	fmt.Fprintf(os.Stderr, "       phailure provides HTTP endpoints for runtime management:\n\n")
	fmt.Fprintf(os.Stderr, "       GET /_chaos/stats\n")
	fmt.Fprintf(os.Stderr, "              Get request statistics and chaos injection counts\n\n")
	fmt.Fprintf(os.Stderr, "       GET /_chaos/config\n")
	fmt.Fprintf(os.Stderr, "              Get current chaos configuration\n\n")
	fmt.Fprintf(os.Stderr, "       POST /_chaos/config\n")
	fmt.Fprintf(os.Stderr, "              Update chaos configuration at runtime\n\n")
	fmt.Fprintf(os.Stderr, "       GET /_chaos/health\n")
	fmt.Fprintf(os.Stderr, "              Health check endpoint\n\n")

	fmt.Fprintf(os.Stderr, "BEST PRACTICES\n")
	fmt.Fprintf(os.Stderr, "       • Start Small: Begin with low probability values (0.01-0.05)\n")
	fmt.Fprintf(os.Stderr, "       • Monitor Everything: Watch your application metrics during testing\n")
	fmt.Fprintf(os.Stderr, "       • Test in Stages: Gradually increase chaos levels\n")
	fmt.Fprintf(os.Stderr, "       • Use Realistic Values: Base probabilities on real-world failure rates\n")
	fmt.Fprintf(os.Stderr, "       • Document Findings: Record how your system responds to chaos\n\n")

	fmt.Fprintf(os.Stderr, "TROUBLESHOOTING\n")
	fmt.Fprintf(os.Stderr, "       phailure won't start:\n")
	fmt.Fprintf(os.Stderr, "       • Ensure the target URL is accessible\n")
	fmt.Fprintf(os.Stderr, "       • Check that the specified port is available\n")
	fmt.Fprintf(os.Stderr, "       • Confirm your target service is running\n\n")
	fmt.Fprintf(os.Stderr, "       No chaos being injected:\n")
	fmt.Fprintf(os.Stderr, "       • Verify probability values are greater than 0\n")
	fmt.Fprintf(os.Stderr, "       • Check chaos injection is enabled in configuration\n")
	fmt.Fprintf(os.Stderr, "       • Use /_chaos/stats to confirm requests are flowing\n\n")
	fmt.Fprintf(os.Stderr, "       Too much chaos:\n")
	fmt.Fprintf(os.Stderr, "       • Reduce probability values\n")
	fmt.Fprintf(os.Stderr, "       • Disable specific chaos types via /_chaos/config\n\n")

	fmt.Fprintf(os.Stderr, "DEBUGGING\n")
	fmt.Fprintf(os.Stderr, "       Use management endpoints to debug issues:\n\n")
	fmt.Fprintf(os.Stderr, "           # Check request flow\n")
	fmt.Fprintf(os.Stderr, "           curl http://localhost:8080/_chaos/stats\n\n")
	fmt.Fprintf(os.Stderr, "           # Verify configuration\n")
	fmt.Fprintf(os.Stderr, "           curl http://localhost:8080/_chaos/config\n\n")
	fmt.Fprintf(os.Stderr, "           # Health check\n")
	fmt.Fprintf(os.Stderr, "           curl http://localhost:8080/_chaos/health\n\n")

	fmt.Fprintf(os.Stderr, "FILES\n")
	fmt.Fprintf(os.Stderr, "       ~/.phailure/config.json\n")
	fmt.Fprintf(os.Stderr, "              Default configuration file location\n\n")

	fmt.Fprintf(os.Stderr, "EXIT STATUS\n")
	fmt.Fprintf(os.Stderr, "       0      Success\n")
	fmt.Fprintf(os.Stderr, "       1      General error\n")
	fmt.Fprintf(os.Stderr, "       2      Invalid arguments\n\n")

	fmt.Fprintf(os.Stderr, "AUTHOR\n")
	fmt.Fprintf(os.Stderr, "       Written by pgaijin66.\n\n")

	fmt.Fprintf(os.Stderr, "REPORTING BUGS\n")
	fmt.Fprintf(os.Stderr, "       Report bugs at: https://github.com/pgaijin66/phailure/issues\n\n")

	fmt.Fprintf(os.Stderr, "SEE ALSO\n")
	fmt.Fprintf(os.Stderr, "       curl(1), hey(1), wrk(1)\n")
	fmt.Fprintf(os.Stderr, "       Online documentation: https://github.com/pgaijin66/phailure\n\n")

	fmt.Fprintf(os.Stderr, "phailure 1.0.0                      2024                           phailure(1)\n")
}
