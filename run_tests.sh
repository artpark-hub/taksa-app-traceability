#!/bin/bash
set -e # Exit immediately if a command fails

echo "🚀 Building Go Kratos Traceability API..."
go build -o traceability ./cmd/traceability/

echo "🛑 Stopping old server instances..."
pkill traceability 2>/dev/null || true

echo "🏃 Starting server in background..."
./traceability -conf ./configs/config.yaml > server.log 2>&1 &

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
