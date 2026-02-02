# 07 - Versioning Strategy

## Purpose

The Versioning Strategy defines how the Partner API evolves over time while maintaining stability for existing partners. It covers version identification, lifecycle management, transformation implementation, and the operational processes around version changes.

This plan operationalizes the versioning concepts from the root PLAN.md, providing concrete guidance for implementation.

---

## Role in the Platform

Versioning serves as the **stability and evolution layer**:

| Need | How Versioning Addresses It |
|------|----------------------------|
| Partners need stability | Released versions never change behavior |
| Internal services need to evolve | Transformations absorb internal changes |
| Platform needs to move forward | New versions introduce improvements |
| Old versions need retirement | Lifecycle management with clear timelines |

Versioning does NOT:
- Prevent internal teams from making changes (it absorbs them)
- Require partners to upgrade immediately (grace periods provided)
- Version individual endpoints (whole API versioned together)

---

## Version Identification

### Version Format

Partner API versions use date-based identifiers: `YYYY-MM-DD`

**Examples**: `2024-01-15`, `2024-06-01`, `2025-01-15`

**Why date-based**:
- Clear chronological ordering
- No ambiguity about major/minor/patch semantics
- Industry precedent (Stripe, Twilio)
- Communicates when the version was released

**No patch versions**: The versioning scheme uses dated versions only. There are no patch versions (e.g., no `2024-01-15.1`). Bug fixes that are non-breaking are applied to the current version without changing the version identifier. Breaking bug fixes are included in the next dated version release or trigger an early release if critical.

### How Partners Specify Version

Partners include their desired API version in request headers:

```
GET /partner/orders/123
X-Partner-API-Version: 2024-01-15
Authorization: Bearer <token>
```

**If header is omitted**: Partner receives their default version (set during onboarding, stored in OAuth client configuration).

**If header specifies unknown version**: Gateway returns 400 Bad Request with message listing available versions.

**If header specifies retired version**: Gateway returns 410 Gone with message indicating version is no longer available and suggesting upgrade.

### Version in Responses

Every response includes the version that was used:

```
HTTP/1.1 200 OK
X-Partner-API-Version: 2024-01-15
Content-Type: application/json

{"order_id": "123", ...}
```

This confirms which version processed the request, useful for debugging.

---

## Version Lifecycle

Each Partner API version progresses through defined states:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│   ACTIVE ──────► SUPPORTED ──────► DEPRECATED ──────► RETIRED              │
│                                                                             │
│   (current)      (still works)     (works + warnings)  (rejected)          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### State Definitions

| State | Description | Gateway Behavior |
|-------|-------------|------------------|
| **Active** | Current recommended version for new integrations | Normal processing |
| **Supported** | Previous versions that work fully | Normal processing |
| **Deprecated** | Announced end-of-life; still functional | Process + deprecation warning header |
| **Retired** | No longer available | Reject with 410 Gone |

### Lifecycle Timing

| Transition | Typical Timeline | Trigger |
|------------|------------------|---------|
| Active → Supported | When new version released | New version becomes Active |
| Supported → Deprecated | 12-18 months after release | Platform team decision |
| Deprecated → Retired | 6 months after deprecation | Scheduled date |

**Minimum guarantees**:
- A version remains Active for at least 6 months
- A version remains Supported for at least 12 months
- Deprecation notice given at least 6 months before retirement

**No emergency retirement process**: All versions, including those with security issues, follow the standard deprecation timeline. This ensures consistency and predictability for partners. If a security issue is discovered, the standard 6-month deprecation notice applies. Partners rely on this predictability for planning their integration work.

### Deprecation Warnings

When a partner uses a deprecated version, responses include warning headers:

```
HTTP/1.1 200 OK
X-Partner-API-Version: 2024-01-15
Deprecation: true
Sunset: Sat, 15 Jun 2025 00:00:00 GMT
Link: </docs/migration/2024-01-15-to-2024-06-01>; rel="deprecation"
```

| Header | Purpose |
|--------|---------|
| `Deprecation` | Indicates version is deprecated (RFC 8594) |
| `Sunset` | Date when version will be retired (RFC 8594) |
| `Link` | URL to migration guide |

These headers allow partners to programmatically detect deprecation even if they don't check emails or visit the portal.

---

## Partner Deprecation Notification

Partners are notified of upcoming deprecations through multiple channels to ensure they don't miss important lifecycle changes.

### Notification Channels

