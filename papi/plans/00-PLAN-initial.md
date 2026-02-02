# Partner API Platform - Strategic Plan

## Vision

The Partner API Platform establishes a unified, stable interface for external partners to consume selected capabilities from internal company services. Rather than exposing internal APIs directly—with their varying conventions, versioning schemes, and implementation details—the platform presents a curated, consistent API surface that evolves independently of internal changes.

This approach serves two primary stakeholders:

1. **External Partners**: Receive a clean, well-documented API with predictable versioning and minimal breaking changes
2. **Internal Teams**: Retain autonomy over their service implementations while contributing selected capabilities to the partner-facing surface

---

## Problem Statement

### Current Challenges

**For Partners:**
- Internal APIs were not designed for external consumption
- Inconsistent naming conventions, error formats, and authentication patterns across services
- Internal versioning changes create unexpected breaking changes for consumers
- No unified documentation or developer experience

**For Internal Teams:**
- Pressure to maintain backward compatibility limits internal evolution
- No clear boundary between internal and external concerns
- Ad-hoc partner integrations create maintenance burden
- Lack of visibility into how partners use their APIs

**For the Organization:**
- No centralized control over what is exposed to partners
- Security and compliance concerns with direct internal API access
- Difficulty enforcing rate limits, quotas, and usage policies
- No unified analytics on partner API consumption

### Desired Outcomes

1. **Decoupled Evolution**: Internal services evolve freely; Partner API maintains stability
2. **Centralized Governance**: Single point of control for partner-facing surface
3. **Operational Visibility**: Per-partner usage tracking, rate limiting, and analytics
4. **Developer Experience**: Consistent, well-documented API for partners and clear contribution model for internal teams

---

## Architectural Approach

### Why an API Gateway Pattern

The API Gateway pattern places a dedicated layer between external consumers and internal services. This architectural choice addresses our core challenges:

**Abstraction**: The gateway presents a partner-centric API design that may differ from internal representations. Field names, response structures, and endpoint paths can be transformed to match partner expectations without modifying internal services.

**Stability**: By owning the partner-facing contract, the gateway can absorb internal changes. When an internal service introduces a breaking change, the gateway's transformation layer adapts, preserving partner compatibility.

**Control**: Authentication, rate limiting, and access policies are enforced at a single point. This simplifies security auditing and enables consistent policy application across all partner interactions.

**Observability**: All partner traffic flows through the gateway, providing comprehensive visibility into usage patterns, error rates, and performance characteristics per partner.

### Why Hybrid Spec Management

Internal teams understand their domains best and should own their API specifications. However, the platform team must control what reaches partners and how it's presented. The hybrid model resolves this tension:

**Team Ownership**: Internal teams maintain OpenAPI specifications in their own repositories, alongside their service code. This preserves existing workflows and ensures specs stay synchronized with implementations.

**Platform Control**: The governance repository references team specs and defines transformation rules. The platform team controls which endpoints are exposed, how paths are mapped, and what fields are included or excluded.

**Automated Aggregation**: A pipeline fetches team specs, applies transformations, and produces the unified Partner API specification. This automation ensures consistency and reduces manual coordination overhead.

This model avoids two problematic alternatives:
- **Fully centralized**: Teams lose ownership; specs drift from implementations
- **Fully federated**: Platform loses control; inconsistent partner experience

### Why Build-Time Composition

The Partner API specification is assembled at build/deploy time rather than dynamically at runtime:

**Predictability**: The exact API surface is known before deployment. Changes are reviewed, tested, and approved through standard merge request workflows.

**Validation**: Contract tests and breaking change detection run against the composed specification. Issues are caught before reaching production.

**Documentation**: Partner documentation is generated from the composed specification, ensuring accuracy and completeness.

**Auditability**: The governance repository provides a complete history of what was exposed to partners and when.

Runtime composition would introduce uncertainty about the actual API surface and complicate validation and documentation generation.

---

## Core Components

### API Gateway

