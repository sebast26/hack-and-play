# 02 - Spec Aggregation Pipeline

## Purpose

The Spec Aggregation Pipeline is the build-time component that transforms and combines OpenAPI specifications from internal teams into the unified Partner API specification. It also generates the gateway configuration needed for runtime routing and transformation.

This pipeline is the bridge between internal team autonomy (they own their specs) and platform control (what partners see). It enforces that only approved, validated, properly-transformed endpoints reach the Partner API.

---

## Role in the Platform

The pipeline serves as the **composition and validation point** for the Partner API:

| Input | Output |
|-------|--------|
| Internal team OpenAPI specs (from their repos) | Unified Partner API OpenAPI spec |
| Transformation rules (from governance repo) | Gateway configuration (intermediate format) |
| Manifest defining what to include | Generated documentation source |

The pipeline does NOT:
- Run at request time (it's a build-time process)
- Modify internal team repositories (read-only access)
- Deploy anything (outputs are consumed by other processes)
- Make decisions about what to expose (follows manifest and transformation rules)

---

## Core Requirements

### Spec Storage Model

The governance repository stores individual team specs alongside the merged output. This provides a complete audit trail and enables incremental updates. See `03-governance-workflow.md` for the complete repository structure.

**Key Directories for Aggregation**:

```
partner-api-governance/
├── manifest.yaml                    # Defines sources and spec mappings
├── spec-mappings/                   # Rules for transforming specs (build-time)
│   ├── orders-mapping.yaml
│   ├── users-mapping.yaml
│   └── payments-mapping.yaml
├── specs/                           # Stored copies of internal team specs
│   ├── orders/
│   │   └── openapi.yaml
│   ├── users/
│   │   └── openapi.yaml
│   └── payments/
│       └── openapi.yaml
├── output/
│   ├── partner-api.yaml             # Merged Partner API specification
│   └── krakend.json                 # Generated gateway configuration
└── versions/
    └── changelog/
        └── latest.md                # Auto-generated changelog
```

**Why store team specs in governance repo**:
- Complete audit trail of exactly which version of each team spec was merged
- Pipeline can run without network access to team repos (useful for debugging/reproducing builds)
- Git history shows when each team's spec changed
- Enables incremental updates: only fetch what changed, read rest from repo

---

### Spec Fetching

The pipeline retrieves OpenAPI specifications from internal team repositories when changes occur.

**Why this matters**: Internal teams maintain specs alongside their code. The pipeline must fetch these specs reliably and update the governance repo.

**Requirements**:
- Fetch specs from GitLab repositories via API or git clone
- Support specifying branch, tag, or commit for each source
- Support different file paths per team (teams may organize repos differently)
- Authenticate using service account with read access to internal repos
- Update stored spec in governance repo after successful fetch
- Detect when a spec has changed by comparing against stored version

**Incremental Fetch Strategy**:

When a webhook indicates a team's spec changed:
1. Fetch **only** the changed team's spec from their repo
2. Compare against the stored spec in `specs/{team}/openapi.yaml`
3. If changed, update the stored spec in governance repo
4. Read all other team specs from governance repo (already stored)
5. Proceed with transformation and merge

This avoids unnecessary network calls while maintaining a complete, auditable record of all specs.

---

### Manifest-Driven Composition

The manifest defines which internal team specs are included and how they're mapped to the Partner API.

**Why this matters**: The platform team must have explicit control over what enters the Partner API. The manifest is the source of truth for this.

**Manifest Must Define**:

| Field | Purpose |
|-------|---------|
| Source repository URL | Where to fetch the spec |
| File path within repo | Location of OpenAPI spec file |
| Branch/tag/commit | Which version of the spec to use |
| Spec mapping reference | Which mapping rules to apply |
| Enabled flag | Ability to temporarily disable a source |

**Example Manifest Structure** (conceptual, not implementation):

```yaml
sources:
  - name: orders
    repository: https://gitlab.company.com/teams/orders-service
    path: api/openapi.yaml
    ref: main
    spec_mapping: orders-mapping.yaml
    enabled: true
    
  - name: users
    repository: https://gitlab.company.com/teams/users-service
    path: docs/api/spec.yaml
    ref: release-v2
    spec_mapping: users-mapping.yaml
    enabled: true
```

---

### Spec Mapping Rules

Spec mapping rules define how each internal spec is adapted for partner consumption at build time. These are distinct from version transformations (which happen at runtime for older API versions—see `07-versioning-strategy.md`).

**Why this matters**: Internal APIs are not designed for external use. Field names, paths, and structures must be mapped to match Partner API standards.

**Spec Mappings Must Support**:

| Mapping Type | Example |
|--------------|---------|
| Path remapping | `/api/v3/orders` → `/partner/orders` |
| Path prefix stripping/adding | Remove `/internal` prefix |
| Endpoint filtering by tag | Include only endpoints tagged `partner` |
| Endpoint filtering by path pattern | Exclude `/admin/*` paths |
| Field filtering in schemas | Remove `internal_id` field from response schemas |
| Field renaming in schemas | `customerId` → `customer_id` |
| Schema prefixing | Add `Orders_` prefix to avoid naming collisions |
| Description overrides | Replace internal descriptions with partner-friendly text |

**Implementation Approach**:

Spec mappings use a two-tier model (consistent with gateway transformation approach):

| Option | When to Use |
|--------|-------------|
| **YAML-based mapping rules** | First choice for all mappings. Declarative, reviewable, easy to understand. |
| **Sidecar service** | Only when the mapping cannot be expressed in YAML. Provides full programming language capabilities for complex business logic. |

The decision rule is simple: **always use YAML-based rules first**. Only introduce a sidecar when the required mapping is impossible to express declaratively. This keeps mappings auditable and accessible to team members without deep programming expertise.

**Spec Mapping Rule Structure** (conceptual):

```yaml
# spec-mappings/orders-mapping.yaml
source_base_path: /api/v3
target_base_path: /partner/orders

include:
  tags: [partner, public]
  
exclude:
  paths:
    - /api/v3/orders/internal/*
    - /api/v3/admin/*
  tags: [internal, admin]

schema_modifications:
  prefix: Orders
  remove_fields:
    - "**.internalId"
    - "**.debugInfo"
  rename_fields:
    customerId: customer_id
    createdAt: created_at
```

---

### Spec Merging

After transformation, individual specs are merged into a single Partner API specification.

**Why this matters**: Partners interact with one API, not 8-15 separate APIs. The merged spec must be valid, consistent, and conflict-free.

**Requirements**:
- Merge paths from all transformed specs into single OpenAPI document
- Merge schemas, avoiding naming collisions (via prefixing)
- Merge security schemes (though Partner API likely has one: OAuth)
- Set Partner API metadata (title, version, description, contact)
- Validate merged spec is valid OpenAPI 3.x
- Fail if path conflicts exist (same path from multiple sources)

---

### Validation

The pipeline validates specs at multiple stages.

**Why this matters**: Invalid or non-compliant specs should not reach partners. Catching issues early prevents runtime surprises.

**Validation Points**:

| Stage | What's Validated | Failure Behavior |
|-------|------------------|------------------|
| After fetch | Internal spec is valid OpenAPI | Pipeline fails; internal team notified |
| After fetch | Internal spec passes Spectral linting | Pipeline fails; internal team notified |
| After transform | Transformed spec is valid OpenAPI | Pipeline fails; platform team investigates |
| After merge | Merged Partner API spec is valid | Pipeline fails; platform team investigates |
| After merge | No breaking changes vs. previous version | Pipeline fails; triggers coordination workflow |

**Breaking Change Detection**:
- Compare newly generated Partner API spec against current production spec
- Use oasdiff to identify breaking changes
- If breaking changes detected from an internal team's update, block and notify
- This is the enforcement point for the Breaking Change Lifecycle workflow

---

### Gateway Configuration Generation

The pipeline generates gateway configuration from the merged spec.

**Why this matters**: The gateway needs routing rules, backend addresses, and transformation references. Generating this from the spec ensures gateway config stays in sync with the API definition.

**Intermediate Configuration Format**:

Rather than generating KrakenD-specific configuration directly, the pipeline produces a gateway-agnostic intermediate format. A separate step converts this to KrakenD config. This provides:
- Easier understanding for internal teams (no KrakenD knowledge needed)
- Flexibility to change gateway technology in future
- Clear separation between "what" (intermediate) and "how" (KrakenD)

**Intermediate Format Must Capture**:

| Element | Purpose |
|---------|---------|
| Endpoint definitions | Partner path, method, description |
| Backend mappings | Internal service URL for each endpoint |
| Timeout settings | Per-endpoint or per-backend timeouts |
| Rate limit tier | Which rate limit applies to this endpoint |
| Transformation reference | Which version transformations apply |
| Aggregation config | For fan-out endpoints: list of backends, merge strategy |
| Authentication requirement | Whether endpoint requires JWT (most do) |

**Example Intermediate Format** (conceptual):

```yaml
endpoints:
  - path: /partner/orders/{id}
    method: GET
    backend:
      service: orders-service.orders.svc.cluster.local
      port: 8080
      path: /api/v3/orders/{id}
    timeout: 5s
    rate_limit_tier: standard
    auth_required: true
    transformations:
      - version_before: "2024-06-01"
        ref: orders-v1-transform
        
  - path: /partner/dashboard/{partner_id}
    method: GET
    aggregation:
      strategy: parallel
      backends:
        - service: users-service.users.svc.cluster.local
          path: /api/users/{partner_id}
          response_key: user
        - service: orders-service.orders.svc.cluster.local
          path: /api/orders/recent?partner={partner_id}
          response_key: recent_orders
    timeout: 3s
    rate_limit_tier: standard
    auth_required: true
```

---

### Pipeline Triggers

The pipeline runs when relevant changes occur.

**Why this matters**: The Partner API spec should stay current with internal team changes, but unnecessary runs waste resources.

**Trigger Conditions**:

| Trigger | Source | Action |
|---------|--------|--------|
| Internal team spec change | Webhook from internal team repo | Fetch that team's spec, run full pipeline |
| Manifest change | Merge to governance repo | Run full pipeline |
| Transformation rule change | Merge to governance repo | Run full pipeline |
| Manual trigger | Platform team action | Run full pipeline |

**Webhook Security**:

Webhooks from internal team repos must be verified to prevent spoofing:

- Each internal team repo is configured with a **webhook secret** (shared secret token)
- GitLab includes the secret in the `X-Gitlab-Token` header
- Pipeline verifies the header matches the expected secret before processing
- Webhooks with invalid or missing tokens are rejected

This uses GitLab's native webhook secret token support, which is the standard approach for webhook authentication.

**Webhook Processing Flow**:

1. Webhook arrives with `X-Gitlab-Token` header
2. Pipeline verifies token matches configured secret for that repo
3. Pipeline identifies which team's spec changed from webhook payload
4. Pipeline fetches **only** the changed team's spec
5. Pipeline compares against stored spec in `specs/{team}/openapi.yaml`
6. If changed, updates stored spec in governance repo
7. Pipeline reads all other specs from governance repo
8. Full transformation, validation, and merge runs
9. If breaking change detected, pipeline fails and notifies via Slack

---

### Changelog Generation

The pipeline auto-generates detailed changelogs when specs change.

**Why this matters**: Partners need to know what changed between versions to plan their integrations and upgrades.

**Changelog Detail Level**:

The changelog includes **field-level detail**, not just endpoint-level summaries. This provides partners with precise information about what changed.

**Example Changelog Output**:

```markdown
## Changes detected on 2025-01-29

### Added Endpoints
- `GET /partner/invoices/{id}` - Retrieve invoice by ID
- `GET /partner/invoices` - List invoices with pagination

### Modified Endpoints

#### GET /partner/orders/{id}
Response changes:
- Added field: `metadata` (object, optional)
- Added field: `shipping_carrier` (string, optional)

#### POST /partner/orders
Request changes:
- Added field: `priority` (string, optional) - Order priority level

### Removed Endpoints
None

### Schema Changes
- Added schema: `Invoice`
- Added schema: `InvoiceLineItem`
- Modified schema: `Order` - added `metadata` and `shipping_carrier` fields
```

**Generation Approach**:
- Use oasdiff or similar tool to compare previous and new Partner API spec
- Generate detailed diff including endpoint and field-level changes
- Format as Markdown for inclusion in Partner Portal
- Store in `changelog/latest.md` and archive previous changelogs

---

## Pipeline Outputs

The pipeline produces several artifacts:

| Output | Consumer | Purpose |
|--------|----------|---------|
| `specs/{team}/openapi.yaml` | Pipeline itself, audit | Stored copies of each team's spec |
| `output/partner-api.yaml` | Partner Portal, external tools | The unified Partner API OpenAPI specification |
| `output/gateway-config.yaml` | Gateway config converter | Intermediate gateway configuration |
| `changelog/latest.md` | Partner Portal | Auto-generated detailed changelog |
| Validation report | CI logs, notifications | Details of any validation failures |

---

## Integration Points

### From Internal Team Repositories

The pipeline reads from internal team repos:
- OpenAPI specification files (read-only access)
- Triggered by webhooks on spec changes (verified via webhook secrets)

### From Governance Repository

The pipeline is configured by governance repo contents:
- Manifest defining sources
- Spec mapping rules per source
- Stored team specs (for incremental updates)
- Shared schema definitions
- Spectral linting rules
- Previous Partner API spec (for breaking change comparison)

### To Gateway

The pipeline produces configuration consumed by gateway deployment:
- Intermediate config converted to KrakenD config
- Version transformation references (for runtime response transformations)

### To Partner Portal

The pipeline produces documentation source:
- Partner API OpenAPI spec (for API reference generation)
- Detailed changelog content

### To Internal Teams (Feedback)

The pipeline notifies teams of issues:
- Validation failures (Slack, GitLab comments)
- Breaking change detection (triggers coordination workflow)

---

## Failure Handling

| Failure Scenario | Pipeline Behavior |
|------------------|-------------------|
| Webhook token invalid | Reject webhook; log security event; do not process |
| Cannot fetch internal spec | Fail pipeline; notify internal team; do not update Partner API |
| Internal spec invalid OpenAPI | Fail pipeline; notify internal team with validation errors |
| Internal spec fails linting | Fail pipeline; notify internal team with linting errors |
| Spec mapping produces invalid spec | Fail pipeline; notify platform team (mapping rule bug) |
| Merge conflict (duplicate paths) | Fail pipeline; notify platform team (manifest/mapping issue) |
| Breaking change detected | Fail pipeline; trigger Breaking Change Lifecycle workflow |
| Gateway config generation fails | Fail pipeline; notify platform team |

**Principle**: The pipeline should never silently produce a degraded Partner API. Failures are explicit and block updates until resolved.

---

## Relationship to Contract Testing

The spec aggregation pipeline and contract testing are complementary:

| Concern | Handled By |
|---------|------------|
| Is the internal spec well-formed and compliant? | Spec aggregation pipeline (Spectral linting) |
| Would internal changes break Partner API? | Spec aggregation pipeline (oasdiff) |
| Does the internal service actually match its spec? | Contract testing in internal team CI (Schemathesis) |

The pipeline ensures specs are correct; contract testing ensures implementations match specs.

---

## Success Criteria

The spec aggregation pipeline is complete when:

| Criteria | Validation |
|----------|------------|
| Stores team specs in governance repo | Each team's spec stored in `specs/{team}/openapi.yaml` |
| Fetches only changed specs | Webhook triggers fetch of single team's spec, not all |
| Verifies webhook authenticity | Invalid webhook tokens are rejected |
| Applies spec mappings correctly | Mapped specs have correct paths, filtered fields, renamed schemas |
| Produces valid merged spec | `partner-api.yaml` passes OpenAPI validation |
| Detects breaking changes | Intentional breaking change is caught and blocks pipeline |
| Generates detailed changelog | Changelog includes field-level changes, not just endpoint list |
| Generates gateway config | Intermediate config contains all endpoints with correct backend mappings |
| Notifies on failures | Validation failures result in Slack notification to appropriate team |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 | 2026-01-28 | Initial draft |
| 0.2 | 2026-01-29 | Resolved open questions: webhook security via secrets; team specs stored in governance repo (no separate caching); YAML transformations with sidecar escape hatch; detailed field-level changelog generation |
| 0.3 | 2026-01-30 | Terminology update: renamed "transformations" to "spec mappings" to distinguish from runtime version transformations; updated to reference canonical repository structure in 03-governance-workflow.md |