| Channel | When Used | Details |
|---------|-----------|---------|
| **API Response Headers** | Every request to deprecated version | `Deprecation` and `Sunset` headers included automatically |
| **Partner Portal** | Persistent banner | Warning displayed when logged in; shows affected applications |
| **Email** | At deprecation announcement and periodic reminders | Sent to registered partner contacts |

### Notification Timeline

| Time Before Retirement | Notification Actions |
|------------------------|---------------------|
| 6 months (deprecation start) | Email announcement; portal banner activated; API headers begin |
| 3 months | Reminder email to partners still using deprecated version |
| 1 month | Final reminder email; increased portal banner prominence |
| Retirement date | Version returns 410 Gone; final notification sent |

### Email Content

Deprecation emails include:
- Which version is being deprecated
- Retirement date
- Link to migration guide
- List of partner's applications using the deprecated version
- Contact information for support

---

## What Triggers a New Version

Not every change requires a new Partner API version. See PLAN.md for the canonical definition of breaking vs. non-breaking changes.

### Changes That DON'T Require New Version

Non-breaking changes flow through to all versions automatically:
- New optional response fields
- New endpoints
- New optional request parameters
- Performance improvements
- Bug fixes aligning with documented contract

### Changes That DO Require New Version

Breaking changes require a new version:
- Removed or renamed response fields
- Removed or renamed endpoints
- Changed field types
- Added required request fields
- Changed default behavior

When these changes are needed, they're introduced in a new version while transformations maintain the old behavior for previous versions.

---

## Transformation Implementation

### Transformation Architecture

When a partner requests an old version, the gateway transforms the current internal response to match the old version's format.

```
Internal Service     Gateway                              Partner (v2024-01-15)
      │                 │                                        │
      │◄── Request ─────┤◄────────── Request ────────────────────┤
      │                 │            X-Partner-API-Version:      │
      │                 │            2024-01-15                  │
      │                 │                                        │
      ├── Response ────►│                                        │
      │   (current      ├── Apply transformation ───►            │
      │    format)      │   (current → 2024-01-15)               │
      │                 │                                        │
      │                 ├── Transformed Response ───────────────►│
      │                 │   (2024-01-15 format)                  │
```

### Transformation Storage

Transformations are stored in the governance repository (`partner-api-governance`):

```
partner-api-governance/
└── versions/
    ├── transformations/
    │   ├── 2024-06-01/           # Transformations for partners before 2024-06-01
    │   │   ├── orders.yaml       # Orders domain transformations (KrakenD config)
    │   │   ├── users.yaml        # Users domain transformations
    │   │   └── manifest.yaml     # What this version changes
    │   │
    │   └── 2025-01-15/           # Transformations for partners before 2025-01-15
    │       ├── orders.yaml
    │       └── manifest.yaml
    │
    ├── golden-files/             # Expected responses for regression testing
    │   ├── 2024-01-15/
    │   │   ├── orders-get.json
    │   │   └── users-list.json
    │   └── 2024-06-01/
    │       ├── orders-get.json
    │       └── users-list.json
    │
    └── changelog/
        ├── 2024-06-01.md         # What changed in this version
        └── 2025-01-15.md
```

### Transformation Manifest

Each version's manifest documents what changed and what transformations apply:

```yaml
# versions/transformations/2024-06-01/manifest.yaml
version: "2024-06-01"
description: "Restructured customer data, standardized timestamps"

changes:
  - endpoint: /partner/orders/*
    type: response_restructure
    description: "Flattened nested customer object"
    transformation: orders.yaml
    
  - endpoint: /partner/users/*
    type: field_rename
    description: "Renamed createdAt to created_at"
    transformation: users.yaml

applies_to_versions_before: "2024-06-01"
```

### Transformation Chaining

Partners on very old versions may need multiple transformations applied in sequence:

```
Partner on 2024-01-15 requests /partner/orders/123

Current response (as of 2025-01-15):
{
  "order_id": "123",
  "customer": {"id": "c1", "name": "John"},
  "amount": {"value": 99.99, "currency": "USD"}
}

Apply transformation 2025-01-15 (if partner version < 2025-01-15):
{
  "order_id": "123",
  "customer": {"id": "c1", "name": "John"},
  "total": 99.99,                              // amount restructured
  "currency": "USD"
}

Apply transformation 2024-06-01 (if partner version < 2024-06-01):
{
  "order_id": "123",
  "customer_id": "c1",                         // customer flattened
  "customer_name": "John",
  "total": 99.99,
  "currency": "USD"
}

Partner receives final response matching 2024-01-15 contract.
```

### Transformation Performance

Transformation chains add latency proportional to chain depth. Each transformation typically adds sub-millisecond to single-digit milliseconds of processing time.

