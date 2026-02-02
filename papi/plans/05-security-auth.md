# 05 - Security and Authentication

## Purpose

Security and Authentication ensures that only authorized partners can access the Partner API, that usage is tracked and limited appropriately, and that communication between platform components is secure.

This plan defines the identity federation model, application registration process, rate limiting strategy, and integration with existing organizational security practices.

---

## Role in the Platform

Security serves as the **trust and access control layer**:

| Concern | How It's Addressed |
|---------|-------------------|
| Partner identity verification | Auth0 via Keycloak identity brokering |
| Application credentials | Keycloak OAuth clients with Client Credentials flow |
| Credential management | Self-service via Partner Portal |
| Usage limits | Three-tier rate limiting per application |
| Abuse prevention | Rate limiting per application |
| Internal communication | Existing shared secrets mechanism |

Security does NOT:
- Create new partner accounts (partners use existing OLX/Auth0 credentials)
- Handle partner billing (separate concern)
- Manage partner business relationships (sales/support concern)
- Replace internal team authentication mechanisms

---

## Identity Model Overview

The authentication model separates **partner identity** from **application credentials**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           IDENTITY MODEL                                    │
│                                                                             │
│   PARTNER IDENTITY                    APPLICATION CREDENTIALS               │
│   ─────────────────                   ────────────────────────              │
│                                                                             │
│   Partner logs in with                Partner registers applications        │
│   existing OLX credentials            that receive OAuth credentials        │
│                                                                             │
│   ┌─────────┐      ┌─────────┐       ┌─────────────────────────────────┐   │
│   │  Auth0  │◄────►│Keycloak │       │  Application 1: Production      │   │
│   │  (OLX)  │      │(Broker) │       │  - client_id: acme-corp_a1b2c3  │   │
│   └─────────┘      └────┬────┘       │  - client_secret: ****          │   │
│                         │            └─────────────────────────────────┘   │
│                         │                                                   │
│                         │            ┌─────────────────────────────────┐   │
│                         └───────────►│  Application 2: Staging         │   │
│                                      │  - client_id: acme-corp_d4e5f6  │   │
│                                      │  - client_secret: ****          │   │
│                                      └─────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Why this separation**:
- Partners don't need new credentials (use what they already have)
- Partners can have multiple applications with separate credentials
- Credential compromise affects only one application, not the partner's identity
- Clear audit trail: which partner owns which application

---

## Partner Authentication (Portal Access)

Partners access the Partner Portal using their existing OLX platform credentials through identity federation.

### Authentication Flow

```
Partner                    Partner Portal              Keycloak                 Auth0
   │                            │                         │                       │
   ├── Click "Sign In" ────────►│                         │                       │
   │                            ├── Redirect to Keycloak ►│                       │
   │                            │                         ├── Redirect to Auth0 ─►│
   │                            │                         │                       │
   │◄─────────────────────────────────── Auth0 Login Page ─────────────────────────┤
   │                            │                         │                       │
   ├── Enter OLX credentials ──────────────────────────────────────────────────────►│
   │                            │                         │                       │
   │                            │                         │◄── Auth0 token ───────┤
   │                            │                         │                       │
   │                            │                         ├── Create/update user  │
   │                            │                         │    in Keycloak        │
   │                            │                         │                       │
   │                            │◄── Keycloak session ────┤                       │
   │◄── Authenticated ──────────┤                         │                       │
   │    (can manage apps)       │                         │                       │
```

### Keycloak as Identity Broker

Keycloak acts as an identity broker, trusting Auth0 for partner authentication:

| Keycloak Role | Description |
|---------------|-------------|
| **Identity Broker** | Accepts Auth0 tokens, creates local user sessions |
| **Session Management** | Manages partner sessions for portal access |
| **OAuth Server** | Issues tokens for partner applications (separate from portal sessions) |
| **Client Registry** | Stores application credentials and metadata |

