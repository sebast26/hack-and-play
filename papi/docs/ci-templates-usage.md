# CI Templates Usage Guide

This guide explains how internal teams can use the Partner API CI templates to validate their OpenAPI specs before contributing to the Partner API.

## Available Templates

| Template | Purpose |
|----------|---------|
| `partner-api-checks.yml` | **Recommended** - Combined template that runs all three checks |
| `spec-validation.yml` | Validates OpenAPI spec syntax and runs Spectral linting |
| `breaking-change-detection.yml` | Detects breaking API changes using oasdiff |
| `implementation-testing.yml` | Tests API implementation against spec using Schemathesis |

## Quick Start (Recommended)

For most teams, we recommend using the combined template that runs all checks:

```yaml
include:
  - project: 'platform/partner-api'
    ref: v1  # Pin to major version for stability
    file: '/ci-templates/partner-api-checks.yml'

variables:
  OPENAPI_SPEC_PATH: './api/openapi.yaml'
  BASELINE_SPEC_URL: 'https://gitlab.example.com/api/v4/projects/123/repository/files/governance%2Fbaseline%2Fpartner-api.yaml/raw?ref=main'
  STAGING_SERVICE_URL: 'https://staging.myservice.example.com'
  STAGING_SPEC_URL: 'https://staging.myservice.example.com/openapi.yaml'
  STAGING_AUTH_TOKEN: $MY_SERVICE_STAGING_TOKEN  # From CI variables
```

This runs all three validation layers in order:
1. **Spec Validation** - Fast syntax and linting checks
2. **Breaking Change Detection** - Compares against baseline
3. **Implementation Testing** - Tests your running service

If earlier stages fail, later stages are skipped for faster feedback.

## Individual Templates

You can also use templates individually if you only need specific checks.

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

## Template: implementation-testing.yml