**Approach**: Accept linear overhead and optimize only if needed.

| Guideline | Details |
|-----------|---------|
| Default behavior | Apply transformations sequentially in chain |
| Monitoring | Track transformation latency in NewRelic per version |
| Optimization trigger | Consider optimization if transformation adds >10ms |
| Optimization options | Pre-compile chains, direct transformations for old versions |

For most use cases, the linear overhead is negligible. If monitoring reveals performance issues with deep chains, optimization strategies can be applied selectively.

### Implementation Options

Transformations follow the same two-tier model as spec transformations and gateway transformations:

| Option | When to Use |
|--------|-------------|
| **KrakenD declarative configuration** | First choice. Use for field renames, restructuring, filtering. |
| **Sidecar service** | When declarative configuration cannot express the transformation. |

---

## Breaking Change Coordination

When an internal team needs to make a breaking change, the platform team creates a transformation.

### Workflow

1. **Internal team submits request** (JIRA ticket)
   - What's changing
   - Why it's needed
   - Expected timeline
   - Example current and previous format responses

2. **Platform team assesses impact**
   - Which Partner API versions affected
   - Complexity of transformation

3. **Platform team implements transformation**
   - Write transformation (KrakenD config or sidecar)
   - Test with provided examples
   - Test with golden files for all affected versions

4. **Platform team deploys transformation**
   - Merge to governance repo
   - ArgoCD deploys to gateway
   - Verify in staging

5. **Platform team unblocks internal team**
   - Comment on JIRA ticket
   - Internal team can proceed with their deployment

### Testing Transformations

Each transformation must pass multiple test layers:

| Test Type | Purpose | How |
|-----------|---------|-----|
| Unit tests | Verify transformation logic | Test with example inputs/outputs |
| Integration tests | Verify against real service responses | Call staging service, apply transformation |
| Golden file tests | Ensure all versions still work | Compare against recorded expected responses |

### Golden File Testing

Golden file testing ensures all supported versions continue working as internal services evolve.

**How it works**:

1. **Record**: For each supported version, record expected responses for key endpoints
   - Store in `versions/golden-files/{version}/`
   - Include representative requests (GET order, list users, etc.)

2. **Replay**: Before each release, replay requests and compare against golden files
   - Make request through gateway for each supported version
   - Compare response against stored golden file
   - Fail if response doesn't match expected format

3. **Update**: When intentionally changing a version's behavior
   - Update golden file to reflect new expected response
   - Requires explicit approval in PR

**Golden file structure**:

```json
// versions/golden-files/2024-01-15/orders-get.json
{
  "request": {
    "method": "GET",
    "path": "/partner/orders/test-123",
    "version": "2024-01-15"
  },
  "expected_response": {
    "status": 200,
    "body": {
      "order_id": "test-123",
      "customer_id": "c1",
      "customer_name": "Test Customer",
      "total": 99.99,
      "currency": "USD"
    }
  }
}
```

**When golden file tests run**:
- On every PR to governance repo
- Before deploying new transformations
- After internal service deployments (triggered via webhook)

---

## Version Release Process

When accumulated changes warrant a new Partner API version:

### Decision Criteria

Release a new version when:
- Significant new capabilities added
- Multiple beneficial changes accumulated
- Major internal service upgrades completed
- Scheduled release cycle (e.g., quarterly)

### Release Steps

1. **Decide on version date**: Choose release date (becomes version identifier)

2. **Finalize transformations**: Ensure all transformations for changes in this version are complete and tested

3. **Update golden files**: Create golden files for the new version

4. **Update documentation**:
   - Changelog describing what's new/changed
   - Migration guide from previous version
   - Updated API reference

5. **Update version configuration**:
   - Add new version to gateway config
   - Set new version as Active
   - Move previous Active version to Supported

6. **Deploy**:
   - Merge changes to governance repo
   - ArgoCD deploys updates
   - Verify in staging
   - Run golden file tests for all versions

7. **Announce** (via all channels):
   - Notification banner in Partner Portal
   - Email to all partners
   - Slack announcement in partner channel (if exists)

---

## Partner Migration Support

### Migration Guides

For each version transition, provide a migration guide:

