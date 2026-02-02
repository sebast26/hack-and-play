# 01 - Gateway Setup

## Purpose

The API Gateway is the runtime component that sits between external partners and internal services. It is the only entry point for partner traffic, responsible for enforcing security policies, routing requests to appropriate backends, transforming responses for version compatibility, and providing operational visibility.

This plan defines what the gateway must do and why, establishing requirements that implementation tasks must fulfill.

---

## Role in the Platform

The gateway serves as the **enforcement point** for everything the Partner API promises:

| Promise to Partners | Gateway Responsibility |
|---------------------|------------------------|
| Stable API contract | Apply version-specific transformations so partners see consistent responses regardless of internal changes |
| Authentication | Validate OAuth tokens before requests reach internal services |
| Rate limits | Enforce per-partner quotas to protect platform and ensure fair usage |
| Low latency | Route efficiently to internal services; minimize processing overhead |
| High availability | Remain operational even when some internal services fail |

The gateway does NOT:
- Generate its own configuration (receives config from spec aggregation pipeline)
- Issue OAuth tokens (Keycloak handles token issuance)
- Store state (stateless for horizontal scaling)
- Know about partner business details (only sees partner ID from token for rate limiting/logging)

---

## Core Requirements

### Request Routing

The gateway must route incoming partner requests to the correct internal service endpoints.

**Why this matters**: Partners call `/partner/orders/{id}`, but the internal service might expose `/api/v3/orders/{id}` on a different host. The gateway abstracts this mapping.

**Requirements**:
- Map partner-facing paths to internal service paths
- Support path parameter forwarding (e.g., `{id}` passed through)
- Support query parameter forwarding (selective or all)
- Support header forwarding (selective, with ability to inject/remove headers)
- Route to internal services via Kubernetes DNS (`service.namespace.svc.cluster.local`)

---

### Response Aggregation (Fan-out)

Some Partner API endpoints require data from multiple internal services combined into a single response.

**Why this matters**: Partners expect a single call for a "dashboard" view, but internally this data lives across user service, order service, and notification service.

**Requirements**:
- Call multiple backends in parallel for a single partner request
- Merge responses into a unified JSON structure
- Support grouping backend responses under named keys (e.g., `{"user": {...}, "orders": {...}}`)
- Handle partial failures gracefully (return available data with indicator that some backends failed)

**Known Patterns**:

| Pattern | Description | Example |
|---------|-------------|---------|
| **Parallel independent** | Multiple backends called simultaneously; responses merged | Dashboard endpoint combining user profile, recent orders, notifications |
| **Parallel with shared parameter** | Same parameter passed to multiple backends | Partner ID used to fetch data from multiple domain services |
| **Sequential dependent** | Second call depends on first call's response | Fetch order, then fetch shipping details using carrier ID from order |

Sequential dependent calls are more complex and may require a dedicated composition service. The gateway should support parallel patterns natively; sequential patterns are handled case-by-case.

**Partial Response Format**:

When one or more backends fail during fan-out, the gateway returns a partial response with a `_meta` field indicating the incomplete state:

```json
{
  "user": {
    "id": "u123",
    "name": "John Doe"
  },
  "orders": null,
  "_meta": {
    "partial": true,
    "failed": ["orders"],
    "message": "Some data could not be retrieved"
  }
}
```

Partners can check for the presence of `_meta.partial` to detect incomplete responses and decide whether the available data is sufficient for their use case.

---

### Authentication Validation

The gateway validates that incoming requests carry valid OAuth tokens before routing to internal services.

**Why this matters**: Internal services should not need to understand partner authentication. They trust the gateway to only forward authenticated requests.

**Requirements**:
- Validate JWT tokens on every request (except explicitly public endpoints like health checks)
- Verify token signature against Keycloak's public keys (JWKS endpoint)
- Check token expiration
- Extract partner identifier from token claims for rate limiting and logging
- Reject requests with invalid/expired tokens with appropriate HTTP status (401)
- Cache JWKS to avoid per-request calls to Keycloak

**Header Propagation**:
- Extract claims from JWT (e.g., `partner_id`, `scopes`)
- Inject extracted values as headers to internal services (e.g., `X-Partner-ID`)
- Internal services can use these headers for logging/auditing but do not validate them

---

