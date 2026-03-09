# Taksa App - Traceability Service

An enterprise-grade manufacturing traceability and genealogy service built with Go and the Kratos framework. This service implements ISA-95 physical equipment hierarchies and utilizes Directed Acyclic Graphs (DAGs) backed by PostgreSQL Recursive CTEs to provide complete forward and backward lot genealogy. Equipment telemetry and process history are managed via TimescaleDB.

## Features

* **Physical Factory Hierarchy (ISA-95):** Manage Enterprises, Sites, Areas, Production Lines, and Production Units.
* **Genealogy & Traceability Engine:**
* **Backward Trace:** Recursive root-cause analysis from finished goods to raw materials.
* **Forward Trace:** Recursive impact analysis from raw materials to finished goods for targeted recalls.
* **Full Genealogy:** Bidirectional node and edge graph generation for UI visualization.


* **Time-Series Telemetry:** Correlates exact machine conditions (e.g., temperatures, RPM) with specific material lot processing windows using TimescaleDB.
* **Automated Integration Testing:** Comprehensive 20-endpoint end-to-end test suite powered by Bruno CLI and Python-based dynamic test generation.
* **Clean Architecture:** Strict separation of transport (API), service, business logic (Biz), and data layers.

## Prerequisites

* Go 1.24+
* PostgreSQL 12+ with TimescaleDB extension
* Protoc compiler (for gRPC/HTTP code generation)
* Make
* Python 3.x (for test generation)
* Node.js & npm (for the `@usebruno/cli` test runner)

## Installation & Setup

### 1. Install Dependencies

```bash
# Download Go dependencies
make init

# Install Kratos and Wire
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
go install github.com/google/wire/cmd/wire@latest

```

### 2. Database Setup

The service includes a consolidated schema and a seed file that models an EV Battery Pack Manufacturing plant for immediate testing.

```bash
# Create database
psql -U postgres -c "CREATE DATABASE traceability_db;"

# Apply master schema (includes TimescaleDB hypertables)
psql -U postgres -d traceability_db -f database/schema/schema.sql

# Load the factory simulation seed data
psql -U postgres -d traceability_db -f database/schema/seed.sql

```

### 3. Configuration

The service uses `configs/config.yaml`. Database credentials can be overridden via environment variables:

```bash
export DB_SOURCE="postgres://user:password@localhost:5432/traceability_db?sslmode=disable"

```

## Building

```bash
# Generate proto files
make api

# Generate Wire dependency injection
cd cmd/traceability && wire && cd ../..

# Build the service binary
go build -o traceability ./cmd/traceability

```

## Running

### Standard Execution

```bash
./traceability -conf ./configs/config.yaml

```

The service will start listening on:

* **HTTP**: `http://localhost:8000`
* **gRPC**: `localhost:9000`

### Docker Deployment

```bash
docker build -t taksa/traceability:latest .

docker run -d \
  --name traceability \
  -p 8000:8000 \
  -p 9000:9000 \
  -v $(pwd)/configs:/data/conf \
  -e DB_SOURCE="postgres://user:password@db:5432/traceability_db" \
  taksa/traceability:latest

```

## API Endpoints

The API provides 20 endpoints categorized into three core domains:

### 1. Master Data & Hierarchy

CRUD operations for the physical factory model and definitions:

* `/api/v1/traceability/enterprises`
* `/api/v1/traceability/sites`
* `/api/v1/traceability/areas`
* `/api/v1/traceability/lines`
* `/api/v1/traceability/production-units`
* `/api/v1/traceability/equipment-classes`
* `/api/v1/traceability/equipment`
* `/api/v1/traceability/operators`
* `/api/v1/traceability/material-definitions`

### 2. Execution & Production

Records the actual manufacturing processes:

* `/api/v1/traceability/lots` (Material batch registration)
* `/api/v1/traceability/work-orders` (Execution planning and actuals)
* `/api/v1/traceability/genealogy` (Links consumed parent lots to produced child lots)

### 3. Traceability Queries

Advanced recursive data retrieval:

* `GET /api/v1/traceability/trace/backward/{lot_id}`
* `GET /api/v1/traceability/trace/forward/{lot_id}`
* `GET /api/v1/traceability/trace/full/{lot_id}`
* `GET /api/v1/traceability/trace/equipment-history/{lot_id}`

## Testing

The project includes an automated, sequential integration test suite using the Bruno REST client.

### Running the Automated Suite

A shell script is provided to compile the application, start a background instance, generate the dynamic Bruno test files, and execute the test runner.

```bash
chmod +x run_tests.sh

./run_tests.sh

```

### Manual Test Generation

If you wish to generate the test files without running the full suite:

```bash
python3 generate_tests.py
cd tests/api
npx @usebruno/cli run --env "Dev VM"

```

## Project Structure

```text
├── api/                  # Protocol Buffer definitions (gRPC/HTTP contracts)
│   └── traceability/v1/  
├── cmd/                  # Application entry points and Wire injection
│   └── traceability/     
├── configs/              # Environment configurations (YAML)
├── database/             
│   ├── queries/          # Reference SQL queries
│   └── schema/           # DDL schema and simulation seed data
├── internal/             
│   ├── biz/              # Business logic, domain models, and interfaces
│   ├── data/             # PostgreSQL/GORM data access layer
│   ├── service/          # API handlers and DTO mapping
│   ├── conf/             # Configuration structs
│   └── server/           # HTTP/gRPC server initialization
├── tests/                
│   └── api/              # Auto-generated Bruno REST API tests
├── generate_tests.py     # Python script for dynamic Bruno test scaffolding
├── run_tests.sh          # CI/CD automation script for testing
└── Dockerfile            # Multi-stage Docker build

```

## License

Copyright (c) 2026 The Taksa Project Authors. All rights reserved.