**Why Keycloak (not direct Auth0)**:
- Keycloak manages OAuth clients (applications) locally
- No need for Auth0 Machine-to-Machine setup per partner application
- Full control over token claims and client attributes
- Self-hosted on Kubernetes (organizational preference)

### Auth0 Integration Configuration

Keycloak connects to Auth0 as an OpenID Connect identity provider. A **dedicated Auth0 application** is created within the existing OLX Auth0 tenant specifically for the Partner API portal.

| Configuration | Value |
|---------------|-------|
| Provider Type | OpenID Connect v1.0 |
| Discovery URL | `https://{tenant}.auth0.com/.well-known/openid-configuration` |
| Client ID | Configured in Auth0 for Partner Portal (new application) |
| Client Secret | Stored in Kubernetes secret |
| Default Scopes | `openid profile email` |

**Auth0 Setup Requirements**:
- New Auth0 application created within existing OLX tenant
- Callback URLs configured for Keycloak
- Requires coordination with Auth0 administrators
- Separate from other OLX applications for clean configuration

**Attribute Mapping**: Auth0 claims are mapped to Keycloak user attributes:

| Auth0 Claim | Keycloak Attribute | Purpose |
|-------------|-------------------|---------|
| `email` | username | Unique identifier |
| `email` | email | Contact information |
| `given_name` | firstName | Display name |
| `family_name` | lastName | Display name |
| `sub` | auth0_id | Link to Auth0 identity |

### Portal Session Duration

Partner portal sessions last **24 hours** before requiring re-authentication.

| Setting | Value | Rationale |
|---------|-------|-----------|
| Session timeout | 24 hours | Full workday coverage across time zones |
| Idle timeout | 24 hours | Same as session timeout |
| Remember me | Not implemented | Standard session duration is sufficient |

After 24 hours, partners are redirected to login again. This balances security with convenience for developers working on integrations.

### First Login Behavior

When a partner authenticates via Auth0 for the first time:

1. Keycloak receives Auth0 token
2. Keycloak creates local user linked to Auth0 identity
3. User attributes populated from Auth0 claims
4. Partner lands in portal, can immediately register applications

No manual approval required—if partner has valid Auth0 credentials, they can access the portal.

---

## Application Registration

Partners register applications through the Partner Portal to obtain API credentials.

### Application Concept

An "application" represents a set of credentials for API access:

| Attribute | Description |
|-----------|-------------|
| Name | Human-readable identifier (e.g., "Production Integration") |
| client_id | Unique identifier for OAuth (format: `{partner-slug}_{random}`) |
| client_secret | Secret for Client Credentials flow |
| Owner | Partner who created the application |
| Rate Limit Tier | standard, enhanced, or premium |
| Created Date | When application was registered |
| API Version | Default Partner API version |

**One partner can have multiple applications**:
- Separate credentials for production/staging environments
- Different integrations with different purposes
- Credential rotation without affecting other applications
- **No limit** on the number of applications per partner

### Client ID Format

Generated client_ids follow the pattern: **`{partner-slug}_{random-suffix}`**

Examples:
- `acme-corp_a1b2c3d4`
- `widgets-inc_x7y8z9w0`
- `big-retailer_m3n4o5p6`

**Why this format**:
- Human-readable: immediately identify which partner in logs and dashboards
- Unique: random suffix prevents collisions
- Traceable: easy to filter metrics and logs by partner
- Short enough: doesn't create overly long identifiers

The partner slug is derived from the partner's Auth0 profile (company name or email domain), normalized to lowercase with hyphens.

### Registration Flow

