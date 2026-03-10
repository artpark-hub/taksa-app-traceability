# Taksa App - Traceability Service

An enterprise-grade manufacturing traceability and genealogy service built with Go and the Kratos framework. This service implements ISA-95 physical equipment hierarchies and utilizes Directed Acyclic Graphs (DAGs) backed by PostgreSQL Recursive CTEs to provide complete forward and backward lot genealogy. Equipment telemetry and process history are managed via TimescaleDB.

> **Version:** 0.0.x | **Status:** Feature-complete for S006 (Product Traceability) and S014 (Historical Analytics)

### **Features**

* **Physical Factory Hierarchy (ISA-95):** Manage Enterprises, Sites, Areas, Production Lines, Production Units, and Equipment with full CRUD operations.
* **Equipment Registry & Capabilities:** Register machines with serial numbers, classify by type, track operational status, record technical capabilities and runtime properties.
* **Material Management:** Catalog system for raw materials, work-in-progress, finished goods, and packaging with lot-level batch tracking.
* **Genealogy & Traceability Engine:**
  * **Backward Trace:** Recursive root-cause analysis from finished goods to raw materials (recursive CTE).
  * **Forward Trace:** Recursive impact analysis from raw materials to finished goods (critical for product recalls).
  * **Full Genealogy Tree:** Bidirectional node and edge graph generation for UI visualization.
  * **Equipment Process History:** Correlate exact machine telemetry (temperature, RPM, etc.) with specific material lot processing windows using TimescaleDB.
* **Time-Series Telemetry:** TimescaleDB hypertables for sensor data with intelligent window-based aggregation per processing step.
* **Work Order Management:** Complete production execution records linking machines, operators, input lots, and output lots.
* **Operator Tracking:** Link every production step to responsible operator and shift for accountability and investigation.
* **Automated Integration Testing:** Comprehensive 24-request end-to-end test suite with 47 assertions powered by Bruno CLI and Python-based dynamic test generation.
* **Clean Architecture:** Strict separation of transport (API), service, business logic (Biz), and data layers per Kratos framework.

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
psql -U postgres -d traceability_db -f database/schema/002_traceability_features.sql


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

## API Endpoints (30 Total)

### Master Data & Hierarchy (9 endpoints)

CRUD operations for the physical factory model and definitions:

* `POST /api/v1/traceability/enterprises` — Create enterprise
* `GET /api/v1/traceability/enterprises` — List enterprises
* `PATCH /api/v1/traceability/enterprises/{id}` — Update enterprise
* `DELETE /api/v1/traceability/enterprises/{id}` — Delete enterprise
* `POST /api/v1/traceability/sites` — Create site
* `GET /api/v1/traceability/sites?enterprise_id=X` — List sites
* `PATCH /api/v1/traceability/sites/{id}` — Update site
* `DELETE /api/v1/traceability/sites/{id}` — Delete site
* `POST /api/v1/traceability/areas` — Create area
* `GET /api/v1/traceability/areas?site_id=X` — List areas
* `POST /api/v1/traceability/lines` — Create production line
* `GET /api/v1/traceability/lines?area_id=X` — List lines
* `POST /api/v1/traceability/production-units` — Create production unit
* `GET /api/v1/traceability/production-units?line_id=X` — List units
* `POST /api/v1/traceability/equipment-classes` — Register equipment class
* `GET /api/v1/traceability/equipment-classes` — List classes
* `POST /api/v1/traceability/equipment` — Register equipment
* `GET /api/v1/traceability/equipment?line_id=X` — List equipment
* `PATCH /api/v1/traceability/equipment/{id}` — Update equipment

### Execution & Production (8 endpoints)

Records the actual manufacturing processes:

* `POST /api/v1/traceability/material-definitions` — Add material definition
* `GET /api/v1/traceability/material-definitions?material_type=X` — List definitions
* `POST /api/v1/traceability/lots` — Register material lot
* `GET /api/v1/traceability/lots/{lot_id}` — Get lot details
* `GET /api/v1/traceability/lots?status=X&material_type=Y` — List lots (with filters)
* `PATCH /api/v1/traceability/lots/{lot_id}/status` — Update lot status
* `POST /api/v1/traceability/operators` — Register operator
* `GET /api/v1/traceability/operators?shift=X&status=Y` — List operators
* `POST /api/v1/traceability/work-orders` — Create work order
* `GET /api/v1/traceability/work-orders/{id}` — Get work order detail (with input lots)
* `GET /api/v1/traceability/work-orders?status=X&equipment_id=Y` — List work orders
* `PATCH /api/v1/traceability/work-orders/{id}/status` — Update work order status
* `POST /api/v1/traceability/genealogy` — Register parent→child lot link

### Traceability Queries — S006 (4 endpoints)

Advanced recursive genealogy data retrieval:

* `GET /api/v1/traceability/trace/backward/{lot_id}` — Backward trace (all ancestors with depth, machine, operator, timestamps)
* `GET /api/v1/traceability/trace/forward/{lot_id}` — Forward trace (all descendants, critical for product recall)
* `GET /api/v1/traceability/trace/full/{lot_id}` — Full genealogy tree (nodes + edges graph for visualization)
* `GET /api/v1/traceability/trace/equipment-history/{lot_id}` — Equipment process history (sensor telemetry aggregated per processing window)

### Historical Analytics — S014 (4 endpoints)

Machine KPIs, comparison, trends, and factory dashboard:

* `GET /api/v1/traceability/analytics/machine-performance?equipment_id=X&from_time=T1&to_time=T2` — Single machine KPIs
* `POST /api/v1/traceability/analytics/machine-comparison` — Compare multiple machines side-by-side
* `GET /api/v1/traceability/analytics/production-trends?from_time=T1&to_time=T2&granularity=daily|weekly|monthly&equipment_id=X` — Time-series output trends
* `GET /api/v1/traceability/analytics/dashboard?from_time=T1&to_time=T2` — Factory dashboard summary