```markdown
# Migrating from 2024-01-15 to 2024-06-01

## Summary
This version restructures customer data and standardizes timestamp formats.

## Breaking Changes

### Customer data restructure
Before (2024-01-15):
  "customer_id": "c1",
  "customer_name": "John Doe"

After (2024-06-01):
  "customer": {
    "id": "c1",
    "name": "John Doe"
  }

### Timestamp format
Before: "created_at": "2024-01-15T10:30:00"
After: "created_at": "2024-01-15T10:30:00Z" (explicit UTC)

## Migration Steps
1. Update response parsing to handle nested customer object
2. Update timestamp parsing to handle Z suffix
3. Test with sandbox environment
4. Update X-Partner-API-Version header to 2024-06-01
```

### Sandbox Testing

Partners can test new versions before committing:
- Request any version via header, regardless of their default
- Sandbox environment available for testing (same as production API, separate data)
- No changes to their credentials needed

### Version Pinning

Partners control their version:
- Default version stored in OAuth client config
- Can override per-request via header
- Platform team can update default on partner request
- Self-service portal allows changing default version

---

## Monitoring Version Usage

### Metrics to Track

| Metric | Purpose |
|--------|---------|
| Requests per version | Understand adoption |
| Partners per version | Identify migration progress |
| Deprecated version usage | Target outreach for migration |
| Version in errors | Identify version-specific issues |
| Transformation latency | Monitor performance overhead |

### Dashboard Views

- **Version distribution**: Pie chart of requests by version
- **Migration progress**: Line chart of deprecated version usage over time
- **Partner version list**: Table of partners and their versions
- **Transformation performance**: Latency percentiles per version

### Alerting

| Condition | Action |
|-----------|--------|
| Deprecated version >10% traffic | Review migration outreach |
| Partner stuck on deprecated for >3 months | Direct outreach |
| Retired version requests | Should be zero; investigate if not |
| Transformation latency >10ms | Investigate optimization |

---

## Edge Cases and Error Handling

### Unknown Version Requested

```
Request: X-Partner-API-Version: 2099-01-01

Response:
HTTP/1.1 400 Bad Request
{
  "error": "unknown_version",
  "message": "Version 2099-01-01 is not recognized",
  "available_versions": ["2024-01-15", "2024-06-01", "2025-01-15"],
  "current_version": "2025-01-15"
}
```

### Retired Version Requested

```
Request: X-Partner-API-Version: 2023-06-01

Response:
HTTP/1.1 410 Gone
{
  "error": "version_retired",
  "message": "Version 2023-06-01 is no longer available",
  "retired_on": "2024-06-01",
  "minimum_supported_version": "2024-01-15",
  "migration_guide": "https://docs.partner-api.com/migration/2023-06-01"
}
```

### Transformation Failure

If a transformation fails (bug, unexpected data):

```
Response:
HTTP/1.1 500 Internal Server Error
{
  "error": "transformation_error",
  "message": "Unable to process response for requested version",
  "request_id": "req-abc123",
  "trace_id": "trace-xyz789"
}
```

Log includes full details for debugging. Platform team alerted.

---

## Integration Points

### With Gateway

- Gateway reads version from header
- Gateway applies transformation chain
- Gateway returns version and deprecation headers
- Gateway tracks transformation latency

### With Governance Repository

- Transformation code stored in governance repo
- Version configuration in governance repo
- Golden files for regression testing in governance repo
- Changelog and migration guides in governance repo

### With Partner Portal

- Portal displays documentation per version
- Version switcher in API reference
- Migration guides accessible
- Changelog displayed
- Deprecation warnings shown for affected applications
- Partners can see their current version in dashboard
- Partners can change their default version
- Partners can see version lifecycle status

### With Breaking Change Lifecycle

- New transformations created through Breaking Change Lifecycle
- Transformation deployment unblocks internal teams

---

## Success Criteria

Versioning strategy is complete when:

| Criteria | Validation |
|----------|------------|
| Version header respected | Different versions return different formats |
| Transformations work | Old version request receives old format |
| Transformation chain works | Very old version gets multiple transformations |
| Deprecation warnings sent | Deprecated version returns `Deprecation` and `Sunset` headers |
| Deprecation emails sent | Partners receive email notifications |
| Portal shows deprecation | Logged-in partners see deprecation banner |
| Retired versions rejected | Retired version returns 410 Gone |
| Documentation versioned | Portal shows docs per version |
| Migration guides available | Each version transition has guide |
| Golden file tests pass | All supported versions return expected responses |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 | 2026-01-28 | Initial draft |
| 0.2 | 2026-01-30 | Resolved open questions: dated versions only (no patches); accept linear transformation overhead; golden file testing for version regression; portal + email + API headers for notifications; no emergency retirement process |
| 0.3 | 2026-01-30 | Merged Self-Service Portal into Partner Portal integration; simplified breaking change lists to reference PLAN.md; standardized repository naming |