The gateway handles runtime concerns: routing requests to appropriate internal services, enforcing authentication and rate limits, transforming requests and responses, and managing resilience patterns like circuit breaking.

**Selection Criteria**: The gateway must support declarative configuration, response aggregation from multiple backends, JWT validation, and integration with Kubernetes service discovery. It should be stateless for horizontal scaling and support the organization's language preferences for any custom extensions.

**Tool Selection**: KrakenD (Go-based) meets these criteria with native support for multi-backend aggregation, JWT validation, and declarative configuration. Its stateless architecture simplifies Kubernetes deployment.

### Spec Aggregation Pipeline

The pipeline transforms and combines internal team specifications into the unified Partner API specification.

**Responsibilities**:
- Fetch specifications from internal team repositories based on manifest configuration
- Apply transformation rules: path remapping, field filtering, schema prefixing
- Merge specifications into a single OpenAPI document
- Validate the result against organizational standards
- Generate gateway configuration from the specification

**Tool Selection**: openapi-merge-cli handles specification merging with path modification and conflict resolution. Redocly CLI provides transformation capabilities through decorators. These tools integrate into standard CI/CD pipelines.

### Governance Repository

The central repository where the platform team maintains control over the Partner API surface.

**Contents**:
- Manifest defining which team specifications to include
- Transformation rules for each team's contribution
- Shared schema definitions used across multiple domains
- Aggregation pipeline configuration
- Gateway configuration templates
- CI/CD pipeline definitions

**Workflow**: Changes to the Partner API surface—whether adding new endpoints, modifying transformations, or updating policies—flow through merge requests in this repository, enabling review and approval processes.

### Contract Testing Framework

Contract testing ensures internal services fulfill the contracts they've committed to in the Partner API.

**Approach**: Provider-driven contract testing where internal teams validate their implementations match their specifications. Each internal team runs contract tests in their own CI pipelines using reusable templates provided by the platform team. This differs from consumer-driven testing (like Pact) which would require the platform team to define expectations for 8-15 internal services—an impractical coordination burden.

**Components**:
- **Breaking change detection**: Runs in internal team CI; compares specification changes against the Partner API spec, blocking merges that would break existing partner integrations
- **Implementation validation**: Runs in internal team CI; tests that running services respond according to their specifications
- **Specification linting**: Runs in internal team CI; enforces organizational standards for API design consistency

The platform team provides CI templates that internal teams include in their pipelines. This ensures consistent testing standards while keeping execution distributed—internal teams are responsible for their own contract compliance.

**Tool Selection**: oasdiff detects breaking changes in OpenAPI specifications. Schemathesis performs property-based testing of implementations against specifications. Spectral provides customizable linting rules. All tools are packaged into reusable GitLab CI templates.

### Authentication and Authorization

Partners authenticate using a federated identity model that leverages existing OLX platform credentials.

**Identity Federation**: Partners already have accounts in the OLX platform (managed by Auth0). Rather than creating separate credentials, partners log into the Partner API Portal using their existing OLX credentials. Keycloak acts as an identity broker, trusting Auth0 as the upstream identity provider.

**Application Registration**: Once authenticated, partners can register one or more "applications" through the portal. Each application receives its own `client_id` and `client_secret`. This allows partners to separate credentials for different environments (production, staging) or different integrations.

**API Authentication**: Partner applications use the OAuth 2.0 Client Credentials flow to obtain access tokens. The gateway validates these tokens and extracts application identity for rate limiting and analytics.

**Why This Model**:
- Partners use credentials they already have (no new passwords to manage)
- Clear separation between partner identity (Auth0) and application credentials (Keycloak)
- Partners can manage multiple applications independently
- Instant credential provisioning without manual approval

**Tool Selection**: Keycloak deployed on Kubernetes provides identity brokering with Auth0 and OAuth client management. A middleware service handles application registration logic and enforces per-partner quotas.

### Partner Portal

A unified portal combining public API documentation with authenticated self-service capabilities.

**Why a Unified Portal**: Partners benefit from a single destination for all their needs—exploring API capabilities, reading documentation, registering applications, and managing credentials. Separating documentation from self-service creates friction and navigation confusion.