Tests your API implementation against its OpenAPI spec using [Schemathesis](https://schemathesis.io/).

### What It Does

1. **Fetches your OpenAPI spec** from your staging service
2. **Generates test cases** based on the spec (property-based testing)
3. **Sends requests** to your staging service
4. **Validates responses** match the spec
5. **Retries on transient failures** (up to 2 retries by default)
6. **Reports failures** with reproduction commands

### Jobs Included

| Job | Purpose | Fails on errors? |
|-----|---------|------------------|
| `implementation-testing` | Full test suite with retries | Yes (configurable) |
| `implementation-smoke-test` | Quick validation | No |

### Required Variables

| Variable | Description |
|----------|-------------|
| `STAGING_SERVICE_URL` | Base URL of your staging service |
| `STAGING_SPEC_URL` | URL to fetch OpenAPI spec from staging |

### Optional Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `STAGING_AUTH_TOKEN` | (none) | Bearer token for authenticated endpoints (store as masked CI variable) |
| `SCHEMATHESIS_CHECKS` | Core checks | Comma-separated checks to run |
| `SCHEMATHESIS_WORKERS` | `4` | Number of parallel test workers |
| `SCHEMATHESIS_MAX_EXAMPLES` | `100` | Max test cases per endpoint |
| `SCHEMATHESIS_MAX_FAILURES` | `20` | Stop after N failures |
| `TEST_TIMEOUT_MINUTES` | `60` | Global timeout for test job |
| `TEST_RETRIES` | `2` | Number of retries on failure |

### Example Usage

**Basic usage:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/implementation-testing.yml'

variables:
  STAGING_SERVICE_URL: 'https://staging.myservice.example.com'
  STAGING_SPEC_URL: 'https://staging.myservice.example.com/openapi.yaml'
```

**With authentication:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: main
    file: '/ci-templates/implementation-testing.yml'

variables:
  STAGING_SERVICE_URL: 'https://staging.myservice.example.com'
  STAGING_SPEC_URL: 'https://staging.myservice.example.com/openapi.yaml'

implementation-testing:
  script:
    - |
      schemathesis run "$STAGING_SPEC_URL" \
        --url "$STAGING_SERVICE_URL" \
        --header "Authorization: Bearer $STAGING_AUTH_TOKEN" \
        --checks "$SCHEMATHESIS_CHECKS" \
        --report junit --report-dir ./schemathesis-reports
```

### Available Checks

| Check | Description |
|-------|-------------|
| `not_a_server_error` | Response is not 5xx |
| `status_code_conformance` | Status code matches spec |
| `content_type_conformance` | Content-Type matches spec |
| `response_schema_conformance` | Response body matches schema |
| `response_headers_conformance` | Headers match spec |
| `negative_data_rejection` | Invalid data returns 4xx |
| `all` | All checks |

### Common Issues

**5xx errors for valid requests**

Your implementation may have bugs. Check the reproduction command in the output.

**Response doesn't match schema**

Common causes:
- Missing required fields
- Wrong field types
- Extra fields not in spec (if `additionalProperties: false`)

**Timeouts**

Increase `TEST_TIMEOUT_MINUTES` or optimize your service. The default is 60 minutes.

---

## Template: partner-api-checks.yml (Combined)

The combined template runs all three validation layers in sequence. This is the **recommended** approach for most teams.

### What It Does

1. **Spec Validation** (fastest) - Validates syntax and runs Spectral linting
2. **Breaking Change Detection** (medium) - Compares against baseline spec
3. **Implementation Testing** (slowest) - Tests your running service

The pipeline **short-circuits**: if spec validation fails, later stages are skipped.

### Required Variables

| Variable | Description |
|----------|-------------|
| `OPENAPI_SPEC_PATH` | Path to your OpenAPI spec file |
| `BASELINE_SPEC_URL` | URL to fetch baseline spec for breaking change detection |
| `STAGING_SERVICE_URL` | Base URL of your staging service |
| `STAGING_SPEC_URL` | URL to fetch OpenAPI spec from staging |

### Optional Variables

All optional variables from the individual templates are supported, plus:

| Variable | Default | Description |
|----------|---------|-------------|
| `SKIP_SPEC_VALIDATION` | (empty) | Set to `true` to skip spec validation |
| `SKIP_BREAKING_CHANGE` | (empty) | Set to `true` to skip breaking change detection |
| `SKIP_IMPLEMENTATION_TESTING` | (empty) | Set to `true` to skip implementation testing |

### Example Usage

**Full configuration:**

```yaml
include:
  - project: 'platform/partner-api'
    ref: v1
    file: '/ci-templates/partner-api-checks.yml'

variables:
  # Required
  OPENAPI_SPEC_PATH: './api/openapi.yaml'
  BASELINE_SPEC_URL: 'https://gitlab.example.com/api/v4/projects/123/repository/files/governance%2Fbaseline%2Fpartner-api.yaml/raw?ref=main'
  STAGING_SERVICE_URL: 'https://staging.orders.example.com'
  STAGING_SPEC_URL: 'https://staging.orders.example.com/openapi.yaml'

  # Optional - Authentication
  STAGING_AUTH_TOKEN: $ORDERS_STAGING_TOKEN  # Masked CI variable

  # Optional - Customize behavior
  SPECTRAL_FAIL_SEVERITY: 'error'
  FAIL_ON_SEVERITY: 'WARN'
  TEST_RETRIES: '2'
```

**Skip implementation testing (for spec-only changes):**

```yaml
include:
  - project: 'platform/partner-api'
    ref: v1
    file: '/ci-templates/partner-api-checks.yml'

variables:
  OPENAPI_SPEC_PATH: './api/openapi.yaml'
  BASELINE_SPEC_URL: 'https://example.com/baseline.yaml'
  STAGING_SERVICE_URL: 'https://staging.example.com'
  STAGING_SPEC_URL: 'https://staging.example.com/openapi.yaml'
  SKIP_IMPLEMENTATION_TESTING: 'true'
```

**Without breaking change detection (no baseline yet):**

```yaml
include:
  - project: 'platform/partner-api'
    ref: v1
    file: '/ci-templates/partner-api-checks.yml'

variables:
  OPENAPI_SPEC_PATH: './api/openapi.yaml'
  # BASELINE_SPEC_URL not set - breaking change detection skipped automatically
  STAGING_SERVICE_URL: 'https://staging.example.com'
  STAGING_SPEC_URL: 'https://staging.example.com/openapi.yaml'
```

### What You'll See

**On success:**

```
✓ Spec Validation: Passed
  - OpenAPI syntax valid
  - 0 linting errors

✓ Breaking Change Detection: Passed
  - No breaking changes to Partner API

✓ Implementation Testing: Passed
  - All responses match OpenAPI spec
```

**On spec validation failure:**

```
✗ Spec Validation: Failed

Linting errors found:
  api/openapi.yaml:45:5 - operation-description: Operation must have a description

Fix these errors and push again.
```

**On breaking change failure:**

```
✗ Breaking Change Detection: FAILED

Breaking changes detected in partner-exposed endpoints.

Next steps:
  1. Post in #partner-api-changes Slack channel
  2. Create a JIRA ticket with change details
  3. See: docs/breaking-change-workflow.md

DO NOT merge until coordinated with platform team.
```

**On implementation testing failure (after retries):**

```
✗ Implementation Testing: FAILED

Your service responses don't match your OpenAPI spec.
Review the failures above and either:
  - Update your code to match the spec, or
  - Update the spec to match your code

Note: Tests run against your deployed staging service.
```

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

### Implementation Testing

To run Schemathesis tests locally:

```bash
# Using Docker (recommended)
docker run --rm --network host \
  -v ./openapi.yaml:/spec.yaml:ro \
  schemathesis/schemathesis:stable run /spec.yaml \
  --url http://localhost:8080 \
  --checks all

# Using the platform script
./scripts/schemathesis-test.sh \
  --spec ./openapi.yaml \
  --url http://localhost:8080

# With authentication
docker run --rm --network host \
  -v ./openapi.yaml:/spec.yaml:ro \
  schemathesis/schemathesis:stable run /spec.yaml \
  --url http://localhost:8080 \
  --header "Authorization: Bearer YOUR_TOKEN" \
  --checks all
```

### All Three Checks Locally

Run all checks in sequence locally:

```bash
# 1. Spec Validation
npm install -g @stoplight/spectral-cli
spectral lint ./openapi.yaml

# 2. Breaking Change Detection
brew install oasdiff  # or download binary
oasdiff breaking baseline.yaml ./openapi.yaml --format text --fail-on WARN

# 3. Implementation Testing
./scripts/schemathesis-test.sh --spec ./openapi.yaml --url http://localhost:8080
```

---

## Authentication Setup

For implementation testing against authenticated endpoints:

1. **Create a service account** or test user in your staging environment
2. **Generate a token** with access to all Partner API endpoints
3. **Store as masked CI variable** in your GitLab project:
   - Go to Settings → CI/CD → Variables
   - Add variable: `MY_SERVICE_STAGING_TOKEN`
   - Check "Mask variable" to hide in logs
   - Check "Protect variable" if only needed on protected branches
4. **Reference in your `.gitlab-ci.yml`**:
   ```yaml
   variables:
     STAGING_AUTH_TOKEN: $MY_SERVICE_STAGING_TOKEN
   ```

---

## Getting Help

- **Slack**: #partner-api-support
- **Documentation**: See the [Onboarding Guide](./onboarding-guide.md)
- **Issues**: File in the partner-api repository
