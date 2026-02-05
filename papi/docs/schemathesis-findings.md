# Schemathesis Test Findings

This document tracks issues found by Schemathesis implementation testing and their resolution status.

## Summary

| Issue | Severity | Status | Component |
|-------|----------|--------|-----------|
| 404 responses converted to 500 | High | Open | KrakenD Gateway |
| Missing 500 response in specs | Medium | Open | OpenAPI Specs |
| Missing 400 response in specs | Low | Open | OpenAPI Specs |
| TRACE method missing Allow header | Low | Won't Fix | KrakenD Gateway |

## Issue Details

### 1. 404 Responses Converted to 500

**Severity:** High

**Description:**
When a backend service returns 404 (Not Found), KrakenD gateway converts it to 500 (Internal Server Error) with the message "invalid status code".

**Reproduction:**
```bash
# Backend returns 404 correctly:
docker exec papi-orders-1 wget -qO- http://localhost:8081/api/v3/orders/nonexistent
# Returns: 404 Not Found

# But gateway returns 500:
curl http://localhost:8080/partner/orders/nonexistent
# Returns: 500 Internal Server Error - "invalid status code"
```

**Root Cause:**
KrakenD by default expects 2xx responses from backends and treats other status codes as errors.

**Resolution:**
Update `gateway/krakend.json` to configure proper error handling:

```json
{
  "backend": [{
    "url_pattern": "/api/v3/orders/{id}",
    "extra_config": {
      "backend/http": {
        "return_error_code": true
      }
    }
  }]
}
```

**Status:** Open - Requires gateway configuration update

---

### 2. Missing 500 Response in Specs

**Severity:** Medium

**Description:**
OpenAPI specs don't document 500 error responses, but the API can return them.

**Affected Endpoints:**
- `GET /partner/orders/{id}`
- `GET /partner/users/{id}`

**Resolution:**
Add 500 response to all endpoints in mock service specs:

```yaml
responses:
  "500":
    description: Internal server error
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/Error"
```

**Status:** Open - Requires spec updates

---

### 3. Missing 400 Response in Specs

**Severity:** Low

**Description:**
Invalid URL escapes (malformed URLs) return 400 Bad Request, which is not documented.

**Example:**
```bash
curl http://localhost:8080/partner/users/%09%0F%25lB
# Returns: 400 Bad Request - "error: invalid URL escape"
```

**Resolution:**
Document 400 response for endpoints with path parameters:

```yaml
responses:
  "400":
    description: Bad request - invalid path parameter format
```

**Status:** Open - Requires spec updates

---

### 4. TRACE Method Missing Allow Header

**Severity:** Low

**Description:**
When TRACE HTTP method returns 405 Method Not Allowed, the `Allow` header is missing (RFC 9110 requirement).

**Note:** TRACE is rarely used and testing it is often excessive. This check is excluded from default Schemathesis configuration.

**Resolution:**
Not planned - TRACE support is not required for Partner API.

**Status:** Won't Fix

---

## Recommended Actions

### Immediate (Gateway Fix)

1. Update KrakenD configuration to pass through 4xx status codes from backends
2. Test that 404 responses work correctly after the fix

### Short-term (Spec Updates)

1. Add 500 error response to all endpoints
2. Add 400 error response to endpoints with path parameters
3. Re-run Schemathesis to verify fixes

### Configuration

For CI, exclude overly strict checks until fixes are deployed:

```yaml
variables:
  SCHEMATHESIS_CHECKS: "not_a_server_error,content_type_conformance,response_schema_conformance"
```

## Test Commands

Run Schemathesis locally to verify fixes:

```bash
# Start services
docker compose up -d

# Run tests with core checks
./scripts/schemathesis-test.sh \
  --spec governance/output/partner-api.yaml \
  --url http://localhost:8080 \
  --checks "not_a_server_error,status_code_conformance,response_schema_conformance"

# Run all checks (will find more issues)
./scripts/schemathesis-test.sh \
  --spec governance/output/partner-api.yaml \
  --url http://localhost:8080 \
  --checks all
```