**Public Section** (no authentication required):
- Auto-generated API reference from OpenAPI specification
- Use-case recipes and integration guides
- Multi-version documentation with version switcher
- Changelog and migration guides
- Search across all content

**Authenticated Section** (requires login via Auth0):
- Application registration and management
- Credential viewing and rotation
- Usage dashboard showing requests against rate limits
- API version selection per application

**Access Model**: Documentation is publicly accessible—partners can explore capabilities before committing. Authentication is only required for application management, triggered when partners click "Register Application" or "Sign In".

**Tool Selection**: A hybrid approach using Docusaurus for documentation and a React-based dashboard for self-service, unified under consistent navigation and branding. Partners authenticate via Keycloak, which redirects to Auth0.

### Observability

Observability follows organizational standards using NewRelic. The gateway exports metrics and traces via OpenTelemetry, providing per-partner visibility into usage patterns, error rates, and performance. Detailed observability design is deferred to implementation phase, leveraging existing organizational practices and tooling.

---

## Tooling Overview

The following tools support the platform's core functions. Selection criteria prioritized self-hosted/open-source options, compatibility with Java/Kotlin/Go ecosystems, and CI/CD integration.

| Category | Tool | Purpose |
|----------|------|---------|
| **Gateway** | KrakenD Community Edition | Request routing, response aggregation, JWT validation, rate limiting, transformations |
| **Spec Aggregation** | openapi-merge-cli | Merge multiple OpenAPI specs with path modification and conflict resolution |
| **Spec Transformation** | Redocly CLI | Filter endpoints, remove internal fields, apply decorators for partner-facing specs |
| **Breaking Change Detection** | oasdiff | Compare OpenAPI specs to detect breaking changes; blocks MRs in CI |
| **Spec Linting** | Spectral | Enforce API design standards and organizational conventions |
| **Implementation Testing** | Schemathesis | Property-based testing validating services match their OpenAPI specs |
| **Identity Broker** | Keycloak | Identity brokering with Auth0, OAuth client management, token issuance |
| **Upstream Identity Provider** | Auth0 (existing) | Partner authentication using existing OLX platform credentials |
| **Portal Documentation** | Docusaurus | Static site generator with versioning, Markdown content, and plugin ecosystem |
| **Portal Self-Service** | React + Keycloak JS | Authenticated dashboard for application management |
| **API Reference Generation** | docusaurus-plugin-openapi-docs | Auto-generate interactive API docs from OpenAPI spec |
| **Portal Search** | @easyops-cn/docusaurus-search-local | Offline/local search without external service dependencies |
| **Observability** | OpenTelemetry + NewRelic | Metrics and traces using organizational standard |
| **GitOps Deployment** | ArgoCD | Declarative, Git-based continuous deployment |
| **Kubernetes** | AWS EKS | Managed Kubernetes for all platform components |

---

## Hybrid Spec Management Model

### Conceptual Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         INTERNAL TEAM REPOSITORIES                          │
│                                                                             │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                    │
│   │   Orders    │    │    Users    │    │  Payments   │     ...            │
│   │  Team Repo  │    │  Team Repo  │    │  Team Repo  │                    │
│   │             │    │             │    │             │                    │
│   │ openapi.yaml│    │ openapi.yaml│    │ openapi.yaml│                    │
│   └──────┬──────┘    └──────┬──────┘    └──────┬──────┘                    │
│          │                  │                  │                            │
└──────────┼──────────────────┼──────────────────┼────────────────────────────┘
           │                  │                  │
           │    ┌─────────────┴─────────────┐    │
           │    │      MANIFEST DEFINES     │    │
           └────┤   - Which specs to fetch  ├────┘
                │   - Source locations      │
                │   - Version/branch refs   │
                └─────────────┬─────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        GOVERNANCE REPOSITORY                                │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                    TRANSFORMATION RULES                              │  │
