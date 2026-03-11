#!/bin/bash
set -e # Exit immediately if a command fails

# -------------------------------------------------------------------
# Database credentials — set via environment variables for security.
# These are sourced from the Taksa platform deployment config.
# Override by setting these vars before running: export DB_USER=...
# -------------------------------------------------------------------
export DB_HOST="${DB_HOST:-127.0.0.1}"
export DB_PORT="${DB_PORT:-5433}"
export DB_USER="${DB_USER:-taksa}"
export DB_PASSWORD="${DB_PASSWORD:-taksa_123}"
export DB_NAME="${DB_NAME:-taksa}"

echo "🚀 Building Go Kratos Traceability API..."
go build -o traceability ./cmd/traceability/

echo "🛑 Stopping old server instances..."
pkill -x traceability 2>/dev/null || true

echo "🏃 Starting server in background..."
# Substitute environment variables into config before starting
envsubst < ./configs/config.yaml > /tmp/traceability_config_runtime.yaml
./traceability -conf /tmp/traceability_config_runtime.yaml > server.log 2>&1 &

echo "⏳ Waiting for server to start..."
sleep 3

echo "🧪 Generating Bruno test files..."
python3 generate_tests.py

echo "✅ Running API test suite..."
cd tests/api
npx @usebruno/cli run --env "Dev VM"

echo "🎉 All done! Server is still running in the background. Check server.log for details."

echo "🧹 Cleaning up temporary files..."
rm -f traceability