## Documentation

Complete documentation is provided in the `docs/` directory:

| File | Purpose |
|---|---|
| **docs/01_database.md** | Database schema reference — all 16 tables with columns, types, FKs, indexes, and real factory examples |
| **docs/02_api_reference.md** | API endpoint reference — all 30 endpoints with request/response JSON examples, query params, valid enums |
| **docs/03_features.md** | Feature guide — 14 features explained with factory scenarios, S006/S014 compliance matrices, README feature verification |
| **docs/04_sql_queries.md** | SQL queries guide — all 18 queries (material/WO, genealogy traces, equipment history, analytics) with how-it-works explanations |
| **docs/05_testing.md** | Testing guide — Bruno `.bru` format, test generation flow, running tests, CI/CD integration |

## Testing

The project includes an automated, sequential integration test suite using the Bruno REST client.

### Running the Automated Suite

A shell script is provided to compile the application, start a background instance, generate the dynamic Bruno test files, and execute the test runner.

```bash
chmod +x run_tests.sh
./run_tests.sh
```

**Expected output:**
```
📊 Execution Summary
┌───────────────┬────────────────┐
│ Status        │     ✓ PASS     │
│ Requests      │ 24 (24 Passed) │
│ Assertions    │     47/47      │
└───────────────┴────────────────┘
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
├── api/traceability/v1/              # Protocol Buffer definitions (gRPC/HTTP contracts)
│   ├── traceability.proto            # Service interface definition
│   ├── traceability.pb.go            # Generated protobuf code
│   ├── traceability_grpc.pb.go       # Generated gRPC server/client stubs
│   └── traceability_http.pb.go       # Generated HTTP gateway code
│
├── cmd/traceability/                 # Application entry point
│   ├── main.go                       # Server startup, Wire injection
│   └── wire_gen.go                   # Dependency injection (generated)
│
├── configs/                          # Environment configurations
│   ├── config.yaml                   # Development config
│   └── config_docker.yaml            # Docker config
│
├── database/
│   ├── queries/                      # SQL query reference documentation
│   │   └── traceability_queries.sql  # All 18 queries with comments
│   └── schema/
│       ├── schema.sql                # Master DDL (16 tables, indexes, hypertables)
│       ├── 002_traceability_features.sql # S006/S014-specific schema
│       └── seed.sql                  # EV Battery plant simulation data
│
├── docs/                             # Comprehensive documentation (NEW)
│   ├── 01_database.md                # Table reference + real examples
│   ├── 02_api_reference.md           # All 30 endpoints + examples
│   ├── 03_features.md                # 14 features + S006/S014 matrices
│   ├── 04_sql_queries.md             # All 18 queries explained
│   └── 05_testing.md                 # Bruno testing guide
│
├── internal/
│   ├── biz/                          # Business logic, domain models, interfaces
│   │   ├── traceability.go           # Core interfaces (CreateLot, TraceBackward, etc.)
│   │   ├── biz.go                    # Biz layer initialization
│   │   └── README.md                 # Architecture notes
│   │
│   ├── data/                         # PostgreSQL/GORM data access layer
│   │   ├── traceability.go           # All CRUD and CTE query implementations
│   │   ├── data.go                   # Data layer initialization
│   │   └── README.md                 # GORM patterns
│   │
│   ├── service/                      # HTTP API handlers and DTO mapping
│   │   ├── service.go                # Service layer initialization
│   │   └── traceability.go           # API request/response handlers
│   │
│   ├── server/                       # HTTP/gRPC server initialization
│   │   ├── http.go                   # HTTP server with middleware
│   │   ├── grpc.go                   # gRPC server
│   │   └── server.go                 # Server setup
│   │
│   └── conf/                         # Configuration structs
│       ├── conf.pb.go                # Generated from protobuf config
│       └── conf.proto                # Config schema
│
├── tests/api/                        # Auto-generated Bruno integration tests
│   ├── 01_Enterprise/ through 17_Analytics/  # 14 test folders (auto-generated)
│   ├── environments/Dev VM.bru       # Base URL variables
│   └── bruno.json                    # Collection metadata
│
├── third_party/                      # Protobuf third-party dependencies
│
├── generate_tests.py                 # Python test generator (dynamically creates .bru files)
├── run_tests.sh                      # CI/CD test automation script
├── Dockerfile                        # Multi-stage Docker build
├── Makefile                          # Build targets (api, wire, build, etc.)
├── go.mod / go.sum                   # Go module dependencies
├── openapi.yaml                      # OpenAPI 3.0 spec (auto-generated from protos)
└── README.md                         # This file
```

## Key Implementation Details

### Database Schema
* **16 tables** covering hierarchy, inventory, production, genealogy, and telemetry
* **Recursive CTEs** for backward/forward traces (PostgreSQL 12+)
* **TimescaleDB hypertable** (`equipment_telemetry`) for efficient time-series queries
* **Cascading deletes** on hierarchy to maintain referential integrity

### API Design
* **RESTful** with protobuf-defined JSON request/response bodies
* **Dual transport:** HTTP (REST gateway) + gRPC for clients
* **Time-window-based analytics:** All trend/comparison queries are parameterized by date range
* **Flexible filtering:** List endpoints support optional `where` clauses (status, material_type, etc.)

### Performance Optimizations
* **Indexes on all foreign keys** and frequently-filtered columns
* **Recursive CTE limits** to prevent infinite loops in genealogy (configurable depth)
* **TimescaleDB compression** for old telemetry data
* **Batch operations** for high-volume lot/event inserts

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on how to submit pull requests, report issues, and suggest improvements.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.