│   │                                                                      │  │
│   │   • Path remapping:    /api/v2/orders → /partner/orders             │  │
│   │   • Field filtering:   Remove internal fields                        │  │
│   │   • Schema prefixing:  Avoid naming collisions                       │  │
│   │   • Tag-based selection: Include only partner-tagged endpoints       │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                              │                                              │
│                              ▼                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                    AGGREGATION PIPELINE                              │  │
│   │                                                                      │  │
│   │   Fetch → Transform → Merge → Validate → Generate                    │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                              │                                              │
│                              ▼                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                         OUTPUTS                                      │  │
│   │                                                                      │  │
│   │   • partner-api.yaml    (Unified OpenAPI specification)             │  │
│   │   • krakend.json        (Gateway routing configuration)              │  │
│   │   • Documentation       (Generated partner docs)                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Manifest Purpose

The manifest serves as the declarative source of truth for what constitutes the Partner API:

- **Source References**: Repository URLs, file paths, and branch/tag references for each team's specification
- **Transformation Bindings**: Links each source to its transformation rules
- **Version Tracking**: Enables pinning to specific versions for stability or tracking HEAD for continuous updates

### Transformation Rules Purpose

Transformation rules define how each team's specification is adapted for partner consumption:

- **Path Remapping**: Internal paths may follow team conventions; partner paths follow platform standards
- **Field Filtering**: Internal responses may include fields not appropriate for external exposure
- **Endpoint Selection**: Teams may have many endpoints; only some are partner-relevant
- **Schema Namespacing**: Multiple teams may have similarly-named schemas; prefixing prevents collisions

### Sync Mechanisms

Two triggers initiate specification synchronization:

1. **Manifest Changes**: When the platform team modifies the manifest (adding teams, changing versions), the pipeline runs
2. **Team Spec Changes**: Webhooks notify the governance repository when internal team specifications change, triggering synchronization

The second mechanism ensures the Partner API stays current with internal team updates without requiring manual intervention.

---

## Implementation Phases

### Phase 1: Foundation

**Objective**: Establish the core infrastructure and validate the architectural approach with a minimal set of internal services.

**Why This Phase First**: Before building governance workflows or onboarding many teams, we need confidence that the gateway routes correctly, the aggregation pipeline produces valid output, and the overall architecture meets requirements. Starting with 2-3 internal services limits blast radius while proving the concept.

**Key Activities**:

1. **Gateway Deployment**
   - Deploy gateway to Kubernetes with basic configuration
   - Configure routing to 2-3 selected internal services
   - Validate request forwarding and response handling
   - Establish health checking and basic monitoring

2. **Repository Structure**
   - Create governance repository with planned directory structure
   - Establish manifest format and initial entries
   - Define transformation rule schema
   - Set up CI/CD pipeline skeleton

3. **Initial Aggregation Pipeline**
   - Implement spec fetching from internal team repositories
   - Configure basic transformation rules
   - Produce aggregated Partner API specification
   - Generate gateway configuration from specification

4. **Validation Foundation**
   - Implement specification linting with organizational rules
   - Validate pipeline output against OpenAPI standards
   - Establish baseline tests for gateway routing

**Phase Exit Criteria**:
- Partner can call gateway endpoints and receive responses from internal services
- Aggregation pipeline produces valid, linted OpenAPI specification
- Gateway configuration is generated from specification
- CI pipeline runs validation on changes

**Risks to Address**:
- Internal service network accessibility from gateway
- Specification format variations across teams
- Transformation rule expressiveness for initial use cases

---

### Phase 2: Governance and Contract Testing

**Objective**: Establish the workflows and tooling that enable internal teams to contribute safely to the Partner API.

**Why This Phase Second**: With the foundation proven, we can now build the processes that scale to 8-15 teams. Governance workflows and contract testing must be in place before broader onboarding to prevent quality and stability issues.

**Key Activities**:

1. **Governance Workflow**
   - Configure branch protection and approval requirements
   - Establish CODEOWNERS for different repository areas
   - Document contribution process for internal teams
   - Create templates for adding new team specifications