```
Partner (authenticated)           Portal                    Middleware              Keycloak
        │                           │                           │                      │
        ├── "Create Application" ──►│                           │                      │
        │   name: "My Prod App"     │                           │                      │
        │                           ├── Validate session ──────►│                      │
        │                           │                           │                      │
        │                           │◄── Partner ID ────────────┤                      │
        │                           │                           │                      │
        │                           ├── Create app request ────►│                      │
        │                           │                           │                      │
        │                           │                           ├── Generate client_id │
        │                           │                           │   (partner-slug +    │
        │                           │                           │    random suffix)    │
        │                           │                           │                      │
        │                           │                           ├── Create OAuth ─────►│
        │                           │                           │   client             │
        │                           │                           │                      │
        │                           │                           │◄── client_id, ───────┤
        │                           │                           │    client_secret     │
        │                           │                           │                      │
        │                           │                           ├── Store metadata     │
        │                           │                           │   (owner, name)      │
        │                           │                           │                      │
        │                           │◄── Credentials ───────────┤                      │
        │◄── client_id, ────────────┤                           │                      │
        │    client_secret          │                           │                      │
        │    (shown once)           │                           │                      │
```

### Application Registration Middleware

A middleware service handles application registration logic:

**Why middleware (not direct Keycloak access)**:
- Applies naming conventions to client_id (partner slug + random)
- Stores ownership metadata
- Provides audit logging
- Abstracts Keycloak Admin API complexity

**Middleware Responsibilities**:

| Responsibility | Description |
|----------------|-------------|
| Session validation | Verify partner is authenticated |
| Partner slug extraction | Derive slug from partner profile |
| Client ID generation | Generate `{partner-slug}_{random}` format |
| Keycloak client creation | Call Admin API to create OAuth client |
| Metadata storage | Track owner, name, creation date |
| Audit logging | Log all credential operations |

### Instant Provisioning

Application registration is instant—no approval workflow:

**Why instant**:
- Reduces friction for partner onboarding
- Partners can start integrating immediately
- Rate limits provide protection against abuse

**Guardrails**:
- Rate limit on registration endpoint (e.g., 5 per minute)
- All applications start at Standard rate limit tier
- Tier upgrades require manual request

---

## API Authentication (Client Credentials Flow)

Partner applications authenticate to the API using OAuth 2.0 Client Credentials flow.

### Authentication Flow

```
Partner Application                    Keycloak                         Gateway
        │                                 │                                │
        ├── POST /token ─────────────────►│                                │
        │   grant_type=client_credentials │                                │
        │   client_id=acme-corp_a1b2c3d4  │                                │
        │   client_secret=****            │                                │
        │                                 │                                │
        │◄── access_token ────────────────┤                                │
        │    (JWT, expires in 1 hour)     │                                │
        │                                 │                                │
        ├── GET /partner/orders ──────────┼───────────────────────────────►│
        │   Authorization: Bearer <token> │                                │
        │                                 │                                ├── Validate JWT
        │                                 │                                ├── Extract app_id
        │                                 │                                ├── Check rate limits
        │                                 │                                ├── Route to backend
        │◄────────────────────────────────┼────────────────────────────────┤
        │   Response                      │                                │
```

### Token Contents

The JWT access token contains claims used by the gateway:

| Claim | Purpose |
|-------|---------|
| `sub` | Application identifier (client_id) |
| `iss` | Token issuer (Keycloak URL) |
| `aud` | Intended audience (Partner API identifier) |
| `exp` | Expiration timestamp |
| `iat` | Issued-at timestamp |
| `azp` | Authorized party (client_id) |
| `rate_limit_tier` | Application's rate limit tier |
| `owner_partner_id` | Partner who owns this application |

### Token Validation

The gateway validates tokens on every request:

| Validation | Action on Failure |
|------------|-------------------|
| Token present | 401 Unauthorized |
| Signature valid (JWKS) | 401 Unauthorized |
| Token not expired | 401 Unauthorized |
| Issuer matches expected | 401 Unauthorized |
| Audience matches expected | 401 Unauthorized |

**JWKS caching**: Gateway caches Keycloak's public keys (JWKS) to avoid per-request calls. Cache refreshes periodically and on signature validation failure.

