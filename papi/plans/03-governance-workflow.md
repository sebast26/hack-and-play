# 03 - Governance Workflow

## Purpose

The Governance Workflow defines how changes flow into the Partner API—who can change what, what approvals are required, and how internal teams contribute without needing to understand the platform's internals.

The goal is a system where internal teams focus on their own APIs while the platform team maintains control over what partners see.

---

## Role in the Platform

The governance workflow serves as the **control and coordination layer**:

| Concern | How It's Addressed |
|---------|-------------------|
| What endpoints are exposed to partners | Manifest in governance repo (platform team controlled) |
| How internal APIs are transformed | Transformation rules in governance repo (platform team controlled) |
| Who approves changes | Approval levels based on change type |
| How internal teams contribute | Via their own repos; no governance repo knowledge needed |

The governance workflow does NOT:
- Require internal teams to access the governance repository
- Require internal teams to understand transformation rules
- Block internal teams from evolving their internal APIs (only partner-breaking changes require coordination)

---

## Design Principle: Invisible Governance

Internal teams should be able to work as if the Partner API doesn't exist—until they make a change that would break it.

**For internal teams, the experience is**:
1. They maintain their OpenAPI spec in their own repo
2. They tag endpoints intended for partners (e.g., `x-partner: true` or tag `partner`)
3. They push changes normally
4. If a change would break Partner API, their CI fails with clear instructions
5. Otherwise, their changes flow through automatically

**They don't need to know**:
- That a governance repository exists
- How transformation rules work
- How the aggregation pipeline operates
- How gateway configuration is generated

This reduces cognitive load and lets teams focus on their domain.

---

## Governance Repository Structure

The governance repository (`partner-api-governance`) is owned and operated by the platform team. It contains everything needed to build the Partner API.

**Complete Repository Structure**:

```
partner-api-governance/
├── manifest.yaml                    # Defines which internal specs to include
├── spec-mappings/                   # Rules for transforming internal specs → Partner API spec
│   ├── orders-mapping.yaml          # Orders team spec mapping rules
│   ├── users-mapping.yaml           # Users team spec mapping rules
│   └── payments-mapping.yaml        # Payments team spec mapping rules
├── specs/                           # Stored copies of internal team specs (auto-updated)
│   ├── orders/
│   │   └── openapi.yaml
│   ├── users/
│   │   └── openapi.yaml
│   └── payments/
│       └── openapi.yaml
├── shared-schemas/                  # Common schemas used across Partner API
│   └── common.yaml
├── versions/                        # Version-specific transformations and testing
│   ├── transformations/             # Runtime response transformations per version
│   │   ├── 2024-06-01/
│   │   │   ├── orders.yaml          # KrakenD config for orders transformations
│   │   │   └── manifest.yaml        # What this version changes
│   │   └── 2025-01-15/
│   │       └── ...
│   ├── golden-files/                # Expected responses for regression testing
│   │   ├── 2024-01-15/
│   │   │   └── orders-get.json
│   │   └── 2024-06-01/
│   │       └── orders-get.json
│   └── changelog/                   # Per-version changelog entries
│       ├── 2024-06-01.md
│       └── 2025-01-15.md
├── spectral/                        # Linting rules for API standards
│   └── .spectral.yaml
├── gateway/                         # Gateway-specific configuration templates
│   └── krakend-template.json
├── output/                          # Generated artifacts (gitignored or CI-only)
│   ├── partner-api.yaml             # Merged Partner API specification
│   └── krakend.json                 # Generated gateway configuration
├── ci-templates/                    # Reusable CI templates for internal teams
│   ├── spec-validation.yml
│   ├── breaking-change-detection.yml
│   ├── implementation-testing.yml
│   └── partner-api-checks.yml
├── docs/                            # Internal documentation for platform team
│   └── onboarding-guide.md
└── .gitlab-ci.yml                   # Pipeline definitions
```

**Two Types of Transformations**:

| Type | Purpose | Location |
|------|---------|----------|
| **Spec mappings** | Transform internal OpenAPI specs into Partner API spec at build time (path remapping, field filtering, endpoint selection) | `spec-mappings/` |
| **Version transformations** | Transform API responses at runtime for partners using older API versions | `versions/transformations/` |

**Repository Contents Summary**:

| Directory/File | Purpose | Who Modifies |
|----------------|---------|--------------|
| `manifest.yaml` | Defines which internal specs to include | Platform team |
| `spec-mappings/` | Spec-to-spec transformation rules per internal team | Platform team |
| `specs/` | Stored copies of internal team specs | Pipeline (automated) |
| `shared-schemas/` | Common schemas used across Partner API | Platform team |
| `versions/` | Version transformations, golden files, changelogs | Platform team |
| `spectral/` | Linting rules for API standards | Platform team |
| `gateway/` | Gateway-specific configuration templates | Platform team |
| `output/` | Generated artifacts | Pipeline (automated) |
| `ci-templates/` | Reusable CI templates for internal teams | Platform team |
| `docs/` | Internal documentation for platform team | Platform team |
| `.gitlab-ci.yml` | Pipeline definitions | Platform team |