2. **Breaking Change Detection**
   - Integrate breaking change detection into CI pipeline
   - Define what constitutes a breaking change in organizational context
   - Establish process for handling necessary breaking changes
   - Configure blocking behavior for unauthorized breaking changes

3. **Contract Testing Templates**
   - Create reusable CI configuration for internal teams
   - Implement specification validation against running services
   - Document integration steps for internal teams
   - Establish feedback mechanisms for test failures

4. **Initial Team Onboarding**
   - Onboard 2-3 additional internal teams using established process
   - Gather feedback on governance workflow usability
   - Refine documentation and tooling based on experience
   - Validate contract testing catches real issues

**Phase Exit Criteria**:
- Merge requests require appropriate approvals
- Breaking specification changes are blocked in CI
- Internal teams can integrate contract testing with provided templates
- At least 5 internal teams contributing to Partner API

**Risks to Address**:
- Governance overhead discouraging team participation
- False positives in breaking change detection
- Contract test reliability across different service implementations

---

### Phase 3: Identity and Access Control

**Objective**: Implement production-grade authentication using federated identity and self-service application registration.

**Why This Phase Third**: Security implementation requires stable gateway and routing. The identity model builds on existing Auth0 infrastructure, minimizing new credential management burden.

**Key Activities**:

1. **Keycloak Deployment**
   - Deploy Keycloak on Kubernetes using the Keycloak Operator
   - Configure PostgreSQL database for Keycloak persistence
   - Set up TLS and external access via ALB
   - Create `partner-api` realm

2. **Auth0 Identity Brokering**
   - Configure Keycloak to trust Auth0 as OIDC identity provider
   - Set up attribute mapping from Auth0 claims
   - Configure first-login flow for new partners
   - Test end-to-end authentication flow

3. **Application Registration Service**
   - Build middleware service for application management
   - Implement client creation via Keycloak Admin API
   - Add per-partner application quotas
   - Implement credential rotation

4. **Rate Limiting**
   - Design rate limit tiers based on partner needs
   - Configure per-application rate limiting in gateway
   - Implement rate limit headers in responses
   - Store tier information in OAuth client attributes

**Phase Exit Criteria**:
- Partners can log in using Auth0 credentials
- Partners can register applications and receive credentials
- Applications can authenticate via Client Credentials flow
- Gateway validates tokens and enforces rate limits

**Risks to Address**:
- Auth0 integration complexity
- Token validation latency impact on request performance
- Rate limit accuracy with stateless gateway instances

---

### Phase 4: Partner Portal and Production Readiness

**Objective**: Launch the unified Partner Portal and achieve production readiness.

**Why This Phase Fourth**: The portal depends on stable authentication and API surface. Production hardening benefits from understanding real failure modes observed in earlier phases.

**Key Activities**:

1. **Documentation Portal**
   - Deploy Docusaurus with auto-generated API reference
   - Author initial set of use-case recipes
   - Implement version switcher and changelog
   - Configure search functionality

2. **Self-Service Dashboard**
   - Build React dashboard with Keycloak authentication
   - Implement application management UI
   - Add credential display and rotation
   - Integrate usage statistics display

3. **Portal Integration**
   - Unify navigation between docs and dashboard
   - Apply consistent branding
   - Test authentication flows
   - Ensure seamless public/authenticated transitions

4. **Production Deployment**
   - Finalize production Kubernetes manifests
   - Configure horizontal pod autoscaling
   - Establish deployment procedures and rollback processes
   - Complete security review

**Phase Exit Criteria**:
- Partner Portal published with API reference, recipes, and search
- Partners can register applications through portal
- Partners can view and rotate credentials
- All planned internal teams contributing to Partner API

**Risks to Address**:
- Portal UX meeting partner expectations
- Documentation accuracy as platform evolves
- Recipe coverage not addressing actual partner use cases

---

## Versioning Strategy

### Partner API Versioning

The Partner API maintains its own version independent of internal service versions. This decoupling is fundamental to the platform's value proposition.

