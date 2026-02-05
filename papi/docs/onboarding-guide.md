# Partner API Onboarding Guide

This guide explains how internal teams can expose their API endpoints through the Partner API.

## Overview

The Partner API is a unified gateway that exposes selected endpoints from internal services to external partners. As an internal team, you can contribute endpoints to the Partner API by:

1. Tagging endpoints in your OpenAPI spec
2. Adding CI templates to your pipeline
3. Coordinating with the platform team for onboarding

## Prerequisites

Before you begin, ensure you have:

- [ ] An OpenAPI 3.0+ spec for your service
- [ ] Your spec hosted in a GitLab repository
- [ ] Access to the Partner API platform repository

## Step 1: Tag Your Endpoints

Add the `partner` tag to any endpoint you want to expose through the Partner API.

**Example:**

```yaml
paths:
  /api/v3/orders:
    get:
      tags:
        - orders      # Your service tag
        - partner     # Add this to expose via Partner API
      summary: List all orders
      description: Returns a list of all orders
      operationId: listOrders
      responses:
        "200":
          description: Successful response
```

### Tagging Guidelines

| Tag | Meaning |
|-----|---------|
| `partner` | Endpoint will be exposed to partners |
| `internal` | Endpoint is internal-only (explicitly excluded) |
| `admin` | Admin endpoint (explicitly excluded) |

**Important:**
- Only add `partner` tag to endpoints that are safe and intended for external use
- Do not expose endpoints that return sensitive internal data
- Ensure all partner-tagged endpoints have proper documentation (description, operationId)

## Step 2: Follow API Design Standards

All partner-facing endpoints must follow the Partner API design standards:

### Naming Conventions

- **Property names**: Use `snake_case` (e.g., `order_id`, `created_at`)
- **Operation IDs**: Use camelCase (e.g., `listOrders`, `getOrderById`)

```yaml
# Good
properties:
  order_id:
    type: string
  customer_name:
    type: string
  created_at:
    type: string
    format: date-time

# Bad - will fail validation
properties:
  orderId:       # Should be order_id
  customerName:  # Should be customer_name
  createdAt:     # Should be created_at
```

### Required Documentation

Every partner-tagged endpoint must have:

- `summary` - Brief one-line description
- `description` - Detailed explanation of what the endpoint does
- `operationId` - Unique identifier for the operation
- Response schemas with descriptions

### Recommended Error Responses

Define `400` and `500` error responses for all endpoints:

```yaml
responses:
  "200":
    description: Successful response
    # ...
  "400":
    description: Bad request - invalid parameters
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/Error"
  "500":
    description: Internal server error
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/Error"
```

## Step 3: Add CI Templates to Your Pipeline

Add the Partner API CI templates to your `.gitlab-ci.yml` to validate your spec automatically:

```yaml
include:
  # Spec validation (required)
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/spec-validation.yml'

  # Breaking change detection (required)
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/breaking-change-detection.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
  BASELINE_SPEC_URL: 'https://gitlab.example.com/api/v4/projects/platform%2Fpartner-api/repository/files/governance%2Fbaseline%2Fpartner-api.yaml/raw?ref=main'
```

See [CI Templates Usage Guide](./ci-templates-usage.md) for detailed configuration options.

## Step 4: Request Onboarding

Once your spec is ready and passing validation, request onboarding:

1. **Create an Issue**: Open a team onboarding request issue in the partner-api repository
2. **Provide Information**:
   - Team name and contact
   - Repository URL containing your OpenAPI spec
   - List of endpoints to expose (those with `partner` tag)
   - Expected traffic volume
3. **Platform Team Review**: The platform team will review your spec and coordinate onboarding

### What the Platform Team Will Do

1. Create a spec mapping file defining path transformations
2. Add your service to the manifest
3. Configure gateway routing
4. Run integration tests
5. Deploy to staging for validation

## Step 5: Path Transformations

Your internal paths will be transformed to partner paths:

| Internal Path | Partner Path |
|---------------|--------------|
| `/api/v3/orders` | `/partner/orders` |
| `/api/v3/orders/{id}` | `/partner/orders/{id}` |
| `/api/v1/users` | `/partner/users` |
| `/api/v1/users/{id}` | `/partner/users/{id}` |

The platform team defines these mappings in `governance/spec-mappings/<service>-mapping.yaml`.

## Validation Checklist

Before requesting onboarding, verify:

- [ ] All partner-exposed endpoints have the `partner` tag
- [ ] Property names use `snake_case`
- [ ] All operations have `summary`, `description`, and `operationId`
- [ ] CI pipeline passes spec validation
- [ ] CI pipeline passes breaking change detection (or no baseline exists yet)
- [ ] No sensitive internal data is exposed

## Common Issues and Solutions

### "Property must use snake_case"

Rename properties from camelCase to snake_case:

```yaml
# Before
properties:
  orderId: ...

# After
properties:
  order_id: ...
```

### "Operation must have description"

Add a description to your operation:

```yaml
get:
  summary: List orders
  description: Returns a paginated list of all orders for the authenticated partner.
  operationId: listOrders
```

### "Breaking changes detected"

If you need to make a breaking change, follow the [Breaking Change Workflow](./breaking-change-workflow.md).

### Endpoint not appearing in Partner API

Check that:
1. The endpoint has the `partner` tag
2. The endpoint is not excluded by `internal` or `admin` tags
3. Your service is enabled in the manifest

## Getting Help

- **Slack**: #partner-api-support
- **Office Hours**: Platform team holds weekly office hours (see calendar)
- **Documentation**: This guide and [CI Templates Usage](./ci-templates-usage.md)
- **Issues**: File in the partner-api repository

## Next Steps

After onboarding:

1. **Monitor**: Check Partner API gateway logs for your endpoints
2. **Iterate**: Add more endpoints by tagging them with `partner`
3. **Maintain**: Keep your spec up to date; breaking changes require coordination