**Internal teams never directly modify this repository.**

---

## Manifest Structure

The manifest defines which internal team specs are included and how they're mapped to the Partner API. Each source entry includes ownership metadata to enable team-level communication and tracking.

**Manifest Fields**:

| Field | Purpose |
|-------|---------|
| `name` | Unique identifier for this source |
| `repository` | Where to fetch the spec |
| `path` | Location of OpenAPI spec file within repo |
| `ref` | Branch, tag, or commit to use |
| `spec_mapping` | Which spec mapping rules to apply |
| `enabled` | Ability to temporarily disable a source |
| `owner_team` | Team responsible for this source |
| `owner_contact` | Contact email or Slack channel for notifications |

**Example Manifest Structure**:

```yaml
sources:
  - name: orders-core
    repository: https://gitlab.company.com/teams/orders-service
    path: api/openapi.yaml
    ref: main
    spec_mapping: orders-mapping.yaml
    enabled: true
    owner_team: orders
    owner_contact: orders-team@company.com

  - name: orders-fulfillment
    repository: https://gitlab.company.com/teams/orders-fulfillment-service
    path: api/openapi.yaml
    ref: main
    spec_mapping: orders-fulfillment-mapping.yaml
    enabled: true
    owner_team: orders
    owner_contact: orders-team@company.com
    
  - name: users
    repository: https://gitlab.company.com/teams/users-service
    path: docs/api/spec.yaml
    ref: release-v2
    spec_mapping: users-mapping.yaml
    enabled: true
    owner_team: users
    owner_contact: "#users-team"
```

This flat structure (one entry per repo) with ownership metadata allows:
- Teams with multiple repos to have separate entries for each
- Querying all sources owned by a specific team
- Sending notifications to the right team when issues occur

---

## How Internal Teams Contribute

### Tagging Endpoints for Partner Exposure

Internal teams indicate which endpoints should be exposed to partners by tagging them in their OpenAPI spec:

**Option A: Using OpenAPI tags**
```yaml
paths:
  /orders/{id}:
    get:
      tags:
        - partner  # This endpoint is exposed to partners
        - orders
      summary: Get order by ID
```

**Option B: Using vendor extension**
```yaml
paths:
  /orders/{id}:
    get:
      x-partner-api: true  # This endpoint is exposed to partners
      summary: Get order by ID
```

The transformation rules in the governance repo filter based on these markers.

### Requesting New Endpoint Exposure

When an internal team wants to expose a new endpoint to partners:

1. **Team adds the endpoint** to their OpenAPI spec with appropriate partner tag
2. **Webhook triggers pipeline** which detects new endpoint
3. **Pipeline checks** if transformation rules exist for this endpoint's path
4. **If no transformation exists**, pipeline notifies platform team (new endpoint needs onboarding)
5. **Platform team reviews**, creates/updates transformation rules if needed
6. **Next pipeline run** includes the new endpoint in Partner API

For teams already onboarded with broad transformation rules (e.g., "include all endpoints tagged `partner`"), new endpoints may flow through automatically without platform team intervention.

### Breaking Change Flow

When an internal team makes a change that would break the Partner API:

1. **Team pushes change** to their repo
2. **Their CI runs** oasdiff against Partner API spec (provided via CI template)
3. **CI fails** with message: "Breaking change detected. See coordination workflow."
4. **Team posts** in `#partner-api-changes` Slack channel
5. **Team creates** JIRA ticket with change details
6. **Platform team** implements transformation to preserve backward compatibility
7. **Platform team** deploys transformation
8. **Platform team** notifies internal team they can proceed
9. **Team merges** their change

This flow is defined in the Breaking Change Lifecycle (see PLAN.md Versioning Strategy section).

---

## Platform Team Workflows

### Onboarding a New Internal Team

When a new internal team wants to contribute to the Partner API:

**Steps**:

1. **Initial contact**: Team reaches out via Slack or JIRA
2. **Self-service preparation**: Team reviews onboarding guide and runs dry-run validation locally
3. **Spec review**: Platform team reviews their OpenAPI spec for quality and standards
4. **Transformation design**: Platform team creates transformation rules for the team
5. **Manifest update**: Platform team adds team to manifest (with `owner_team` and `owner_contact`)
6. **Webhook setup**: Platform team works with team to configure webhook from their repo
7. **CI template integration**: Team adds contract testing CI template to their pipeline
8. **Validation**: Pipeline runs, team's endpoints appear in Partner API spec
9. **Communication**: Platform team confirms onboarding complete

