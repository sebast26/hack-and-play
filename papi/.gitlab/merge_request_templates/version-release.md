## Partner API Version Release

<!-- Use this template when releasing a new version of the Partner API -->

### Release Information

- **Version**: `v0.0.0`
- **Release Date**:
- **Release Manager**:

### Release Type

<!-- Check one -->

- [ ] **Major** - Breaking changes, new major features
- [ ] **Minor** - New features, backward compatible
- [ ] **Patch** - Bug fixes, documentation updates

### Changes Included

#### New Features
<!-- List new endpoints, capabilities, etc. -->

-

#### Improvements
<!-- List enhancements to existing functionality -->

-

#### Bug Fixes
<!-- List fixes -->

-

#### Breaking Changes
<!-- List any breaking changes - requires Major version bump -->

- [ ] No breaking changes in this release
-

#### Deprecations
<!-- List any newly deprecated features -->

-

### Pre-Release Checklist

**Validation:**
- [ ] All CI checks passing
- [ ] Aggregation pipeline succeeds
- [ ] Breaking change check passes
- [ ] Spectral validation passes
- [ ] Integration tests pass

**Documentation:**
- [ ] Changelog updated
- [ ] API documentation updated
- [ ] Migration guide created (if breaking changes)

**Communication:**
- [ ] Release notes drafted
- [ ] Partner notification prepared
- [ ] Internal announcement prepared

### Release Checklist

**Platform Team Actions:**
- [ ] Update baseline spec (`governance/baseline/partner-api.yaml`)
- [ ] Tag CI templates with version (`v0.0.0`)
- [ ] Create Git tag for release
- [ ] Deploy to staging
- [ ] Verify staging deployment
- [ ] Deploy to production
- [ ] Verify production deployment
- [ ] Send partner notification
- [ ] Update public documentation

### Post-Release

- [ ] Monitor for issues
- [ ] Respond to partner feedback
- [ ] Close related issues

### Rollback Plan

<!-- Describe how to rollback if issues are found -->

1. Revert to previous baseline spec
2. Redeploy gateway with previous configuration
3. Notify partners of temporary rollback

### Files Changed

- `governance/baseline/partner-api.yaml`
- Version tags on `ci-templates/*`

---

/label ~"release" ~"partner-api"
/assign @platform-team
/milestone %"Partner API vX.X"