### Rate Limiting

The gateway enforces request quotas per partner to protect internal services and ensure fair usage.

**Why this matters**: Without rate limiting, a single partner could overwhelm internal services, affecting all partners.

**Requirements**:
- Enforce rate limits per partner (identified by partner ID from JWT)
- Support configurable limits (requests per second/minute)
- Return standard rate limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`)
- Return 429 Too Many Requests when limit exceeded
- Limits should be configurable per endpoint or globally

**Rate Limit Tiers**:

Rate limits are organized into tiers (Standard, Enhanced, Premium). For tier definitions, specific values, and the tier upgrade process, see [05-security-auth.md](./05-security-auth.md). This document covers enforcement mechanics only.

The gateway reads the tier from the JWT token's `rate_limit_tier` claim and applies the corresponding limits.

**Consideration - Stateless Rate Limiting**:

KrakenD Community Edition performs rate limiting per gateway instance (not cluster-wide). With multiple gateway pods, actual limits are multiplied by pod count. This is acceptable if:
- Limits are set conservatively (accounting for pod count)
- Traffic is reasonably distributed across pods
- Precise enforcement is not critical

If precise cluster-wide limits become necessary, options include:
- Redis-backed rate limiting (requires KrakenD Enterprise or custom solution)
- External rate limiting service
- Accepting approximate limits as sufficient

---

### Version-Aware Transformations

The gateway transforms responses to match the Partner API version each partner is using.

**Why this matters**: Partners on older API versions must continue receiving responses in the format they expect, even as internal services evolve.

**Requirements**:
- Extract Partner API version from request header (`X-Partner-API-Version`)
- Apply appropriate transformation based on version
- Support transformation chaining (response passes through multiple transformations for very old versions)
- Handle retired/unsupported versions appropriately
- Pass through unchanged for current version

**Transformation Capabilities Needed**:
- Field renaming (`customer` → `customer_name`)
- Field flattening (`customer.name` → `customer_name`)
- Field removal (hide internal fields)
- Field restructuring (flat → nested, nested → flat)
- Conditional logic (if field exists, transform; otherwise skip)

**Implementation Approach**:

Transformations follow a two-tier model:

| Option | When to Use |
|--------|-------------|
| **KrakenD declarative configuration** | First choice for all transformations. Use KrakenD's built-in response manipulation features (response mapping, filtering, grouping, field extraction). |
| **Sidecar service** | Only when the transformation cannot be expressed using KrakenD's declarative configuration. Provides full programming language capabilities (Kotlin/Go) for complex business logic. |

The decision rule is simple: **always use KrakenD declarative configuration first**. Only introduce a sidecar when the required transformation is impossible to express declaratively. This minimizes operational complexity and keeps most transformations as configuration rather than code.

**Retired Version Handling**:

When a partner requests a version that has been retired, the gateway returns **410 Gone** with a helpful error response:

```json
{
  "error": "version_retired",
  "message": "Version 2023-06-01 is no longer available",
  "retired_on": "2024-06-01",
  "minimum_supported_version": "2024-01-15",
  "migration_guide": "https://developers.company.com/docs/migration/2023-06-01"
}
```

This status code clearly communicates that the version existed previously but has been permanently removed, distinguishing it from invalid version requests (400 Bad Request) or unknown versions.

---

### Resilience

The gateway must remain operational and useful even when some internal services fail.

**Why this matters**: If one internal service is down, partners should still be able to use endpoints that don't depend on it.

**Requirements**:

**Circuit Breaking**:
- Track error rates per backend service
- Open circuit (stop calling) when error threshold exceeded
- Return fast failure instead of waiting for timeouts
- Periodically attempt recovery (half-open state)
- Log circuit state changes for operational visibility

**Timeouts**:
- Configurable timeout per backend
- Timeout on initial connection and on response
- Return 504 Gateway Timeout when backend doesn't respond in time

**Graceful Degradation for Aggregation**:
- When one backend in a fan-out fails, return partial response with remaining data
- Include `_meta` field in response body indicating partial data and which backends failed (see Response Aggregation section)
- Partners can decide whether partial data is acceptable for their use case

---

### Observability

The gateway must emit metrics and traces for operational visibility.

**Why this matters**: The platform team needs to understand traffic patterns, identify issues, and provide per-partner usage analytics.

**Requirements**:

**Metrics** (exported via OpenTelemetry to NewRelic):
- Request count (by endpoint, partner, status code)
- Latency histograms (by endpoint, backend)
- Error rates (by endpoint, error type)
- Rate limit hits (by partner)
- Circuit breaker state changes

**Tracing**:
- Generate trace ID for each request
- Propagate trace ID to internal services
- Include spans for: JWT validation, each backend call, transformation

**Logging**:
- Structured logs (JSON)
- Include: trace ID, partner ID, endpoint, status, latency
- Log level configurable (ERROR in production, DEBUG for troubleshooting)

**Health Endpoints**:
- Liveness probe: gateway process is running
- Readiness probe: gateway can accept traffic (dependencies reachable)

---

## Integration Points

The gateway integrates with other platform components. These integration points define the "contract" between the gateway and other systems.

### From Spec Aggregation Pipeline

The pipeline generates gateway configuration. The gateway expects:
- Endpoint definitions (partner path → internal service mapping)
- Backend service addresses (Kubernetes DNS)
- Timeout values per endpoint/backend
- Transformation references (which transformations apply to which endpoints/versions)

### From Security/Auth Component

The gateway depends on Keycloak for token validation:
- JWKS endpoint URL for public key retrieval
- Expected token issuer and audience values
- Claim names for partner identification

### From Governance Repository

Transformations are stored in the governance repository:
- KrakenD configuration snippets for declarative transformations
- Sidecar service configuration (when declarative is insufficient)
- Version-to-transformation mapping

### To Observability Stack

The gateway exports telemetry:
- OpenTelemetry metrics endpoint (scraped by collector)
- OpenTelemetry traces (pushed to collector)
- Structured logs to stdout (collected by Kubernetes logging)

### To Internal Services

The gateway calls internal services:
- Via Kubernetes service DNS
- Using service-to-service authentication (organizational standard)
- Forwarding relevant headers (partner ID, trace ID)

---

## Deployment Characteristics

### Stateless Operation

The gateway stores no state locally. All configuration comes from mounted ConfigMaps or environment variables. This enables:
- Horizontal scaling (add/remove pods freely)
- Rolling deployments without session concerns
- Pod restarts without data loss

### Configuration Reload

KrakenD does not support hot-reload of configuration. Configuration changes require pod restart. This is acceptable because:
- Configuration changes flow through GitOps (ArgoCD detects changes, triggers rollout)
- Rolling deployment ensures zero downtime
- Configuration changes are infrequent (tied to spec changes)

### Scaling Considerations

- Scale based on CPU utilization (JWT validation and transformation are CPU-bound)
- Minimum 3 replicas for high availability
- Pod anti-affinity to spread across nodes/zones
- Horizontal Pod Autoscaler for traffic-based scaling

---

## What This Plan Does NOT Cover

The following are explicitly out of scope for gateway setup and belong to other component plans:

| Topic | Covered In |
|-------|------------|
| How configuration is generated from OpenAPI spec | 02-spec-aggregation.md |
| OAuth token issuance and partner credential management | 05-security-auth.md |
| Version lifecycle and deprecation policies | 07-versioning-strategy.md |
| Partner Portal integration | 08-partner-portal.md |

---

## Success Criteria

The gateway implementation is complete when:

| Criteria | Validation |
|----------|------------|
| Routes requests to internal services | Partner can call gateway endpoint and receive internal service response |
| Validates JWT tokens | Requests without valid token receive 401; valid tokens pass through |
| Enforces rate limits | Exceeding limit returns 429 with appropriate headers |
| Applies version transformations | Partner on old version receives old response format |
| Rejects retired versions | Retired version returns 410 Gone with migration guidance |
| Handles backend failures | Circuit opens on failures; partial responses returned for fan-out with `_meta` field |
| Emits observability data | Metrics visible in NewRelic; traces show request flow |
| Scales horizontally | Multiple pods handle traffic; no single point of failure |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 | 2026-01-28 | Initial draft |
| 0.2 | 2026-01-29 | Resolved open questions: rate limit tiers reference 05-security-auth.md; transformation approach simplified to KrakenD declarative + sidecar (removed Lua); partial responses use `_meta` field; retired versions return 410 Gone |
