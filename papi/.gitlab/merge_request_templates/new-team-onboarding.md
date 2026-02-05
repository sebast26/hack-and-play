## New Team Onboarding

<!-- Use this template when adding a new internal team to the Partner API -->

### Team Information

- **Team Name**:
- **Team Contact**:
- **Slack Channel**:
- **Service Name**:
- **Repository URL**:

### Endpoints to Expose

<!-- List the endpoints that will be exposed through the Partner API -->

| Internal Path | Partner Path | Methods | Description |
|---------------|--------------|---------|-------------|
| `/api/v1/...` | `/partner/...` | GET | ... |

### Checklist

**Team Preparation:**
- [ ] OpenAPI spec follows Partner API design standards
- [ ] All partner endpoints have `partner` tag
- [ ] Property names use `snake_case`
- [ ] All operations have `summary`, `description`, `operationId`
- [ ] CI templates added to team's pipeline
- [ ] Spec validation passing in team's CI

**Platform Team Actions:**
- [ ] Reviewed OpenAPI spec
- [ ] Created spec mapping file (`governance/spec-mappings/<service>-mapping.yaml`)
- [ ] Added service to manifest (`governance/manifest.yaml`)
- [ ] Updated `governance/openapi-merge.yaml`
- [ ] Ran aggregation pipeline successfully
- [ ] Updated gateway configuration
- [ ] Tested endpoints in staging
- [ ] Updated CODEOWNERS if needed

### Files Changed

<!-- Auto-populated by GitLab, or list manually -->

- `governance/manifest.yaml`
- `governance/spec-mappings/<service>-mapping.yaml`
- `governance/openapi-merge.yaml`

### Testing

**Aggregation Pipeline:**
```bash
./scripts/aggregate.sh
```

**Breaking Change Check:**
```bash
./scripts/breaking-change-check.sh
```

**Gateway Test:**
```bash
docker-compose up -d
curl http://localhost:8080/partner/<endpoint>
```

### Notes

<!-- Any additional context or notes -->

---

/label ~"team-onboarding" ~"partner-api"
/assign @platform-team
