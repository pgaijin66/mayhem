# ChaosKit - API Chaos Engineering Tool

ChaosKit is a powerful chaos engineering tool for APIs that allows you to inject controlled failures, delays, and timeouts into your HTTP services to test their resilience.

## Features

- **HTTP Reverse Proxy**: Acts as a proxy between clients and your target service
- **Delay Injection**: Add configurable delays to simulate network latency
- **Error Injection**: Return random HTTP error codes with custom messages
- **Timeout Simulation**: Simulate request timeouts
- **Real-time Statistics**: Monitor chaos injection statistics
- **Dynamic Configuration**: Update chaos parameters on the fly via REST API
- **CLI Configuration**: Flexible command-line options
- **JSON Configuration**: Load configuration from JSON files

## Quick Start

### Build

```bash
make build
```

### Basic Usage

```bash
# Start chaos proxy pointing to your api 
./chaoskit -target=<YOUR_API_ENDPOINT> -port=8080

# With custom chaos parameters
./chaoskit -target=<YOUR_API_ENDPOINT> -delay-prob=0.3 -error-prob=0.1 -port=8080
```

### Test the Chaos

```bash
# Make requests through the chaos proxy
curl http://localhost:8080/get

# Check statistics
curl http://localhost:8080/_chaos/stats

# Update configuration
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{"delay_probability": 0.5, "error_probability": 0.2}'
```

## Configuration Options

### Command Line Flags

- `-target`: Target service URL (required)
- `-port`: Port to run the chaos proxy on (default: 8080)
- `-delay-min`: Minimum delay duration (default: 100ms)
- `-delay-max`: Maximum delay duration (default: 2s)
- `-delay-prob`: Probability of delay injection 0.0-1.0 (default: 0.1)
- `-error-prob`: Probability of error injection 0.0-1.0 (default: 0.05)
- `-error-codes`: Comma-separated list of error codes (default: "500,502,503,504")
- `-error-msg`: Error message for injected errors
- `-timeout-dur`: Timeout duration (default: 30s)
- `-timeout-prob`: Probability of timeout injection 0.0-1.0 (default: 0.02)
- `-config`: JSON configuration file path

### Configuration File

Create a `chaos-config.json` file:

```bash
make config-example
```

Then run with:

```bash
./chaoskit -target=<YOUR_API_ENDPOINT> -config=chaos-config.json
```

## Management Endpoints

- `GET /_chaos/stats` - View injection statistics
- `GET /_chaos/config` - View current configuration
- `POST /_chaos/config` - Update configuration
- `GET /_chaos/health` - Health check

## Docker Usage

```bash
# Build Docker image
make docker

# Run with Docker
docker run -p 8080:8080 chaoskit:latest -target=<YOUR_API_ENDPOINT>
```

## Use Cases

1. **API Resilience Testing**: Test how your applications handle API failures
2. **Circuit Breaker Testing**: Verify circuit breaker patterns work correctly
3. **Retry Logic Validation**: Ensure retry mechanisms handle failures properly
4. **Timeout Handling**: Test application behavior under slow network conditions
5. **Load Testing Enhancement**: Add realistic failure scenarios to load tests

## Examples

### Testing a Microservice

```bash
# Start your microservice
# Start chaos proxy
./chaoskit -target=<YOUR_API_ENDPOINT> -port=8080 -delay-prob=0.2 -error-prob=0.1

# Your tests now go through localhost:8080 instead of localhost:3000
```

### CI/CD Integration

```bash
# In your test script
./chaoskit -target=$SERVICE_URL -port=8080 -config=test-chaos.json &
CHAOS_PID=$!

# Run your tests against localhost:8080
npm test

# Cleanup
kill $CHAOS_PID
```

## Building from Source

```bash
git clone <your-repo>
cd chaoskit
make deps
make build
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Add tests for new features
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.