**Version Format**: Date-based versions (e.g., `2024-01-15`) following patterns established by Stripe and Twilio. Partners specify their API version via header, defaulting to their onboarding version.

**Backward Compatibility Commitment**: Once a version is released, its contract remains stable. New capabilities are added in new versions; existing capabilities are not modified or removed within a version.

**Version Lifecycle**:
1. **Active**: Current recommended version for new integrations
2. **Supported**: Previous versions maintained for existing partners
3. **Deprecated**: Announced end-of-life with migration timeline
4. **Retired**: No longer available

### Handling Internal Changes

When internal services introduce changes that would affect the Partner API:

**Non-Breaking Changes**: New optional fields, new endpoints, or expanded functionality flow through automatically (if transformation rules permit) without requiring Partner API version changes.

**Breaking Changes**: Removed fields, renamed endpoints, or changed semantics are absorbed by the transformation layer. The gateway maps partner requests to the new internal representation while preserving the partner-facing contract.

**Version Migration**: When breaking changes accumulate or transformation complexity becomes unsustainable, a new Partner API version is released with the updated contract. Partners migrate on their timeline within the deprecation window.

### Transformation Layer Approach

Rather than a separate "adaptation service," transformations are handled at two levels:

1. **Specification Level**: The aggregation pipeline transforms specifications at build time—path remapping, field filtering, schema modifications
2. **Runtime Level**: The gateway handles request/response transformation—header injection, field mapping, response merging

This layered approach keeps simple transformations configuration-driven using KrakenD's declarative configuration. When transformations cannot be expressed declaratively, a sidecar service provides full programming language capabilities. This two-tier model minimizes operational complexity while handling complex scenarios.

---

### Breaking Change Lifecycle

#### The Coordination Challenge

When internal services evolve, some changes would break the Partner API contract if exposed directly. The platform must absorb these changes through transformations, but this requires coordination: the transformation must exist *before* the internal change reaches production.

Without a defined workflow, two failure modes emerge:
- **Internal team blocked indefinitely**: They wait for transformation with no clear process
- **Partner API breaks**: Internal change deploys before transformation is ready

#### Detection: Catching Breaking Changes Early

Internal team CI pipelines include breaking change detection that compares their OpenAPI specification against the Partner API specification. When a change would break the partner contract, the pipeline fails with a clear message directing the team to the coordination workflow.

This detection runs on every merge request, ensuring breaking changes are identified before code is merged—not after deployment when partners are affected.

**What constitutes a breaking change**:
- Removing or renaming fields in responses
- Removing or renaming endpoints
- Changing field types (string → object, etc.)
- Adding required fields to request bodies
- Narrowing allowed values (enum restrictions)
- Changing authentication requirements
- Changing default behavior (e.g., sorting, pagination)

**What is NOT a breaking change**:
- Adding new optional fields to responses
- Adding new endpoints
- Adding optional parameters to requests
- Expanding allowed values
- Performance improvements
- Bug fixes that align behavior with documented contract

#### Coordination Workflow

When breaking change detection blocks a merge request:

1. **Notification**: Internal team posts in the dedicated Slack channel (`#partner-api-changes`) describing the intended change and business rationale

2. **Ticket Creation**: Internal team creates a JIRA ticket containing:
   - Description of the change (what field/endpoint is changing, why)
   - Old format example (request/response as it is today)
   - New format example (request/response after their change)
   - Link to their blocked merge request
   - Target deployment date (if any)

3. **Platform Team Implementation**: Platform team implements the transformation that converts the new internal format back to the old Partner API format for partners on previous versions

4. **Testing**: Platform team runs integration tests against the internal team's staging environment to validate the transformation handles real responses correctly

5. **Deployment Sequencing**:
   - Platform team deploys transformation to production gateway
   - Platform team confirms transformation is active
   - Platform team unblocks internal team (comment on JIRA/MR)
   - Internal team proceeds with their deployment

#### Transformation Storage and Structure

