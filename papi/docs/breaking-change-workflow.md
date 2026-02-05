# Breaking Change Workflow

This document describes the process for handling breaking changes to the Partner API.

## What is a Breaking Change?

A breaking change is any modification to the API that could cause existing partner integrations to fail. Examples include:

| Change Type | Example | Breaking? |
|-------------|---------|-----------|
| Remove endpoint | DELETE `/partner/orders` | Yes |
| Remove field | Remove `order_id` from response | Yes |
| Rename field | `order_id` → `id` | Yes |
| Change field type | `total: string` → `total: number` | Yes |
| Make optional required | `customer_id` becomes required | Yes |
| Add required parameter | New required query param | Yes |
| Change response status | `200` → `201` | Yes |
| Add optional field | Add `metadata` to response | No |
| Add optional parameter | Add optional query param | No |
| Add new endpoint | New `POST /partner/orders` | No |

## When Breaking Changes Are Detected

The CI pipeline automatically detects breaking changes by comparing your spec against the baseline Partner API spec. If breaking changes are detected:

1. The pipeline will fail with a detailed report
2. You'll see which changes are breaking and why
3. You must follow this workflow before the change can be merged

## Breaking Change Process

### Step 1: Evaluate the Change

Ask yourself:
- Is this change absolutely necessary?
- Can it be done in a non-breaking way?
- What is the impact on existing partners?

**Non-breaking alternatives:**

| Instead of... | Consider... |
|---------------|-------------|
| Removing a field | Deprecate it first (keep returning it) |
| Renaming a field | Add new field, keep old one |
| Changing field type | Add new field with new type |
| Making field required | Keep it optional |

### Step 2: Post in Slack

If the breaking change is necessary, post in **#partner-api-changes** Slack channel:

```
🚨 Breaking Change Request

Team: [Your Team Name]
Service: [Service Name]
Change Type: [Remove endpoint / Remove field / etc.]

Description:
[Brief description of what you're changing and why]

Affected Endpoints:
- /partner/orders (removing `legacy_id` field)

Justification:
[Why this change is necessary and cannot be avoided]

Proposed Timeline:
[When you need this deployed]
```

### Step 3: Create a JIRA Ticket

Create a JIRA ticket in the **PARTNER-API** project with:

**Title**: `[Breaking Change] <brief description>`

**Description**:
- Team name and contact
- Detailed description of the change
- List of affected endpoints/fields
- Justification for why it's necessary
- Proposed deprecation timeline
- Partner communication plan

**Labels**: `breaking-change`, `partner-api`, `<your-team>`

### Step 4: Platform Team Review

The platform team will:

1. Review the change request within 2 business days
2. Assess impact on existing partners
3. Propose a deprecation timeline if approved
4. May request modifications or alternatives

### Step 5: Partner Communication

If approved, the platform team will coordinate:

1. **Partner notification**: Email partners about upcoming change
2. **Deprecation period**: Typically 30-90 days depending on impact
3. **Documentation update**: Update Partner API docs with deprecation notice
4. **Migration guide**: Create guide for partners to migrate

### Step 6: Implementation

Once the deprecation period ends:

1. Platform team updates the baseline spec
2. Your breaking change can be merged
3. Aggregation pipeline regenerates the Partner API spec
4. Changes deploy to production

## Deprecation Guidelines

### Deprecation Notice Period

| Impact Level | Minimum Notice |
|--------------|----------------|
| Low (unused endpoints) | 14 days |
| Medium (rarely used) | 30 days |
| High (commonly used) | 60-90 days |
| Critical (all partners use) | 90+ days |

### Marking Deprecation in Specs

Use the `deprecated` flag in your OpenAPI spec:

```yaml
/api/v3/orders:
  get:
    deprecated: true
    summary: List orders (DEPRECATED - use /api/v4/orders)
    description: |
      **DEPRECATED**: This endpoint will be removed on 2024-06-01.
      Please migrate to /api/v4/orders.
```

## Emergency Breaking Changes

In rare cases (security vulnerabilities, legal requirements), expedited breaking changes may be necessary:

1. Contact platform team immediately via Slack DM or @platform-oncall
2. Provide security/legal justification
3. Platform team will assess and may approve immediate deployment
4. Partners will be notified as soon as possible

## Frequently Asked Questions

### Q: How long does the approval process take?
A: Initial review within 2 business days. Full approval timeline depends on impact assessment and partner communication needs.

### Q: Can I bypass the process for small changes?
A: No. All breaking changes, regardless of size, must follow this process. Small changes often have unexpected impacts.

### Q: What if a partner doesn't migrate in time?
A: The platform team will reach out to partners who haven't migrated. In extreme cases, the deprecation period may be extended.

### Q: How do I know if my change is breaking?
A: Run the breaking change detection CI template locally:
```bash
oasdiff breaking baseline.yaml ./openapi.yaml --format text --fail-on WARN
```

### Q: Can I add a field as required for a new endpoint?
A: Yes, new endpoints can have required fields. Breaking change rules only apply to existing endpoints.

## Contact

- **Slack**: #partner-api-changes (breaking changes), #partner-api-support (questions)
- **JIRA**: PARTNER-API project
- **Email**: partner-api-platform@example.com
- **Oncall**: @platform-oncall (emergencies only)
