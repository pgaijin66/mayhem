FROM golang:1.24-alpine AS build_base
ARG LDFLAGS
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64

# Set the Current Working Directory inside the container
WORKDIR /src

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy source code
COPY . .

# Build the ChaosKit application
# hadolint ignore=DL3059
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "${LDFLAGS}" \
    -a -installsuffix cgo \
    -o ./bin/chaoskit \
    main.go

# Start fresh from a smaller image
FROM alpine:3.15

# Add ca-certificates for HTTPS requests and timezone data
RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 chaoskit && \
    adduser -D -s /bin/sh -u 1000 -G chaoskit chaoskit

WORKDIR /app

# Copy the binary from build stage
COPY --from=build_base /src/bin/chaoskit /app/chaoskit

# Change ownership to non-root user
RUN chown -R chaoskit:chaoskit /app

# Switch to non-root user
USER chaoskit

# Make binary executable
RUN chmod +x chaoskit

# This container exposes port 8080 to the outside world (ChaosKit default)
EXPOSE 8080

# Health check to ensure the chaos proxy is running
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/_chaos/health || exit 1

# Set default environment variables
ENV CHAOSKIT_PORT=8080
ENV CHAOSKIT_LOG_LEVEL=info

# Run the ChaosKit binary
ENTRYPOINT ["./chaoskit"]

# Default command line arguments (can be overridden)
CMD ["--port", "8080"]