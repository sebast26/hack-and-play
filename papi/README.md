# Partner API Platform

A unified API gateway for external partners, built with KrakenD and docker-compose.

## Overview

This project implements Phase 1 of the Partner API Platform:
- **KrakenD API Gateway** - Routes partner requests to internal services
- **Mock Backend Services** - Orders and Users services for testing
- **Spec Aggregation Pipeline** - Merges OpenAPI specs and generates gateway config

## Prerequisites

- Docker and docker-compose
- Node.js (for openapi-merge-cli)
- yq (`brew install yq`)
- jq (`brew install jq`)

## Quick Start

### 1. Install Dependencies

```bash
npm install
```

### 2. Run the Spec Aggregation Pipeline

Generate the merged Partner API spec and KrakenD configuration:

```bash
./scripts/aggregate.sh
```

This will:
- Transform internal specs to partner-facing paths
- Merge specs into `governance/output/partner-api.yaml`
- Generate `gateway/krakend.json`

### 3. Start the Services

```bash
docker-compose up --build
```

### 4. Test the Gateway

```bash
# Health check
curl http://localhost:8080/__health

# List orders
curl http://localhost:8080/partner/orders

# Get specific order
curl http://localhost:8080/partner/orders/123

# List users
curl http://localhost:8080/partner/users

# Get specific user
curl http://localhost:8080/partner/users/456
```

## Project Structure

```
papi/
├── docker-compose.yml          # Container orchestration
├── package.json                # npm dependencies
├── gateway/
│   └── krakend.json            # Gateway configuration (generated)
├── mocks/
│   ├── orders/                 # Orders mock service
│   │   ├── main.go
│   │   ├── Dockerfile
│   │   └── openapi.yaml
│   └── users/                  # Users mock service
│       ├── main.go
│       ├── Dockerfile
│       └── openapi.yaml
├── governance/
│   ├── manifest.yaml           # Source definitions
│   ├── openapi-merge.yaml      # Merge configuration
│   ├── spec-mappings/          # Path transformation rules
│   │   ├── orders-mapping.yaml
│   │   └── users-mapping.yaml
│   ├── specs/                  # Transformed specs (generated)
│   │   ├── orders/
│   │   └── users/
│   ├── shared-schemas/
│   │   └── common.yaml
│   └── output/
│       └── partner-api.yaml    # Merged Partner API spec (generated)
└── scripts/
    ├── aggregate.sh            # Main pipeline script
    ├── transform-spec.sh       # Path transformation
    └── generate-krakend.sh     # KrakenD config generator
```

## Endpoints

| Partner API Path | Internal Service | Internal Path |
|------------------|------------------|---------------|
| `GET /partner/orders` | orders:8081 | `/api/v3/orders` |
| `GET /partner/orders/{id}` | orders:8081 | `/api/v3/orders/{id}` |
| `GET /partner/users` | users:8082 | `/api/v1/users` |
| `GET /partner/users/{id}` | users:8082 | `/api/v1/users/{id}` |

## Spec Aggregation Pipeline

The pipeline transforms internal OpenAPI specs into a unified Partner API spec.

### How It Works

1. **Transform** - Reads specs from `mocks/` and applies path mappings from `spec-mappings/`
2. **Merge** - Combines transformed specs using openapi-merge-cli
3. **Generate** - Creates KrakenD configuration from the merged spec

### Running the Pipeline

```bash
# Full pipeline
./scripts/aggregate.sh

# Individual steps
./scripts/transform-spec.sh orders    # Transform orders spec
./scripts/transform-spec.sh users     # Transform users spec
./scripts/generate-krakend.sh         # Generate KrakenD config
```

### Adding a New Service

1. Create the mock service in `mocks/<service-name>/`
2. Add OpenAPI spec with `partner` tag on exposed endpoints
3. Create mapping file in `governance/spec-mappings/<service-name>-mapping.yaml`:
   ```yaml
   source_base_path: /api/v1/<service-name>
   target_base_path: /partner/<service-name>
   backend:
     host: http://<service-name>:<port>
   include:
     tags:
       - partner
   ```
4. Add source to `governance/manifest.yaml`
5. Add input to `governance/openapi-merge.yaml`
6. Run `./scripts/aggregate.sh`

## Verification

### Verify Spec Generation

```bash
# Check transformed specs exist
ls -la governance/specs/*/openapi.yaml

# Check merged spec paths
yq '.paths | keys' governance/output/partner-api.yaml

# Expected output:
# - /partner/orders
# - /partner/orders/{id}
# - /partner/users
# - /partner/users/{id}
```

### Verify KrakenD Configuration

```bash
# Check generated endpoints
cat gateway/krakend.json | jq '.endpoints[].endpoint'

# Validate JSON syntax
cat gateway/krakend.json | jq . > /dev/null && echo "Valid JSON"
```

### Verify End-to-End with Docker

```bash
# Start services
docker-compose up -d --build

# Wait for services to be ready
sleep 5

# Test all endpoints
echo "Testing /partner/orders..."
curl -s http://localhost:8080/partner/orders | jq '.orders | length'

echo "Testing /partner/orders/123..."
curl -s http://localhost:8080/partner/orders/123 | jq '.order_id'

echo "Testing /partner/users..."
curl -s http://localhost:8080/partner/users | jq '.users | length'

echo "Testing /partner/users/456..."
curl -s http://localhost:8080/partner/users/456 | jq '.user_id'

# Check gateway logs
docker-compose logs gateway

# Stop services
docker-compose down
```

## Development Workflow

### Making Changes to Mock Services

1. Edit the service code in `mocks/<service>/main.go`
2. Update `mocks/<service>/openapi.yaml` if API changes
3. Run `./scripts/aggregate.sh` to regenerate configs
4. Restart services: `docker-compose up -d --build`

### Making Changes to Path Mappings

1. Edit `governance/spec-mappings/<service>-mapping.yaml`
2. Run `./scripts/aggregate.sh`
3. Restart gateway: `docker-compose restart gateway`

## Troubleshooting

### Pipeline Fails with "yq not found"

```bash
brew install yq
```

### Pipeline Fails with "openapi-merge-cli not found"

```bash
npm install
```

### Gateway Returns 404

1. Check if endpoint exists in `gateway/krakend.json`
2. Verify path mapping in `spec-mappings/`
3. Re-run `./scripts/aggregate.sh`

### Backend Service Not Reachable

1. Check if service is running: `docker-compose ps`
2. Check service logs: `docker-compose logs <service>`
3. Verify host in `spec-mappings/<service>-mapping.yaml`

## What's Not Included (Deferred)

- Authentication (Keycloak/JWT validation)
- Rate limiting
- Spectral linting validation
- Breaking change detection
- CI/CD integration
