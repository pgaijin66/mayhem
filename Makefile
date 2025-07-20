# Makefile for mayhem
.PHONY: build run clean test podman

# Build variables
APP_NAME := mayhem
BINARY_NAME := mayhem
MAIN_PACKAGE := ./main.go

# Version information
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
TAG := $(shell git describe --tags --exact-match 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) \
                     -X main.GitCommit=$(COMMIT) \
                     -X main.BuildDate=$(BUILD_DATE) \
                     -X main.Tag=$(TAG)"

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
# ==================================================================================== #
# BUILD
# ==================================================================================== #

# Build for current platform
build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o bin/$(APP_NAME) .
	@echo "✅ Built: bin/$(APP_NAME)"

# Build for all platforms
build-all: clean
	@echo "Building $(APP_NAME) $(VERSION) for all platforms..."
	@mkdir -p dist
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(APP_NAME)-linux-arm64 .
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(APP_NAME)-darwin-arm64 .
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP_NAME)-windows-amd64.exe .
	@echo "✅ Binaries built in dist/ for version $(VERSION)"

# Create release archives
package: build-all
	@echo "Creating release packages for $(VERSION)..."
	@mkdir -p releases
	# Linux
	tar -czf releases/$(APP_NAME)-$(VERSION)-linux-amd64.tar.gz -C dist $(APP_NAME)-linux-amd64
	tar -czf releases/$(APP_NAME)-$(VERSION)-linux-arm64.tar.gz -C dist $(APP_NAME)-linux-arm64
	# macOS
	tar -czf releases/$(APP_NAME)-$(VERSION)-darwin-amd64.tar.gz -C dist $(APP_NAME)-darwin-amd64
	tar -czf releases/$(APP_NAME)-$(VERSION)-darwin-arm64.tar.gz -C dist $(APP_NAME)-darwin-arm64
	# Windows
	zip -j releases/$(APP_NAME)-$(VERSION)-windows-amd64.zip dist/$(APP_NAME)-windows-amd64.exe
	@echo "✅ Release packages created in releases/ for $(VERSION)"

# Generate checksums
checksums: package
	@echo "Generating checksums for $(VERSION)..."
	cd releases && sha256sum *.tar.gz *.zip > $(APP_NAME)-$(VERSION)-checksums.txt
	@echo "✅ Checksums generated: releases/$(APP_NAME)-$(VERSION)-checksums.txt"

# ==================================================================================== #
# RELEASE MANAGEMENT
# ==================================================================================== #

# Tag a new version
tag:
	@if [ -z "$(TAG)" ]; then \
		echo "❌ Usage: make tag TAG=v1.0.0"; \
		exit 1; \
	fi
	@echo "🏷️  Creating tag $(TAG)..."
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)
	@echo "✅ Tag $(TAG) created and pushed"

# Create a GitHub release
release: test checksums
	@echo "🚀 Creating GitHub release $(VERSION)..."
	gh release create $(VERSION) releases/* \
		--title "Mayhem $(VERSION)" \
		--notes "Release $(VERSION) of Mayhem API Chaos Engineering Tool" \
		--latest
	@echo "✅ GitHub release $(VERSION) created"

# Release workflow: tag and release
release-workflow:
	@if [ -z "$(TAG)" ]; then \
		echo "❌ Usage: make release-workflow TAG=v1.0.0"; \
		exit 1; \
	fi
	@echo "🚀 Starting release workflow for $(TAG)..."
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)
	@echo "✅ Tag pushed. GitHub Actions will handle the release build."

# ==================================================================================== #
# TESTING
# ==================================================================================== #

# Test
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run with example configuration
run-example: build
	@echo "🚀 Starting mayhem $(VERSION) with example configuration..."
	./bin/$(BINARY_NAME) -target=http://httpbin.org -port=8080 -delay-prob=0.3 -error-prob=0.1

# ==================================================================================== #
# UTILITIES
# ==================================================================================== #

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	rm -rf bin/ dist/ releases/
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

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed"

# Install locally
install: build
	@echo "Installing $(APP_NAME) $(VERSION) to /usr/local/bin/"
	sudo cp bin/$(APP_NAME) /usr/local/bin/
	@echo "✅ $(APP_NAME) $(VERSION) installed successfully"

# Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Show current git status for release
status:
	@echo "📊 Git Status:"
	@echo "Current branch: $(shell git branch --show-current)"
	@echo "Latest tag: $(shell git describe --tags --abbrev=0 2>/dev/null || echo 'No tags')"
	@echo "Current version: $(VERSION)"
	@echo "Uncommitted changes: $(shell git status --porcelain | wc -l)"
	@echo ""
	@echo "Recent commits:"
	@git log --oneline -5

# ==================================================================================== #
# CONTAINER BUILDS
# ==================================================================================== #

# Docker build
docker:
	@echo "🐳 Building Docker image $(APP_NAME):$(VERSION)..."
	docker build -t $(APP_NAME):$(VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "✅ Docker image built: $(APP_NAME):$(VERSION)"

# Podman build
podman:
	@echo "🐳 Building Podman image $(APP_NAME):$(VERSION)..."
	podman build -t $(APP_NAME):$(VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) .
	podman tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "✅ Podman image built: $(APP_NAME):$(VERSION)"

# Run with Docker
docker-run: docker
	@echo "🐳 Running mayhem $(VERSION) in Docker..."
	docker run -p 8080:8080 $(APP_NAME):$(VERSION) -target=http://httpbin.org

# Run with Podman
podman-run: podman
	@echo "🐳 Running mayhem $(VERSION) in Podman..."
	podman run -p 8080:8080 $(APP_NAME):$(VERSION) -target=http://httpbin.org

# ==================================================================================== #
# HELP
# ==================================================================================== #

# Help
help:
	@echo "🔥 Mayhem $(VERSION) - API Chaos Engineering Tool"
	@echo ""
	@echo "📋 Available targets:"
	@echo ""
	@echo "🔨 Build:"
	@echo "  build         - Build for current platform"
	@echo "  build-all     - Build for all platforms"
	@echo "  package       - Create release packages"
	@echo "  checksums     - Generate checksums"
	@echo ""
	@echo "🚀 Release:"
	@echo "  tag TAG=v1.0.0        - Create and push a git tag"
	@echo "  release               - Create GitHub release (requires gh CLI)"
	@echo "  release-workflow TAG=v1.0.0 - Tag and trigger automated release"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test          - Run tests"
	@echo "  run-example   - Run with example configuration"
	@echo ""
	@echo "🐳 Containers:"
	@echo "  docker        - Build Docker image"
	@echo "  docker-run    - Build and run with Docker"
	@echo "  podman        - Build Podman image"
	@echo "  podman-run    - Build and run with Podman"
	@echo ""
	@echo "🛠️  Utilities:"
	@echo "  install       - Install locally to /usr/local/bin"
	@echo "  clean         - Clean build artifacts"
	@echo "  config-example - Create example configuration file"
	@echo "  deps          - Install dependencies"
	@echo "  version       - Show version information"
	@echo "  status        - Show git status and version info"
	@echo ""
	@echo "✅ Quality Control:"
	@echo "  tidy          - Format code and tidy modules"
	@echo "  audit         - Run quality control checks"
	@echo "  vulncheck     - Check for vulnerabilities"
	@echo ""
	@echo "📖 Examples:"
	@echo "  make build                    # Build current platform"
	@echo "  make tag TAG=v1.0.0          # Create release tag"
	@echo "  make release-workflow TAG=v1.0.0  # Full release process"
	@echo "  make docker-run              # Test with Docker"