**Approval required**: Platform team lead approves manifest changes for new teams.

### Onboarding Documentation and Tooling

To enable internal teams to self-serve as much as possible, the platform provides:

**Comprehensive Onboarding Guide** (in `docs/onboarding/`):
- How to tag endpoints for partner exposure (`partner` tag or `x-partner-api` extension)
- How to add CI templates to their pipeline
- What to expect during platform team review
- Common issues and how to resolve them
- Examples from already-onboarded teams

**Dry-Run Validation Tool**:
- Teams can run validation locally before requesting onboarding
- Checks OpenAPI spec validity and Spectral linting
- Simulates transformation to show what would be exposed
- Identifies potential issues before platform team review

**What Remains Manual**:
- Adding entries to manifest (requires platform team review and approval)
- Creating transformation rules (platform team responsibility)
- Webhook secret configuration (requires coordination)

This approach enables teams to prepare independently while keeping manifest control with the platform team.

### Modifying Transformation Rules

When transformation rules need to change (new fields to filter, path changes, etc.):

**Steps**:

1. **Identify need**: Usually triggered by internal team request or Partner API design review
2. **Draft changes**: Platform team member modifies transformation rules
3. **Test locally**: Run pipeline locally to verify transformed output
4. **Create MR**: Submit merge request with changes
5. **Review**: Another platform team member reviews
6. **Merge**: After approval, merge triggers pipeline
7. **Verify**: Check Partner API spec reflects expected changes

**Approval required**: Any platform team member can approve transformation changes.

### Releasing a New Partner API Version

When accumulated changes warrant a new Partner API version:

**Steps**:

1. **Decision**: Platform team decides current changes justify new version
2. **Version bump**: Update Partner API version in manifest metadata
3. **Changelog**: Document changes for this version
4. **Transformation updates**: Ensure version transformations are in place for previous versions
5. **Create MR**: Submit merge request with version release
6. **Review**: Platform team lead reviews
7. **Merge**: Triggers pipeline to generate new version
8. **Documentation**: Partner Portal automatically updates with new version docs
9. **Communication**: Notify partners of new version availability

**Approval required**: Platform team lead approves version releases.

### Deprecating and Retiring Versions

When an old Partner API version reaches end of life:

**Steps**:

1. **Announcement**: Communicate deprecation timeline to partners (6+ months notice)
2. **Deprecation flag**: Mark version as deprecated in configuration
3. **Warning headers**: Gateway returns deprecation warning headers for deprecated versions
4. **Monitor usage**: Track which partners still use deprecated version
5. **Retirement**: After grace period, mark version as retired
6. **Enforcement**: Gateway returns error for retired versions
7. **Cleanup**: Remove version-specific transformations (optional, can keep for reference)

**Approval required**: Platform team lead approves deprecation and retirement decisions.

---

## Approval Levels

Different changes require different approval levels:

| Change Type | Approver | Rationale |
|-------------|----------|-----------|
| New team onboarding | Platform team lead | Significant commitment; affects platform scope |
| Manifest changes (add/remove sources) | Platform team lead | Controls what's in Partner API |
| Transformation rule changes | Any platform team member | Routine maintenance |
| Spectral/linting rule changes | Any platform team member | Standards evolution |
| Version releases | Platform team lead | Partner-facing impact |
| Version deprecation/retirement | Platform team lead | Partner-facing impact |
| CI template changes (non-breaking) | Any platform team member | Internal tooling |
| CI template changes (breaking) | Platform team lead | Affects all internal teams |
| Documentation updates | Any platform team member | No runtime impact |

**Emergency Changes**: There is no special expedited process for emergency changes. All changes, including urgent fixes, follow the standard approval workflow. This ensures consistency, maintains audit trails, and prevents shortcuts that could introduce additional issues. If a change is truly urgent, team members should be available to review and approve quickly through the normal process.

---

## CI Templates for Internal Teams

The platform team provides reusable CI templates that internal teams include in their pipelines.

### Template Versioning

CI templates use **semantic versioning** to prevent breaking changes from disrupting internal teams:

- Templates are tagged with versions: `v1`, `v2`, `v1.1.0`, etc.
- Teams pin to a **major version** (e.g., `ref: v1`)
- Non-breaking changes (bug fixes, new optional features) are released within the major version
- Breaking changes require a new major version (e.g., `v2`)
- Platform team announces new major versions and provides migration guidance
- Previous major versions are supported for a deprecation period (minimum 6 months)

### What Templates Provide

| Template | Purpose |
|----------|---------|
| `spec-validation.yml` | Validates OpenAPI spec syntax and Spectral linting |
| `breaking-change-detection.yml` | Compares spec changes against Partner API; fails on breaking changes |
| `implementation-testing.yml` | Runs Schemathesis against staging service |

