## Team Onboarding Request

<!--
Use this template to request onboarding your team to the Partner API.
Please read docs/onboarding-guide.md before submitting.
-->

### Team Information

- **Team Name**:
- **Team Lead**:
- **Technical Contact**:
- **Slack Channel**:
- **Email Distribution List**:

### Service Information

- **Service Name**:
- **Service Description**:
- **Repository URL**:
- **OpenAPI Spec Path**: (e.g., `/openapi.yaml`)

### Endpoints to Expose

<!-- List the endpoints you want to expose through the Partner API -->

| Internal Path | HTTP Method | Description | Use Case |
|---------------|-------------|-------------|----------|
| `/api/v1/...` | GET | ... | ... |
| `/api/v1/...` | POST | ... | ... |

### Expected Usage

- **Expected traffic volume**: requests/day
- **Expected number of partners using this**:
- **Launch timeline**:

### Pre-Onboarding Checklist

**Spec Preparation:**
- [ ] OpenAPI 3.0+ spec exists
- [ ] Spec is valid YAML/JSON
- [ ] All target endpoints have `partner` tag
- [ ] Property names use `snake_case`
- [ ] All operations have `summary`, `description`, `operationId`

**CI Integration:**
- [ ] Added `spec-validation.yml` CI template to pipeline
- [ ] Spec validation is passing
- [ ] Added `breaking-change-detection.yml` CI template (optional for new teams)

**Documentation:**
- [ ] I have read [docs/onboarding-guide.md](../docs/onboarding-guide.md)
- [ ] I have read [docs/ci-templates-usage.md](../docs/ci-templates-usage.md)

### Questions

<!-- Any questions for the platform team -->

### Additional Context

<!-- Any other relevant information about your service or requirements -->

---

/label ~"team-onboarding" ~"needs-triage"
/assign @platform-team
