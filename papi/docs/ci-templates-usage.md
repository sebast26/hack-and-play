# CI Templates Usage Guide

This guide explains how internal teams can use the Partner API CI templates to validate their OpenAPI specs before contributing to the Partner API.

## Available Templates

| Template | Purpose |
|----------|---------|
| `spec-validation.yml` | Validates OpenAPI spec syntax and runs Spectral linting |

## Quick Start

Add the following to your `.gitlab-ci.yml`:

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/spec-validation.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
```

This will run spec validation on every merge request and branch push.

---

## Template: spec-validation.yml

Validates your OpenAPI spec against the Partner API design standards.

### What It Does

1. **Syntax Validation**: Checks that your OpenAPI spec is valid YAML/JSON
2. **Spectral Linting**: Enforces Partner API design rules:
   - All operations must have descriptions and operationIds
   - Property names must use `snake_case`
   - Error responses (400, 500) should be defined
   - Schema examples must be valid

### Required Variables

| Variable | Description |
|----------|-------------|
| `OPENAPI_SPEC_PATH` | Path to your OpenAPI spec file (e.g., `./openapi.yaml`) |

### Optional Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SPECTRAL_FAIL_SEVERITY` | `error` | Minimum severity to fail the job. Options: `error`, `warn`, `info`, `hint` |
| `SPECTRAL_RULESET_URL` | (embedded) | URL to a custom Spectral ruleset. If not set, uses the platform's standard ruleset |

### Example Usage

**Basic usage:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/spec-validation.yml'

variables:
  OPENAPI_SPEC_PATH: './api/openapi.yaml'
```

**Fail on warnings too:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/spec-validation.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
  SPECTRAL_FAIL_SEVERITY: 'warn'
```

**Use a custom ruleset:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/spec-validation.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
  SPECTRAL_RULESET_URL: 'https://example.com/my-custom-ruleset.yaml'
```

**Pin to a specific version:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: v1.0.0  # Use a tagged version for stability
    file: '/ci-templates/spec-validation.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
```

### Customizing the Stage

By default, the job runs in the `validate` stage. To use a different stage:

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/spec-validation.yml'

stages:
  - lint
  - build
  - test

spec-validation:
  stage: lint

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
```

### Common Issues

**Error: OPENAPI_SPEC_PATH variable is required**

You must set the `OPENAPI_SPEC_PATH` variable to point to your spec file.

**Error: OpenAPI spec not found**

Check that the path is correct relative to your repository root.

**Validation failures**

Review the Spectral output for specific rule violations. Common fixes:
- Add `description` to all operations
- Add `operationId` to all operations
- Rename properties to use `snake_case` (e.g., `userId` → `user_id`)
- Add `400` and `500` response definitions

---

## Running Locally

To validate your spec locally before pushing, install Spectral and download the ruleset:

```bash
# Install Spectral
npm install -g @stoplight/spectral-cli

# Download the platform ruleset
curl -o .spectral.yaml https://raw.githubusercontent.com/your-org/partner-api/main/governance/spectral/.spectral.yaml

# Run validation
spectral lint ./openapi.yaml --ruleset .spectral.yaml
```

---

## Getting Help

- **Slack**: #partner-api-support
- **Documentation**: See the [Onboarding Guide](./onboarding-guide.md)
- **Issues**: File in the partner-api repository