---

## Credential Management

Partners manage application credentials through the Partner Portal.

### Viewing Applications

Partners see a list of their applications with:
- Application name
- client_id (always visible)
- Rate limit tier
- Creation date
- Last used date (if tracking available)

**client_secret is NOT shown** after initial creation.

### Credential Rotation

Partners can rotate credentials at any time:

1. Partner clicks "Rotate Secret" for an application
2. Confirmation dialog warns that old secret will be invalidated
3. New client_secret generated
4. Old client_secret immediately invalidated
5. New secret shown once (partner must save it)

**No grace period**: Rotation is immediate. Partners should update their integrations before rotating, or accept brief downtime.

**Why no grace period**:
- Simpler implementation
- If secret is compromised, immediate invalidation is desired
- Partners control when to rotate

### Deleting Applications

Partners can delete applications they no longer need:

1. Partner clicks "Delete" for an application
2. Confirmation dialog warns this is irreversible
3. OAuth client deleted from Keycloak
4. Application removed from partner's list
5. Any tokens issued for this client become invalid

---

## Rate Limiting

### Three-Tier Model

Rate limits protect the platform and ensure fair usage across partners.

| Tier | Requests/Second | Requests/Day | Intended For |
|------|-----------------|--------------|--------------|
| **Standard** | 10 | 50,000 | New applications, evaluation, small integrations |
| **Enhanced** | 50 | 250,000 | Production integrations, medium volume |
| **Premium** | 200 | 1,000,000 | High-volume partners, strategic integrations |

**Tier assignment**:
- New applications start at Standard
- Tier upgrades requested via support process (see below)
- Tier stored as attribute in Keycloak OAuth client
- Gateway reads tier from JWT token claim

### Tier Upgrade Process

Partners request tier upgrades through the standard support process:

**Process**:
1. Partner submits support ticket (JIRA) requesting tier upgrade
2. Ticket includes: application name, current tier, requested tier, justification
3. Platform team reviews request
4. If approved, platform team updates tier in Keycloak
5. Partner notified of approval; new limits effective immediately

**Approval Criteria**:
- Legitimate business need for higher limits
- Good standing (no abuse history)
- For Premium tier: may require business/sales involvement

**Tracking**: All tier changes tracked in JIRA with full audit trail.

### Rate Limit Enforcement

Gateway enforces rate limits per application (not per partner):

| Check | Scope | Response on Exceed |
|-------|-------|-------------------|
| Requests per second | Rolling window | 429 Too Many Requests |
| Requests per day | Calendar day (UTC) | 429 Too Many Requests |

**Implementation Approach**:

KrakenD Community Edition uses **per-instance rate limiting** (each gateway pod tracks its own counts). With multiple pods, the effective cluster limit is the configured limit multiplied by pod count.

| Configuration | Description |
|---------------|-------------|
| Approach | Per-instance (stateless) |
| Limit calculation | Configured limit = desired limit ÷ expected pod count |
| Example | For 10 req/s with 3 pods, configure ~3 req/s per instance |

**Why per-instance (not Redis-backed)**:
- KrakenD CE doesn't support Redis-backed rate limiting (Enterprise feature)
- Per-instance is simpler, no Redis dependency
- Limits are approximate but sufficient for most use cases
- If precise limits become critical, revisit during implementation

**Monitoring**: Track actual rate limit behavior in NewRelic. If inaccuracy causes issues, consider KrakenD Enterprise or alternative solutions.

**Response headers** (always included):

| Header | Description |
|--------|-------------|
| `X-RateLimit-Limit` | Maximum requests allowed in window |
| `X-RateLimit-Remaining` | Requests remaining in current window |
| `X-RateLimit-Reset` | Unix timestamp when window resets |

**429 response body**:
```json
{
  "error": "rate_limit_exceeded",
  "message": "Request limit exceeded. Please retry after the reset time.",
  "retry_after": 1706472000
}
```

