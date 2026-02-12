# Taksa App - Traceability Service

A comprehensive REST/gRPC service for manufacturing equipment traceability, built with Go, Kratos framework, and PostgreSQL with TimescaleDB.

## Features

- **38 RESTful API endpoints** across 10 entity types (Enterprise, Site, Area, Production Line, Production Unit, Equipment Class, Equipment, Capability, Property, Logs)
- **gRPC server** for high-performance communication
- **PostgreSQL database** with TimescaleDB support for event logging
- **Clean architecture** with layered design (biz/data/service)
- **Comprehensive test suite** with Bruno REST client

## Prerequisites

- Go 1.24+
- PostgreSQL 12+ with TimescaleDB extension
- Protoc compiler for proto generation
- Make

## Installation & Setup

### 1. Install Dependencies

```bash
# Download Go dependencies
make init

# Install Kratos and other tools
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
go get github.com/google/wire/cmd/wire
```

### 2. Database Setup

```bash
# Create database and schema
psql -U postgres -c "CREATE DATABASE traceability_db;"
psql -U traceability_db < database/schema/schema.sql

# Or with custom credentials
psql -U $DB_USER -c "CREATE DATABASE traceability_db;"
psql -U $DB_USER traceability_db < database/schema/schema.sql
```

### 3. Configuration

The service uses `configs/config.yaml` for configuration. Database credentials can be provided via environment variable:

```bash
export DB_SOURCE="postgres://user:password@localhost:5432/traceability_db?sslmode=disable"
```

Or the default configuration will be used:
```yaml
data:
  database:
    driver: postgres
    source: ${DB_SOURCE:-postgres://taksa:taksa_123@127.0.0.1:5433/traceability_db?sslmode=disable}
```

## Building

```bash
# Generate proto files
make api

# Generate Wire dependency injection
cd cmd/traceability && wire && cd ../..

# Build the service
go build -o traceability ./cmd/traceability

# Build everything
make all
```

## Running

```bash
# Start the service
./traceability -conf ./configs

# Or with Docker
docker build -t traceability:latest .
docker run --rm -p 8000:8000 -p 9000:9000 \
  -v $(pwd)/configs:/data/conf \
  --env DB_SOURCE="postgres://user:password@db:5432/traceability_db" \
  traceability:latest
```

The service will start on:
- **HTTP**: `http://localhost:8000`
- **gRPC**: `localhost:9000`

## API Endpoints

The API provides comprehensive CRUD operations for the following entities:

1. **Enterprise** - Top-level organizational units
2. **Site** - Physical manufacturing sites
3. **Area** - Functional areas within sites
4. **Production Line** - Equipment lines within areas
5. **Production Unit** - Individual workstations
6. **Equipment Class** - Classification of equipment types
7. **Equipment** - Individual equipment instances
8. **Capability** - Equipment capabilities and specifications
9. **Property** - Equipment property values
10. **Logs** - Traceability event logs

Each entity supports:
- `Create` - Add new entity
- `List` - Retrieve filtered list
- `Update` - Modify existing entity
- `Delete` - Remove entity

## Testing

### Bruno REST Client

Comprehensive REST API tests are provided in the `tests/api/` directory:

```bash
# Install Bruno: https://www.usebruno.com/
# Open the collection
tests/api/bruno.json

# Select Dev VM environment (http://localhost:8000)
# Run the test suite in order: 1_Enterprise → ... → 9_Logs → 99_Cleanup
```

### Manual Testing

Example with curl:

```bash
# Create an enterprise
curl -X POST http://localhost:8000/api/v1/traceability/enterprises \
  -H "Content-Type: application/json" \
  -d '{"name":"My Company","description":"Main facility"}'

# List enterprises
curl http://localhost:8000/api/v1/traceability/enterprises
```

## Project Structure

```
├── api/                          # Protocol Buffer definitions
│   └── traceability/v1/         # Traceability service protos
├── cmd/                         # Entry points
│   └── traceability/            # Main service binary
├── configs/                     # Configuration files
├── database/                    # Database schemas
│   └── schema/                  # SQL migrations
├── internal/                    # Internal packages
│   ├── biz/                     # Business logic layer
│   ├── data/                    # Data/repository layer
│   ├── service/                 # Service/handler layer
│   ├── conf/                    # Configuration interface
│   └── server/                  # Server setup
├── tests/                       # Test files
│   └── api/                     # Bruno REST API tests
└── third_party/                 # External proto definitions
```

## Architecture

The service follows a clean architecture pattern:

```
Proto Definitions (api/)
    ↓
    Service Layer (internal/service/) - HTTP/gRPC handlers
    ↓
    Business Logic (internal/biz/) - Domain logic and validation
    ↓
    Data Layer (internal/data/) - Repository pattern with GORM
    ↓
    Database (PostgreSQL + TimescaleDB)
```

## Development

### Generating Proto Files

After modifying protos in `api/traceability/v1/traceability.proto`:

```bash
make api
```

### Wire Injection

To update dependency injection:

```bash
cd cmd/traceability
wire
```

### Code Generation

```bash
# Regenerate all proto and injection code
go generate ./...
```

## Docker Deployment

```bash
# Build image
docker build -t artpark/traceability:latest .

# Run with environment
docker run -d \
  --name traceability \
  -p 8000:8000 \
  -p 9000:9000 \
  -e DB_SOURCE="postgres://user:password@db:5432/traceability_db" \
  artpark/traceability:latest
```

## License

Copyright (c) 2026 The Taksa Project Authors. All rights reserved.

## Support

For issues and questions, please refer to the project documentation or submit an issue.

