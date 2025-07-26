## Table Of Contents

- [Start your API](#start-your-api)
- [Test Normal API Functionality](#test-normal-api-functionality)
- [Start phailure Proxy](#start-phailure-proxy)
- [Test API Through Chaos Proxy](#test-api-through-chaos-proxy)
- [Run Multiple Tests to Observe Chaos Effects](#run-multiple-tests-to-observe-chaos-effects)
- [Test All Endpoints Through Chaos](#test-all-endpoints-through-chaos)
- [Monitor Chaos Statistics](#monitor-chaos-statistics)
- [Dynamically Adjust Chaos Levels](#dynamically-adjust-chaos-levels)
- [Reduce Chaos to Minimal](#reduce-chaos-to-minimal)
- [Disable All Chaos](#disable-all-chaos)
- [Specific Chaos Scenarios](#specific-chaos-scenarios)
   * [Test Error Resilience](#test-error-resilience)
   * [Test Timeout Behavior](#test-timeout-behavior)
   * [Test High Latency](#test-high-latency)
   * [Mixed Chaos](#mixed-chaos)
- [Load Testing with Chaos](#load-testing-with-chaos)
- [Advanced Testing Patterns](#advanced-testing-patterns)
   * [Gradual Chaos Increase](#gradual-chaos-increase)
   * [Endpoint-Specific Testing](#endpoint-specific-testing)
- [Cleanup and Reset](#cleanup-and-reset)
- [Quick Reference Commands](#quick-reference-commands)
- [What to Observe](#what-to-observe)


### Start your API

```bash
# Terminal 1
go run main.go
```

or any API. Note down the address

Your API will be running on `http://localhost:9090`

### Test Normal API Functionality

Before introducing chaos, verify your API works correctly:


```bash
# Terminal 2:
# Test hello endpoint
curl http://localhost:9090/hello
# Expected: {"message":"world"}

# Test users endpoint
curl http://localhost:9090/users
# Expected: {"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}

# Test creating a user
curl -X POST http://localhost:9090/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Charlie"}'
# Expected: {"message":"User created successfully"}

# Test with verbose output to see response details
curl -v http://localhost:9090/hello
```

### Start phailure Proxy

Open a new terminal and start phailure pointing to your Gin API:

```bash
# Terminal 3
./phailure \
  --target http://localhost:9090 \
  --port 8080 \
  --delay-prob 0.2 \
  --error-prob 0.1 \
  --timeout-prob 0.05 \
  --delay-min 100ms \
  --delay-max 2s
```

You should see startup output like:


```bash
# Terminal 3
$ ./phailure \
  --target http://localhost:9090 \
  --port 8080 \
  --delay-prob 0.2 \
  --error-prob 0.1 \
  --timeout-prob 0.05 \
  --delay-min 100ms \
  --delay-max 2s

🔥 phailure - API Chaos Engineering Tool
=======================================
📡 Proxy listening on: http://localhost:8080
🎯 Target service: http://localhost:9090
⚡  Delay injection: 20.0% (100ms - 2000ms)
💥 Error injection: 10.0% (codes: [500 502 503 504])
⏱️  Timeout injection: 5.0% (30s)

Management endpoints:
📊 Stats: http://localhost:8080/_chaos/stats
⚙️ Config: http://localhost:8080/_chaos/config
❤️ Health: http://localhost:8080/_chaos/health

Press Ctrl+C to stop
2025/07/20 13:39:40 🚀 Starting chaos proxy on port 8080
```

### Test API Through Chaos Proxy

Now test your API through the chaos proxy on port 8080:

```
# Test hello endpoint through chaos proxy
curl http://localhost:8080/hello

# You might get different responses:
# ✅ Normal response: {"message":"world"}
# 💥 Chaos error: {"error":"Chaos engineering fault injection","code":500,"chaos":true}
# ⏰ Delayed response (takes longer than usual)
# ⏱️ Timeout response: {"error":"Request timeout due to chaos engineering","code":504}
```

### Run Multiple Tests to Observe Chaos Effects

Run the same request multiple times to see different chaos scenarios:

```bash
# Run 10 requests to see chaos in action
echo "Testing chaos effects with 10 requests:"
for i in {1..10}; do
  echo "Request $i:"
  curl -w "HTTP Status: %{http_code} | Time: %{time_total}s\n" \
       -s http://localhost:8080/hello
  echo "---"
  sleep 1
done
```

### Test All Endpoints Through Chaos

Test each of your API endpoints through the chaos proxy:

```bash
# Test users endpoint
echo "Testing /users endpoint:"
curl -w "Status: %{http_code} | Time: %{time_total}s\n" \
     http://localhost:8080/users

# Test POST request
echo "Testing POST /users:"
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"TestUser"}' \
  -w "Status: %{http_code} | Time: %{time_total}s\n"

# Test with detailed timing information
echo "Testing with detailed timing:"
curl -w "Response Time: %{time_total}s | DNS: %{time_namelookup}s | Connect: %{time_connect}s | Transfer: %{time_starttransfer}s\n" \
     -s http://localhost:8080/hello
```

### Monitor Chaos Statistics

phailure provides management endpoints to monitor what's happening:

```bash
# Check chaos statistics
echo "Current chaos statistics:"
curl -s http://localhost:8080/_chaos/stats | python3 -m json.tool

# Example output:
# {
#   "total_requests": 25,
#   "delays_injected": 5,
#   "errors_injected": 3,
#   "delay_percentage": 20.0,
#   "error_percentage": 12.0,
#   "uptime": "5m30s",
#   "config": {...}
# }

# Check current configuration
echo "Current chaos configuration:"
curl -s http://localhost:8080/_chaos/config | python3 -m json.tool

# Check health of chaos proxy
echo "Chaos proxy health:"
curl -s http://localhost:8080/_chaos/health | python3 -m json.tool
```

### Dynamically Adjust Chaos Levels

You can update chaos configuration in real-time without restarting:

```bash
echo "Setting HIGH chaos level..."
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "delay_enabled": true,
    "delay_min": "500ms",
    "delay_max": "3s",
    "delay_probability": 0.5,
    "error_enabled": true,
    "error_codes": [500, 502, 503, 504, 429],
    "error_probability": 0.3,
    "error_message": "High chaos mode activated!",
    "timeout_enabled": true,
    "timeout_duration": "10s",
    "timeout_probability": 0.1
  }'

# Test with high chaos
echo "Testing with HIGH chaos:"
for i in {1..5}; do
  curl -w "Request $i: Status %{http_code} | Time %{time_total}s\n" \
       -s http://localhost:8080/hello
done
```

### Reduce Chaos to Minimal

```bash
echo "Setting LOW chaos level..."
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "delay_enabled": true,
    "delay_min": "10ms",
    "delay_max": "100ms", 
    "delay_probability": 0.05,
    "error_enabled": true,
    "error_codes": [500],
    "error_probability": 0.02,
    "error_message": "Minimal chaos mode",
    "timeout_enabled": false,
    "timeout_probability": 0.0
  }'

# Test with low chaos
echo "Testing with LOW chaos:"
for i in {1..5}; do
  curl -w "Request $i: Status %{http_code} | Time %{time_total}s\n" \
       -s http://localhost:8080/hello
done
```

### Disable All Chaos

```bash
echo "Disabling ALL chaos..."
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "delay_enabled": false,
    "error_enabled": false,
    "timeout_enabled": false
  }'

# Test without chaos (should behave normally)
echo "Testing with NO chaos:"
for i in {1..3}; do
  curl -w "Request $i: Status %{http_code} | Time %{time_total}s\n" \
       -s http://localhost:8080/hello
done
```

### Specific Chaos Scenarios

#### Test Error Resilience

```bash
echo "=== SCENARIO 1: High Error Rate ==="
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "error_probability": 0.8,
    "error_enabled": true,
    "error_codes": [500, 502, 503, 504],
    "delay_enabled": false,
    "timeout_enabled": false
  }'

echo "Testing with 80% error rate:"
for i in {1..10}; do
  response=$(curl -s -w "%{http_code}" http://localhost:8080/hello)
  status_code=$(echo "$response" | tail -c 4)
  echo "Request $i: HTTP $status_code"
done
```

#### Test Timeout Behavior

```bash
echo "=== SCENARIO 2: Forced Timeouts ==="
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "timeout_enabled": true,
    "timeout_probability": 1.0,
    "timeout_duration": "5s",
    "error_enabled": false,
    "delay_enabled": false
  }'

echo "Testing with guaranteed timeouts (5s):"
for i in {1..3}; do
  echo "Request $i (this will take 5 seconds):"
  time curl --max-time 10 -s http://localhost:8080/hello
  echo ""
done
```

#### Test High Latency

```bash
echo "=== SCENARIO 3: High Latency ==="
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "delay_enabled": true,
    "delay_probability": 1.0,
    "delay_min": "1s",
    "delay_max": "3s",
    "error_enabled": false,
    "timeout_enabled": false
  }'

echo "Testing with guaranteed delays (1-3 seconds):"
for i in {1..5}; do
  echo "Request $i:"
  time curl -s http://localhost:8080/hello | head -c 50
  echo ""
done
```

#### Mixed Chaos

```bash
echo "=== SCENARIO 4: Realistic Mixed Chaos ==="
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "delay_enabled": true,
    "delay_probability": 0.3,
    "delay_min": "200ms",
    "delay_max": "1s",
    "error_enabled": true,
    "error_probability": 0.15,
    "error_codes": [500, 503, 504],
    "timeout_enabled": true,
    "timeout_probability": 0.05,
    "timeout_duration": "8s"
  }'

echo "Testing with realistic mixed chaos:"
for i in {1..15}; do
  start_time=$(date +%s.%N)
  response=$(curl -s -w "%{http_code}" --max-time 10 http://localhost:8080/hello 2>/dev/null || echo "TIMEOUT")
  end_time=$(date +%s.%N)
  duration=$(echo "$end_time - $start_time" | bc)
  
  if [[ "$response" == *"TIMEOUT"* ]]; then
    echo "Request $i: TIMEOUT after ${duration}s"
  else
    status_code=$(echo "$response" | tail -c 4)
    echo "Request $i: HTTP $status_code in ${duration}s"
  fi
done
```

### Load Testing with Chaos

Run a simple concurrent load test while chaos is active:

```
echo "=== LOAD TEST WITH CHAOS ==="

# Set moderate chaos
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "delay_probability": 0.2,
    "error_probability": 0.1,
    "timeout_probability": 0.03,
    "delay_enabled": true,
    "error_enabled": true,
    "timeout_enabled": true
  }'

echo "Running 50 concurrent requests with chaos..."

# Create a function to run a single test
run_test() {
  local id=$1
  local result=$(curl -s -w "%{http_code}:%{time_total}" --max-time 5 \
                 http://localhost:8080/hello 2>/dev/null || echo "TIMEOUT:5.000")
  echo "Request $id: $result"
}

# Export the function so subshells can use it
export -f run_test

# Run 50 concurrent requests
for i in {1..50}; do
  run_test $i &
done

# Wait for all background jobs to complete
wait

echo "Load test completed!"

# Check final statistics
echo "Final chaos statistics:"
curl -s http://localhost:8080/_chaos/stats | python3 -m json.tool
```

### Advanced Testing Patterns

#### Gradual Chaos Increase

```bash
echo "=== GRADUAL CHAOS INCREASE ==="

chaos_levels=(0.0 0.1 0.2 0.5 0.8)

for level in "${chaos_levels[@]}"; do
  echo "Setting chaos level to $level..."
  curl -s -X POST http://localhost:8080/_chaos/config \
    -H "Content-Type: application/json" \
    -d "{
      \"delay_probability\": $level,
      \"error_probability\": $(echo "$level / 2" | bc -l),
      \"delay_enabled\": true,
      \"error_enabled\": true
    }" > /dev/null
  
  echo "Testing with chaos level $level:"
  success=0
  total=10
  
  for i in $(seq 1 $total); do
    if curl -s --max-time 3 http://localhost:8080/hello > /dev/null 2>&1; then
      ((success++))
    fi
  done
  
  success_rate=$(echo "scale=1; $success * 100 / $total" | bc)
  echo "Success rate: $success_rate% ($success/$total)"
  echo ""
done
```

#### Endpoint-Specific Testing

```bash
echo "=== ENDPOINT-SPECIFIC TESTING ==="

endpoints=("/hello" "/users")

# Set moderate chaos
curl -s -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "delay_probability": 0.3,
    "error_probability": 0.2,
    "delay_enabled": true,
    "error_enabled": true
  }' > /dev/null

for endpoint in "${endpoints[@]}"; do
  echo "Testing endpoint: $endpoint"
  
  success=0
  errors=0
  timeouts=0
  total=10
  
  for i in $(seq 1 $total); do
    response=$(curl -s -w "%{http_code}" --max-time 3 \
               http://localhost:8080$endpoint 2>/dev/null || echo "TIMEOUT")
    
    if [[ "$response" == *"TIMEOUT"* ]]; then
      ((timeouts++))
    elif [[ "$response" == *"200"* ]] || [[ "$response" == *"201"* ]]; then
      ((success++))
    else
      ((errors++))
    fi
  done
  
  echo "  Success: $success/$total"
  echo "  Errors:  $errors/$total"
  echo "  Timeouts: $timeouts/$total"
  echo ""
done
```

### Cleanup and Reset

When you're done testing:

```bash
# Reset chaos to minimal
echo "Resetting chaos to minimal levels..."
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{
    "delay_enabled": false,
    "error_enabled": false,
    "timeout_enabled": false
  }'

# Final statistics
echo "Final chaos proxy statistics:"
curl -s http://localhost:8080/_chaos/stats | python3 -m json.tool

# Stop phailure (Ctrl+C in the terminal running it)
# Stop your Gin API (Ctrl+C in the terminal running it)
```

### Quick Reference Commands

```bash
# Start your API
go run main.go

# Start phailure with moderate chaos
./phailure --target http://localhost:9090 --port 8080 --delay-prob 0.2 --error-prob 0.1

# Test normal endpoint
curl http://localhost:9090/hello

# Test through chaos proxy
curl http://localhost:8080/hello

# Check chaos stats
curl -s http://localhost:8080/_chaos/stats | python3 -m json.tool

# Update chaos config to high chaos
curl -X POST http://localhost:8080/_chaos/config \
  -H "Content-Type: application/json" \
  -d '{"error_probability": 0.5, "delay_probability": 0.4}'

# Run 10 test requests
for i in {1..10}; do 
  curl -w "Status: %{http_code} | Time: %{time_total}s\n" \
       -s http://localhost:8080/hello
done

# Test with timeout handling
curl --max-time 5 http://localhost:8080/hello
```

### What to Observe
When testing with chaos, look for:

- Some requests fail with 5xx status codes

- Requests take different amounts of time

- Some requests may timeout completely

- Overall percentage of successful requests

- Which types of errors occur most frequently

- How your application handles failures