### Rate Limit Visibility

Partners can see usage in the Partner Portal:
- Per-application request counts
- Current usage vs. daily limit
- Historical usage graphs

**Usage Data Source**:

Usage data is sourced from **NewRelic with a caching layer**:

| Component | Description |
|-----------|-------------|
| Data source | NewRelic metrics (gateway exports with application ID dimension) |
| Query mechanism | Portal backend queries NewRelic API |
| Caching | Results cached in portal backend (e.g., 5-minute TTL) |
| Displayed metrics | Request counts, rate limit usage, error rates |

**Why NewRelic with caching**:
- Data already collected by gateway observability
- No additional infrastructure needed
- Caching reduces NewRelic API calls and improves portal responsiveness
- Cached data is sufficiently fresh for usage dashboards

---

## Keycloak Deployment

### Deployment Model

Keycloak runs self-hosted on Kubernetes (EKS), managed by the platform team.

| Aspect | Approach |
|--------|----------|
| Deployment method | Keycloak Operator on Kubernetes |
| High availability | Multiple replicas behind load balancer |
| Database | PostgreSQL (managed RDS) |
| TLS | ALB terminates TLS; internal traffic over HTTPS |
| Namespace | Dedicated namespace (`partner-api-auth`) |

**Why Keycloak Operator**:
- Kubernetes-native deployment
- Handles rolling updates, scaling
- CRDs for realm/client configuration
- Community supported

### Realm Configuration

Single realm for Partner API:

| Setting | Value |
|---------|-------|
| Realm name | `partner-api` |
| Login theme | Branded for Partner Portal |
| Token settings | Access token lifespan: 1 hour |
| Session settings | SSO Session Idle: 24 hours; SSO Session Max: 24 hours |

### OAuth Client Configuration

Partner applications are registered as Keycloak clients with these settings:

| Setting | Value | Reason |
|---------|-------|--------|
| Client Protocol | openid-connect | Standard OAuth/OIDC |
| Access Type | confidential | Client secret required |
| Service Accounts Enabled | Yes | Required for Client Credentials |
| Standard Flow Enabled | No | Not using authorization code flow |
| Direct Access Grants | No | Not using password grant |
| Valid Redirect URIs | (none) | No redirect-based flows |
| Token Lifespan | 1 hour | Balance security vs. convenience |

### Custom Attributes

Application-specific data stored as client attributes:

| Attribute | Purpose |
|-----------|---------|
| `owner_partner_id` | Partner who owns this application |
| `rate_limit_tier` | standard, enhanced, or premium |
| `app_name` | Human-readable application name |
| `created_at` | When application was registered |
| `api_version` | Default Partner API version |

These attributes are included in tokens as custom claims via protocol mappers.

---

## Gateway Integration

### JWT Validation Configuration

Gateway needs:

| Configuration | Source |
|---------------|--------|
| JWKS URL | Keycloak realm's JWKS endpoint |
| Expected issuer | Keycloak realm URL |
| Expected audience | Partner API identifier |
| Clock skew tolerance | 30 seconds (for time sync differences) |

### Claim Extraction

Gateway extracts from validated JWT:

| Claim | Used For |
|-------|----------|
| `sub` (application ID) | Rate limiting key, logging, metrics |
| `rate_limit_tier` | Selecting rate limit bucket |
| `owner_partner_id` | Logging, per-partner analytics |
| `api_version` | Default version if not in request header |

### Header Propagation

Gateway adds headers to requests forwarded to internal services:

| Header | Value | Purpose |
|--------|-------|---------|
| `X-Application-ID` | From `sub` claim | Internal logging/auditing |
| `X-Partner-ID` | From `owner_partner_id` claim | Partner identification |
| `X-Request-ID` | Generated UUID | Distributed tracing |

Internal services can use these headers for their own logging but don't need to validate them—they trust the gateway.

