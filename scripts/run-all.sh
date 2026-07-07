#!/usr/bin/env bash
# Starts all four processes: gateway + 3 microservices.
# Each log stream is prefixed with the service name.
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
LOGS="$ROOT/logs"
mkdir -p "$LOGS"

pkill -f "go-eats" 2>/dev/null || true

run() {
  local name=$1; local pkg=$2
  echo "Starting $name..."
  cd "$ROOT"
  go run "./$pkg/..." > "$LOGS/$name.log" 2>&1 &
  echo "$! $name"
}

run "user-svc"        "cmd/user-svc"
run "restaurant-svc"  "cmd/restaurant-svc"
run "order-svc"       "cmd/order-svc"
run "gateway"         "cmd/gateway"

echo ""
echo "All services started. Logs in $LOGS/"
echo "  Gateway        → http://localhost:8080"
echo "  User Service   → http://localhost:8081  | gRPC :9091"
echo "  Restaurant Svc → http://localhost:8082  | gRPC :9092"
echo "  Order Service  → http://localhost:8083"
echo ""
echo "Press Ctrl+C to stop all."
wait