### How Internal Teams Use Templates

Internal teams add a single include to their `.gitlab-ci.yml`:

```yaml
include:
  - project: 'platform/partner-api-governance'
    file: '/ci-templates/partner-api-checks.yml'
    ref: v1  # Pin to major version for stability

variables:
  OPENAPI_SPEC_PATH: "api/openapi.yaml"
  STAGING_SERVICE_URL: "https://orders-staging.internal.company.com"
```

This gives them:
- Spec validation on every MR
- Breaking change detection against Partner API
- Implementation testing against their staging environment

**Teams don't need to understand what these templates do internally**—they just see pass/fail and actionable error messages.

When `v2` is released with breaking changes, teams can migrate on their own schedule within the deprecation period.

---

## Sync Reliability

### Webhook + Scheduled Reconciliation

To ensure the governance repo stays synchronized with internal team specs, the platform uses a dual approach:

**Webhooks (Real-time)**:
- Internal team repos send webhooks when spec files change
- Pipeline processes changes immediately
- Provides near-real-time synchronization

**Scheduled Reconciliation (Safety Net)**:
- A daily scheduled job runs full reconciliation
- Fetches all specs and compares against stored versions in governance repo
- Catches any changes that webhooks might have missed (network issues, service downtime, etc.)
- Updates any specs that are out of sync

This "belt and suspenders" approach ensures:
- Fast updates in normal circumstances (webhooks)
- No changes are permanently missed (daily reconciliation)
- Self-healing if webhook infrastructure has temporary issues

The daily reconciliation job is lightweight—it compares commit SHAs first and only fetches specs that have actually changed.

---

## Communication Channels

| Channel | Purpose | Participants |
|---------|---------|--------------|
| `#partner-api-changes` (Slack) | Breaking change coordination, urgent issues | Platform team, internal teams |
| `#partner-api-announcements` (Slack) | Version releases, deprecations, maintenance, CI template updates | Platform team → all |
| JIRA project | Tracking onboarding, breaking change tickets | Platform team, internal teams |
| GitLab MRs | Code review, approvals | Platform team |

---

## Audit and Traceability

All changes to the Partner API are traceable:

| What | Where It's Recorded |
|------|---------------------|
| Manifest changes | Git history in governance repo |
| Transformation changes | Git history in governance repo |
| Stored team specs | Git history in governance repo (`specs/` directory) |
| Version releases | Git tags + changelog in governance repo |
| Internal team spec changes | Git history in their repos + stored copies in governance repo |
| Pipeline runs | GitLab CI pipeline history |
| Approval decisions | GitLab MR approvals |
| Reconciliation runs | Scheduled pipeline history |

This audit trail supports compliance requirements and incident investigation.

---

## Integration Points

### With Internal Team Repositories

- Webhooks notify governance repo of spec changes
- CI templates are included from governance repo (version-pinned)
- Breaking change detection compares against Partner API spec

### With Spec Aggregation Pipeline

- Pipeline reads manifest and transformation rules from governance repo
- Pipeline triggered by merges to governance repo
- Pipeline triggered by webhooks from internal repos
- Daily reconciliation job ensures sync

### With Gateway

- Gateway configuration generated from governance repo content
- Deployed via ArgoCD watching governance repo outputs

### With Partner Portal

- Partner API spec from governance repo feeds portal generation
- Changelog from governance repo displayed in portal

---

## Success Criteria

The governance workflow is complete when:

| Criteria | Validation |
|----------|------------|
| Internal teams can contribute without governance repo access | Team onboarded using only their own repo + CI templates |
| Breaking changes are caught before merge | Intentional breaking change fails internal team's CI |
| New team onboarding is documented and repeatable | Second team onboards faster than first using self-service guide |
| Dry-run validation works | Teams can validate specs locally before requesting onboarding |
| Approval levels are enforced | GitLab MR rules require appropriate approvers |
| CI template versioning works | Teams pinned to v1 unaffected by v2 release |
| Changes are traceable | Can trace any Partner API change back to source MR |
| Webhook + reconciliation reliable | No spec changes missed for more than 24 hours |
| Communication channels are effective | Breaking change coordination completes without confusion |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 | 2026-01-28 | Initial draft |
| 0.2 | 2026-01-30 | Resolved open questions: webhook + daily reconciliation for reliability; semantic versioning for CI templates; comprehensive onboarding guide with dry-run validation; no special emergency process; flat manifest with owner_team and owner_contact metadata |
| 0.3 | 2026-01-30 | Established canonical repository structure; distinguished spec mappings (build-time) from version transformations (runtime); updated manifest to use spec_mapping field |
