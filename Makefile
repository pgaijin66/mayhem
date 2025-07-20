# Makefile for mayhem
.PHONY: build run clean test podman

# Build variables
BINARY_NAME=mayhem
MAIN_PACKAGE=./main.go

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify

## vulncheck: Check for vulnerabilities
.PHONY: vulncheck
vulncheck:
	govulncheck ./...

# Build the binary
build:
	@echo "🔨 Building mayhem..."
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✅ Build complete: $(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "✅ Multi-platform build complete"

# Run with example configuration
run-example:
	@echo "🚀 Starting mayhem with example configuration..."
	./$(BINARY_NAME) -target=http://httpbin.org -port=8080 -delay-prob=0.3 -error-prob=0.1

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	rm -f $(BINARY_NAME)*
	@echo "✅ Clean complete"

# Create example configuration file
config-example:
	@echo "📄 Creating example configuration..."
	@echo '{' > chaos-config.json
	@echo '  "delay_enabled": true,' >> chaos-config.json
	@echo '  "delay_min": "100ms",' >> chaos-config.json
	@echo '  "delay_max": "2s",' >> chaos-config.json
	@echo '  "delay_probability": 0.2,' >> chaos-config.json
	@echo '  "error_enabled": true,' >> chaos-config.json
	@echo '  "error_codes": [500, 502, 503, 504],' >> chaos-config.json
	@echo '  "error_probability": 0.1,' >> chaos-config.json
	@echo '  "error_message": "Chaos engineering fault injection",' >> chaos-config.json
	@echo '  "timeout_enabled": true,' >> chaos-config.json
	@echo '  "timeout_duration": "30s",' >> chaos-config.json
	@echo '  "timeout_probability": 0.05' >> chaos-config.json
	@echo '}' >> chaos-config.json
	@echo "✅ Example configuration created: chaos-config.json"

# Alternative a
# podman build
podman:
	@echo "🐳 Building podman image..."
	podman build -t mayhem:latest .
	@echo "✅ podman image built: mayhem:latest"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed"

# Run with podman
podman-run:
	@echo "🐳 Running mayhem in podman..."
	podman run -p 8080:8080 mayhem:latest -target=http://httpbin.org

# Help
help:
	@echo "mayhem - API Chaos Engineering Tool"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  run-example  - Run with example configuration"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  config-example - Create example configuration file"
	@echo "  podman       - Build podman image"
	@echo "  podman-run   - Run with podman"
	@echo "  deps         - Install dependencies"
	@echo "  help         - Show this help"
