## Spec Mapping Change

<!-- Use this template when modifying transformation rules for an existing service -->

### Change Summary

- **Service**:
- **Requested By**:
- **Reason for Change**:

### Type of Change

<!-- Check all that apply -->

- [ ] Path transformation change
- [ ] Add new endpoints
- [ ] Remove endpoints
- [ ] Modify filtering rules
- [ ] Schema modifications
- [ ] Backend host change

### Changes Description

<!-- Describe what is changing and why -->

**Before:**
```yaml
# Current mapping configuration
```

**After:**
```yaml
# New mapping configuration
```

### Impact Assessment

**Affected Endpoints:**
| Endpoint | Change | Impact |
|----------|--------|--------|
| `/partner/...` | ... | ... |

**Breaking Changes:**
- [ ] This change introduces breaking changes
- [ ] Breaking change workflow completed (if applicable)
- [ ] Partners have been notified (if breaking)

### Checklist

**Requester:**
- [ ] Change has been discussed with platform team
- [ ] Impact on existing partners assessed
- [ ] JIRA ticket created (if breaking change)

**Platform Team:**
- [ ] Reviewed mapping changes
- [ ] Verified aggregation pipeline succeeds
- [ ] Verified no unintended breaking changes
- [ ] Tested affected endpoints
- [ ] Updated documentation if needed

### Testing

**Aggregation Pipeline:**
```bash
./scripts/aggregate.sh
```

**Breaking Change Check:**
```bash
./scripts/breaking-change-check.sh
```

**Validate Transformed Spec:**
```bash
./scripts/validate.sh
```

### Related Issues

<!-- Link to related issues or JIRA tickets -->

- Closes #
- JIRA:

---

/label ~"spec-mapping" ~"partner-api"
/assign @platform-team