Transformations are version-aware functions that convert current internal responses to previous Partner API formats. They reside in the governance repository alongside the gateway configuration.

Each transformation is:
- **Scoped to a specific version boundary**: e.g., "transform responses for partners on versions before 2024-06-01"
- **Focused on a specific change**: e.g., "flatten customer object to customer_name field"
- **Independently testable**: can be unit tested with example inputs/outputs from the JIRA ticket

When multiple breaking changes accumulate, transformations chain: a response for a partner on `2024-01-15` may pass through multiple transformations to reach the expected format.

#### Testing Transformations

Before deployment, transformations are validated through:

1. **Unit tests**: Using the example request/response pairs from the JIRA ticket
2. **Integration tests**: Calling the internal team's staging environment through the transformation and validating the output matches expected Partner API format

Integration tests run in CI when transformation code changes, ensuring transformations remain valid as they evolve.

#### Failure Modes and Recovery

**If transformation has a bug discovered after deployment**:
- Gateway logs will show transformation errors for affected partner versions
- Platform team can quickly deploy a fix or rollback
- Internal team's change is unaffected (their service is already deployed)

**If internal team deploys before transformation is ready**:
- Partners on old versions receive malformed responses
- This is why the workflow requires platform team confirmation before internal deployment
- If it happens accidentally, platform team treats transformation as urgent hotfix

---

## Success Criteria

### Platform Health Metrics

| Metric | Target | Rationale |
|--------|--------|-----------|
| Availability | 99.9% | Partners depend on API for their applications |
| P95 Latency | <200ms overhead | Gateway should add minimal latency to backend responses |
| Error Rate | <0.1% platform errors | Distinguish platform errors from backend/client errors |

### Process Metrics

| Metric | Target | Rationale |
|--------|--------|-----------|
| Breaking Changes in Production | 0 | Contract testing should catch all breaking changes |
| Spec-to-Production Lead Time | <1 day | Changes should flow quickly once approved |
| Internal Team Onboarding | <1 week | Low friction encourages participation |

### Adoption Metrics

| Metric | Target | Rationale |
|--------|--------|-----------|
| Internal Team Coverage | 100% of planned teams | Platform value increases with comprehensiveness |
| Partner Satisfaction | Measured via feedback | Ultimate measure of platform success |

---

## Component Plans

Detailed implementation guidance for each component resides in separate documents:

| Plan | Scope |
|------|-------|
| [01-gateway-setup.md](./plans/01-gateway-setup.md) | Gateway deployment, configuration patterns, Kubernetes manifests |
| [02-spec-aggregation.md](./plans/02-spec-aggregation.md) | Manifest structure, transformation rules, aggregation pipeline |
| [03-governance-workflow.md](./plans/03-governance-workflow.md) | Repository structure, approval workflows, contribution process |
| [04-contract-testing.md](./plans/04-contract-testing.md) | Testing strategy, CI templates, tool configuration |
| [05-security-auth.md](./plans/05-security-auth.md) | Identity federation, application registration, rate limiting |
| [07-versioning-strategy.md](./plans/07-versioning-strategy.md) | Version policy, transformation patterns, migration procedures |
| [08-partner-portal.md](./plans/08-partner-portal.md) | Unified portal setup, documentation, self-service dashboard |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 | 2026-01-28 | Initial draft |
| 0.2 | 2026-01-28 | Revised to focus on strategy and rationale; removed implementation details |
| 0.3 | 2026-01-28 | Added Developer Portal component for partner documentation and discovery |
| 0.4 | 2026-01-28 | Added Breaking Change Lifecycle workflow within Versioning Strategy |
| 0.5 | 2026-01-28 | Added Tooling Overview section; clarified provider-driven testing runs in internal team pipelines |
| 0.6 | 2026-01-29 | Revised authentication to use Auth0 federation via Keycloak; unified portal combining docs and self-service; deferred observability details to implementation |
| 0.7 | 2026-01-30 | Updated transformation approach to KrakenD declarative + sidecar (removed Lua/CEL); added complete breaking change definitions; established canonical terminology |
