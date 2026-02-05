## Breaking Change Request

<!--
Use this template to request a breaking change to the Partner API.
Please read docs/breaking-change-workflow.md before submitting.
-->

### Requester Information

- **Team**:
- **Contact**:
- **Slack Handle**:

### Change Description

**What are you changing?**
<!-- Describe the change in detail -->

**Which endpoints/fields are affected?**
| Endpoint | Field/Parameter | Current | Proposed |
|----------|-----------------|---------|----------|
| `/partner/...` | ... | ... | ... |

### Justification

**Why is this change necessary?**
<!-- Explain why the change cannot be avoided -->

**Why can't this be done in a non-breaking way?**
<!-- Explain alternatives considered -->

### Impact Assessment

**How many partners use the affected endpoints?**
- [ ] Unknown
- [ ] None (new/unused endpoints)
- [ ] Few (1-5 partners)
- [ ] Some (5-20 partners)
- [ ] Many (20+ partners)

**What will break for partners?**
<!-- Describe what partners will experience -->

**Is there a migration path?**
<!-- Describe how partners can adapt -->

### Proposed Timeline

- **Requested deploy date**:
- **Proposed deprecation period**: 30 / 60 / 90 days
- **Reason for timeline**:

### Communication Plan

- [ ] I will draft partner communication
- [ ] I will create a migration guide
- [ ] I will be available to answer partner questions

### Checklist

- [ ] I have read [docs/breaking-change-workflow.md](../docs/breaking-change-workflow.md)
- [ ] I have posted in #partner-api-changes Slack channel
- [ ] I have considered non-breaking alternatives
- [ ] I understand this requires platform team approval

### Additional Context

<!-- Any other relevant information -->

---

/label ~"breaking-change" ~"needs-review"
/assign @platform-team