---

## Service-to-Service Authentication

Communication between gateway and internal services uses the organization's existing shared secrets mechanism.

**What this means**:
- Gateway authenticates to internal services using established patterns
- No new service-to-service auth mechanism introduced
- Internal teams don't need to change their auth setup
- Platform team configures gateway with appropriate secrets

**Configuration needed**:
- Gateway must be provisioned with secrets for each internal service
- Secrets management follows organizational standards
- Secret rotation follows organizational procedures

This plan doesn't prescribe implementation details for service-to-service auth since the organization has established practices.

---

## Security Considerations

### Token Security

| Concern | Mitigation |
|---------|------------|
| Token theft | Short expiration (1 hour); partners must secure tokens |
| Credential theft | Client secrets shown only once; rotation available |
| Replay attacks | Token expiration |
| Token tampering | JWT signature validation |

### Keycloak Security

| Concern | Mitigation |
|---------|------------|
| Admin access | Limited to platform team; MFA required |
| Database security | Encrypted at rest; Kubernetes secrets |
| Network exposure | Internal network only; ALB for external access |
| Audit logging | All admin actions logged |

### Application Registration Security

| Concern | Mitigation |
|---------|------------|
| Unauthorized registration | Requires Auth0 authentication |
| Runaway creation | Rate limiting on registration endpoint |
| Abuse | Rate limiting; monitoring for anomalies |
| Credential exposure | Secrets shown once; must re-rotate if lost |

### Auth0 Integration Security

| Concern | Mitigation |
|---------|------------|
| Auth0 compromise | Keycloak validates tokens cryptographically |
| Token interception | TLS everywhere |
| Session hijacking | Standard session security practices; 24-hour session limit |

---

## Integration Points

### With Auth0 (Upstream)

- Auth0 authenticates partner users
- Keycloak trusts Auth0 as OIDC identity provider
- New Auth0 application created within existing OLX tenant
- Auth0 configuration managed by existing OLX team (requires coordination)

### With Gateway

- Gateway validates JWTs using Keycloak JWKS
- Gateway extracts claims for rate limiting and logging
- Gateway enforces rate limits based on tier claim

### With Partner Portal

- Portal authenticates users via Keycloak (which uses Auth0)
- Portal sessions last 24 hours
- Portal creates/manages applications via middleware
- Portal displays application list and usage data (from NewRelic)

### With Middleware Service

- Middleware validates partner sessions
- Middleware calls Keycloak Admin API to create clients
- Middleware generates client_ids in partner-slug format
- Middleware enforces rate limits on registration

---

## Success Criteria

Security and authentication is complete when:

| Criteria | Validation |
|----------|------------|
| Partners can log in via Auth0 | Partner with OLX credentials can access portal |
| Sessions last 24 hours | Partner remains logged in for full day |
| Applications can be registered | Partner creates application, receives credentials |
| Client IDs follow naming convention | Generated client_id matches `{partner-slug}_{random}` format |
| Credentials work for API access | Client Credentials flow produces valid token |
| Invalid auth rejected | Missing/invalid/expired tokens receive 401 |
| Rate limits enforced | Exceeding limits produces 429 with correct headers |
| Tier upgrade process works | Support ticket leads to tier update |
| Credential rotation works | New secret works, old secret rejected |
| Usage visible in portal | Partner can see request counts and limits |
| Multiple applications supported | Partner can create and use multiple applications |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 | 2026-01-28 | Initial draft |
| 0.2 | 2026-01-29 | Revised for Auth0 federation model; partners authenticate via existing OLX credentials; applications registered separately from partner identity |
| 0.3 | 2026-01-30 | Resolved open questions: no application quota limit; support ticket for tier upgrades; new Auth0 app in existing tenant; 24-hour sessions; per-instance rate limiting with monitoring; NewRelic with caching for usage data; client_id format as partner-slug + random suffix |
