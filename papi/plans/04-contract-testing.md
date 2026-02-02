# 04 - Contract Testing

## Purpose

Contract Testing ensures that internal services actually behave according to their OpenAPI specifications, and that changes to those specifications don't break the Partner API contract. It's the safety net that catches mismatches between documentation and reality before they affect partners.

This plan defines the testing strategy, what the CI templates provide, and how internal teams integrate them into their pipelines.

---

## Role in the Platform

Contract testing serves as the **verification layer** between specs and implementations:

| Risk | How Contract Testing Addresses It |
|------|-----------------------------------|
| Internal service doesn't match its spec | Schemathesis tests implementation against spec |
| Spec change would break Partner API | oasdiff compares spec against Partner API contract |
| Spec violates API design standards | Spectral linting enforces organizational rules |
| Changes deployed without verification | CI pipeline blocks deployment on test failures |

Contract testing does NOT:
- Test business logic correctness (that's unit/integration testing)
- Test performance or load handling (that's performance testing)
- Validate Partner API gateway behavior (that's gateway testing)
- Replace internal team's own test suites

---

## Testing Strategy: Provider-Driven

The platform uses **provider-driven contract testing**, where internal teams (providers) validate that their implementations match their specifications.

**Why provider-driven (not consumer-driven)**:
- Consumer-driven (like Pact) requires the Partner API team to define expectations for 8-15 internal services—impractical coordination burden
- Internal teams already own their OpenAPI specs
- OpenAPI spec serves as the contract; no separate contract definition needed
- Testing runs in internal team CI, not centrally

**The contract is the OpenAPI specification**. If the spec says an endpoint returns a `customer_id` field, the implementation must return that field.

---

## Three Layers of Contract Testing

### Layer 1: Specification Linting (Spectral)

**What it validates**: The OpenAPI spec itself is well-formed and follows organizational standards.

**When it runs**: On every merge request that modifies the OpenAPI spec.

**What it catches**:
- Missing descriptions on endpoints or parameters
- Inconsistent naming conventions (camelCase vs. snake_case)
- Missing response schemas
- Security scheme not defined
- Non-standard HTTP status codes
- Missing error response definitions

**Failure behavior**: Merge request blocked until linting errors resolved.

**Example rules**:
| Rule | Purpose |
|------|---------|
| `operation-description` | All operations must have descriptions |
| `oas3-valid-schema-example` | Examples must match their schemas |
| `path-params-defined` | Path parameters must be defined |
| `response-schema-defined` | Responses must have schemas |
| Custom: `partner-api-naming` | Field names must be snake_case |
| Custom: `partner-api-errors` | Must define 400, 401, 500 responses |

---

### Layer 2: Breaking Change Detection (oasdiff)

**What it validates**: Changes to the spec don't break the Partner API contract.

**When it runs**: On every merge request that modifies the OpenAPI spec.

**What it catches**:
Breaking changes as defined in PLAN.md, including: removed endpoints, removed or renamed response fields, changed field types, added required request fields, narrowed enum values, and changed authentication requirements.

**Failure behavior**: Merge request blocked; team directed to Breaking Change Lifecycle workflow.

**Comparison target**: The internal team's spec is compared against the current Partner API spec (fetched from governance repo or published URL).

**Important distinction**:
- If the change affects endpoints NOT exposed to Partner API → allowed
- If the change affects endpoints exposed to Partner API → breaking change detection applies

The CI template filters comparison to only partner-tagged endpoints.

---

### Layer 3: Implementation Validation (Schemathesis)

**What it validates**: The running service actually behaves according to its OpenAPI spec.

**When it runs**: On every merge request, against the team's staging environment.

**What it catches**:
- Endpoints returning fields not in spec
- Endpoints missing fields that spec says are required
- Response status codes not matching spec
- Response formats not matching schema (wrong types, missing properties)
- Request validation not matching spec (accepting invalid input, rejecting valid input)

**Failure behavior**: Merge request blocked until implementation matches spec.

**How Schemathesis works**:
1. Reads the OpenAPI spec
2. Generates test cases based on spec (valid inputs, edge cases, invalid inputs)
3. Calls the actual service endpoints
4. Validates responses match spec schemas
5. Reports mismatches as failures

**Test environment**: Internal team's staging environment. Each team maintains their own staging deployment where their service runs with realistic (but non-production) data.

---

## Two-Phase Testing Approach

To handle cases where the spec in an MR is newer than what's deployed to staging, the CI pipeline uses a two-phase approach:

**Phase 1: Spec Validation (against MR spec)**
- Spec syntax validation
- Spectral linting
- Breaking change detection against Partner API

This phase validates the spec in the MR is correct, well-formed, and won't break partners.

**Phase 2: Implementation Testing (against deployed spec)**
- Schemathesis fetches the OpenAPI spec from the running staging service (e.g., `$STAGING_SERVICE_URL/openapi.yaml`)
- Tests validate that staging implementation matches what staging claims to be

This approach ensures:
- Spec quality is checked immediately on every MR
- Implementation testing doesn't produce false failures when spec changes haven't been deployed yet
- After deployment, implementation tests validate the deployed code matches the deployed spec

**Staging service requirement**: Services must expose their OpenAPI spec at a known endpoint (e.g., `/openapi.yaml` or `/api/docs/openapi.yaml`). The CI template fetches this for implementation testing.

---

## CI Templates

The platform team provides CI templates that internal teams include in their pipelines. Templates encapsulate all tooling so teams don't need to understand the tools.

### Template: `spec-validation.yml`

**Purpose**: Validate OpenAPI spec syntax and organizational standards.

**What it does**:
1. Parses OpenAPI spec to verify valid syntax
2. Runs Spectral with platform-defined ruleset
3. Reports errors with file/line references

**Required variables**:
| Variable | Description | Example |
|----------|-------------|---------|
| `OPENAPI_SPEC_PATH` | Path to spec file in repo | `api/openapi.yaml` |

**Outputs**: Pass/fail with detailed error messages.

---

### Template: `breaking-change-detection.yml`

**Purpose**: Detect changes that would break the Partner API.

**What it does**:
1. Fetches current Partner API spec (from published URL)
2. Extracts partner-tagged endpoints from internal spec
3. Runs oasdiff comparing previous vs. current
4. Fails if breaking changes detected

**Required variables**:
| Variable | Description | Example |
|----------|-------------|---------|
| `OPENAPI_SPEC_PATH` | Path to spec file in repo | `api/openapi.yaml` |
| `PARTNER_TAG` | Tag identifying partner endpoints | `partner` |

**Optional variables**:
| Variable | Description | Default |
|----------|-------------|---------|
| `PARTNER_API_SPEC_URL` | URL to fetch Partner API spec | Platform default URL |

**Outputs**: Pass/fail with list of breaking changes if any.

**On failure**: Message includes:
- What breaking changes were detected
- Link to Breaking Change Lifecycle documentation
- Instructions to coordinate with platform team

---

### Template: `implementation-testing.yml`

**Purpose**: Validate service implementation matches spec.

**What it does**:
1. Fetches OpenAPI spec from staging service (deployed spec)
2. Runs Schemathesis against staging service
3. Tests endpoints based on configuration
4. Validates responses match schemas
5. Reports mismatches

**Required variables**:
| Variable | Description | Example |
|----------|-------------|---------|
| `STAGING_SERVICE_URL` | Base URL of staging service | `https://orders-staging.internal.company.com` |
| `STAGING_SPEC_PATH` | Path to fetch spec from staging | `/openapi.yaml` |

**Optional variables**:
| Variable | Description | Default |
|----------|-------------|---------|
| `STAGING_AUTH_TOKEN` | Bearer token for authenticated endpoints | None |
| `SCHEMATHESIS_WORKERS` | Parallel test workers | `4` |
| `SCHEMATHESIS_MAX_EXAMPLES` | Max test cases per endpoint | `100` |
| `TEST_TIMEOUT_MINUTES` | Global timeout for test job | `60` |
| `TEST_RETRIES` | Number of retries on failure | `2` |

**Authentication**: Teams store authentication tokens as **masked GitLab CI variables** in their project settings. The CI template reads from `$STAGING_AUTH_TOKEN` if provided. Tokens should:
- Be scoped to staging environment only
- Have sufficient permissions to call all tested endpoints
- Be rotated according to team's security practices

**Outputs**: Pass/fail with detailed mismatch report.

---

### Endpoint Test Configuration

Teams can configure how different endpoints are tested by providing a configuration file:

```yaml
# .partner-api-test-config.yaml
endpoints:
  # Test read endpoints fully
  "GET /orders":
    mode: full
    
  "GET /orders/{id}":
    mode: full
    
  # Skip dangerous write endpoints
  "DELETE /orders/{id}":
    mode: skip
    reason: "Deletes production-like data in staging"
    
  # Limit iterations for write endpoints
  "POST /orders":
    mode: limited
    max_examples: 5
    
  "PUT /orders/{id}":
    mode: limited
    max_examples: 5
```

**Test modes**:
| Mode | Description |
|------|-------------|
| `full` | Run all generated test cases (default for GET endpoints) |
| `limited` | Run reduced number of test cases (for write endpoints) |
| `skip` | Skip testing this endpoint (must provide reason) |

If no configuration file exists, defaults apply:
- GET/HEAD endpoints: `full` mode
- POST/PUT/PATCH/DELETE endpoints: `limited` mode with `max_examples: 10`

---

### Test Execution Settings

**Timeout**: All implementation tests have a **global timeout of 60 minutes**. If tests exceed this duration, the job fails. Teams with large APIs should optimize by:
- Reducing `max_examples` for some endpoints
- Skipping non-critical endpoints
- Splitting tests into multiple jobs

**Retries**: Tests are retried up to **2 times** before being considered failed. This handles transient issues like:
- Network hiccups
- Staging service restarts
- Timing-related failures

There is no mechanism to mark endpoints as "known flaky". If an endpoint consistently fails, the team must either fix the underlying issue or adjust their test configuration (e.g., skip with documented reason).

**Parallelism**: Implementation tests run **in parallel with other CI jobs** (unit tests, linting, etc.) to minimize total pipeline duration. The test job is independent and doesn't block or wait for other jobs.

---

### Combined Template: `partner-api-checks.yml`

**Purpose**: Single include that runs all three layers.

**What it does**: Orchestrates spec-validation, breaking-change-detection, and implementation-testing in appropriate order.

**Execution order**:
1. Spec validation (fast, catches syntax/lint errors first) - uses MR spec
2. Breaking change detection (medium, requires spec parsing) - uses MR spec
3. Implementation testing (slow, requires network calls to staging) - uses deployed spec

Pipeline short-circuits: if earlier stage fails, later stages don't run.

**Example usage in internal team's `.gitlab-ci.yml`**:

```yaml
include:
  - project: 'platform/partner-api-governance'
    file: '/ci-templates/partner-api-checks.yml'
    ref: v1  # Pin to major version

variables:
  OPENAPI_SPEC_PATH: "api/openapi.yaml"
  STAGING_SERVICE_URL: "https://orders-staging.internal.company.com"
  STAGING_SPEC_PATH: "/openapi.yaml"
  STAGING_AUTH_TOKEN: $ORDERS_STAGING_TOKEN  # From CI variables
  PARTNER_TAG: "partner"
```

---

## What Internal Teams See

### On Success

All checks pass; merge request can proceed:

```
✓ Spec Validation: Passed
  - OpenAPI syntax valid
  - 0 linting errors
  
✓ Breaking Change Detection: Passed
  - No breaking changes to Partner API
  
✓ Implementation Testing: Passed
  - 47 endpoints tested
  - 2,350 test cases executed
  - All responses match spec
  - Duration: 8m 32s
```

### On Spec Validation Failure

```
✗ Spec Validation: Failed

Linting errors found:
  api/openapi.yaml:45:5 - operation-description: Operation must have a description
  api/openapi.yaml:78:9 - response-schema-defined: Response 200 must have a schema

Fix these errors and push again.
```

### On Breaking Change Detection Failure

```
✗ Breaking Change Detection: Failed

Breaking changes detected in partner-exposed endpoints:

  DELETE /orders/{id}
    - Removed required response field: customer_name
    - This field is used by Partner API version 2024-01-15

  POST /orders
    - Added required request field: priority
    - Existing partners don't send this field

These changes would break the Partner API contract.

Next steps:
  1. Post in #partner-api-changes Slack channel
  2. Create JIRA ticket with change details
  3. Platform team will coordinate transformation
  4. See: https://docs.internal/partner-api/breaking-changes

DO NOT merge until coordinated with platform team.
```

### On Implementation Testing Failure

```
✗ Implementation Testing: Failed (after 2 retries)

Service responses don't match OpenAPI spec:

  GET /orders/{id}
    Response 200:
      - Missing required field: shipping_address
      - Field 'status' has type 'number', expected 'string'
    
  POST /orders
    Response 201:
      - Extra field not in spec: internal_debug_info
      - Field 'created_at' format is 'date', expected 'date-time'

Your service implementation doesn't match your OpenAPI spec.
Either update your code to match the spec, or update the spec to match your code.

Note: Implementation tests run against your deployed staging service.
Make sure your latest code is deployed to staging before running these tests.
```

---

## Staging Environment Requirements

For implementation testing to work, internal teams must maintain a staging environment:

| Requirement | Rationale |
|-------------|-----------|
| Service deployed and accessible | Schemathesis needs to make HTTP calls |
| OpenAPI spec exposed at known path | Implementation tests fetch deployed spec |
| Network accessible from CI runners | CI pipeline must reach staging URL |
| Representative data available | Tests need realistic responses |
| Authentication configured | If endpoints require auth, CI needs valid token |
| Isolated from production | Tests may create/modify data |

**Who maintains staging**: Each internal team is responsible for their own staging environment. Platform team doesn't provision or maintain these.

**Test data management**: Teams are responsible for ensuring staging has representative data. The CI template supports configuring endpoint test modes (full, limited, skip) so teams can control how write endpoints are tested and avoid unwanted data modifications.

---

## Failure Handling

| Failure Type | Blocks Deployment | Retries | Team Action |
|--------------|-------------------|---------|-------------|
| Spec syntax invalid | Yes | No | Fix YAML/JSON errors |
| Linting errors | Yes | No | Fix spec to meet standards |
| Breaking change detected | Yes | No | Coordinate with platform team |
| Implementation mismatch | Yes | 2 | Fix code or update spec |
| Staging unreachable | Yes | 2 | Fix staging environment |
| Timeout exceeded | Yes | No | Optimize tests or increase timeout |

**All failures block deployment**. This ensures:
- Specs are always valid and compliant
- Breaking changes are always coordinated
- Implementations always match specs
- Partners never see unexpected behavior

---

## Integration Points

### With Governance Repository

- CI templates stored in governance repo (version-tagged)
- Spectral ruleset defined in governance repo
- Partner API spec URL configured in governance repo

### With Internal Team Repositories

- Teams include CI templates via GitLab `include:` with version pinning
- Templates run as part of team's existing pipeline (parallel execution)
- Variables configured in team's `.gitlab-ci.yml` and CI settings
- Auth tokens stored as masked CI variables

### With Breaking Change Lifecycle

- Breaking change detection failure triggers the coordination workflow
- Failure message links to documentation and Slack channel

### With Spec Aggregation Pipeline

- Implementation testing validates that what the aggregation pipeline sees (specs) matches reality (running services)
- Provides confidence that transformed Partner API spec reflects actual capabilities

---

## Success Criteria

Contract testing is complete when:

| Criteria | Validation |
|----------|------------|
| CI templates work across teams | At least 3 teams successfully using templates |
| Spec violations caught | Intentional linting error fails pipeline |
| Breaking changes caught | Intentional breaking change fails pipeline |
| Implementation mismatches caught | Intentional spec/code mismatch fails pipeline |
| Two-phase testing works | Spec changes don't cause false failures in implementation tests |
| Endpoint configuration works | Teams can configure test modes per endpoint |
| Retries handle transient failures | Transient network issues don't block MRs |
| Timeout enforced | Tests exceeding 60 minutes are terminated |
| Failure messages are actionable | Teams can resolve failures without platform team help (except breaking changes) |
| No false positives | Tests don't fail for valid code/spec combinations |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 | 2026-01-28 | Initial draft |
| 0.2 | 2026-01-30 | Resolved open questions: GitLab CI variables for auth; existing staging data with endpoint configuration; 2 retries with no flaky marking; 60 minute timeout with parallel execution; two-phase testing (MR spec for validation, deployed spec for implementation) |
