# CI Templates Usage Guide

This guide explains how internal teams can use the Partner API CI templates to validate their OpenAPI specs before contributing to the Partner API.

## Available Templates

| Template | Purpose |
|----------|---------|
| `spec-validation.yml` | Validates OpenAPI spec syntax and runs Spectral linting |
| `breaking-change-detection.yml` | Detects breaking API changes using oasdiff |

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

## Template: breaking-change-detection.yml

Detects breaking API changes by comparing your spec against a baseline using [oasdiff](https://github.com/oasdiff/oasdiff).

### What It Does

1. **Fetches Baseline**: Downloads the baseline spec from a URL you provide
2. **Compares Specs**: Uses oasdiff to detect breaking changes
3. **Reports Results**: Lists any breaking changes found
4. **Changelog** (optional job): Shows a summary of all API changes

### Jobs Included

| Job | Purpose | Fails on breaking changes? |
|-----|---------|---------------------------|
| `breaking-change-detection` | Detects breaking changes | Yes |
| `breaking-change-diff` | Shows changelog of all changes | No (informational) |

### Required Variables

| Variable | Description |
|----------|-------------|
| `OPENAPI_SPEC_PATH` | Path to your current OpenAPI spec file |
| `BASELINE_SPEC_URL` | URL to fetch the baseline spec for comparison |

### Optional Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OASDIFF_VERSION` | `latest` | Version of oasdiff to use |
| `FAIL_ON_SEVERITY` | `WARN` | Minimum severity to fail. Options: `ERR`, `WARN`, `INFO` |

### Example Usage

**Basic usage:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/breaking-change-detection.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
  BASELINE_SPEC_URL: 'https://gitlab.example.com/api/v4/projects/123/repository/files/governance%2Fbaseline%2Fpartner-api.yaml/raw?ref=main'
```

**Using GitLab CI job token for private repos:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/breaking-change-detection.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
  # Use CI_JOB_TOKEN for authentication to private GitLab repos
  BASELINE_SPEC_URL: 'https://gitlab.example.com/api/v4/projects/platform%2Fpartner-api/repository/files/governance%2Fbaseline%2Fpartner-api.yaml/raw?ref=main&private_token=${CI_JOB_TOKEN}'
```

**Only fail on errors (ignore warnings):**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/breaking-change-detection.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
  BASELINE_SPEC_URL: 'https://example.com/baseline-spec.yaml'
  FAIL_ON_SEVERITY: 'ERR'
```

**Disable the changelog job:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/breaking-change-detection.yml'

variables:
  OPENAPI_SPEC_PATH: './openapi.yaml'
  BASELINE_SPEC_URL: 'https://example.com/baseline-spec.yaml'

breaking-change-diff:
  rules:
    - when: never
```

### What Counts as a Breaking Change?

oasdiff detects many types of breaking changes, including:

- Removing an endpoint
- Removing a request parameter
- Making an optional parameter required
- Changing a parameter type
- Removing a response property
- Changing response status codes

For the full list, see the [oasdiff documentation](https://github.com/oasdiff/oasdiff#breaking-changes).

### Common Issues

**Error: Failed to fetch baseline spec**

Check that:
- The `BASELINE_SPEC_URL` is correct
- You have access to the URL (authentication may be required)
- The file exists at the specified ref/branch

**Breaking changes detected unexpectedly**

Review the output carefully. Common causes:
- Renamed a field (oasdiff sees this as remove + add)
- Changed a type (e.g., `string` to `integer`)
- Made a previously optional field required

If the change is intentional, follow the breaking change workflow documented in [breaking-change-workflow.md](./breaking-change-workflow.md).

---

## Running Locally

### Spec Validation

To validate your spec locally before pushing:

```bash
# Install Spectral
npm install -g @stoplight/spectral-cli

# Download the platform ruleset
curl -o .spectral.yaml https://raw.githubusercontent.com/your-org/partner-api/main/governance/spectral/.spectral.yaml

# Run validation
spectral lint ./openapi.yaml --ruleset .spectral.yaml
```

### Breaking Change Detection

To check for breaking changes locally:

```bash
# Install oasdiff (macOS)
brew install oasdiff

# Or download directly
curl -sL https://github.com/oasdiff/oasdiff/releases/latest/download/oasdiff_<version>_<os>_<arch>.tar.gz | tar -xz

# Download the baseline spec
curl -o baseline.yaml https://your-gitlab/path/to/baseline/partner-api.yaml

# Check for breaking changes
oasdiff breaking baseline.yaml ./openapi.yaml --format text --fail-on WARN

# View all changes (changelog)
oasdiff changelog baseline.yaml ./openapi.yaml --format text
```

---

## Getting Help

- **Slack**: #partner-api-support
- **Documentation**: See the [Onboarding Guide](./onboarding-guide.md)
- **Issues**: File in the partner-api repository
