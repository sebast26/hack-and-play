# Branch Protection Setup Guide

This document describes the required GitLab project settings for the Partner API repository.

## Branch Protection Rules

### Main Branch (`main` or `master`)

Configure the following protection rules in GitLab:

**Settings → Repository → Protected branches**

| Setting | Value |
|---------|-------|
| Branch | `main` |
| Allowed to merge | Maintainers |
| Allowed to push | No one |
| Allowed to force push | No |
| Code owner approval | Required |

### How to Configure

1. Go to **Settings → Repository → Protected branches**
2. Select `main` branch
3. Set "Allowed to merge" to **Maintainers**
4. Set "Allowed to push and merge" to **No one**
5. Enable "Require approval from code owners"

## Merge Request Approval Rules

### Required Approvals

Configure approval rules in **Settings → Merge requests → Approval rules**:

#### Rule 1: Platform Team Approval (Default)

| Setting | Value |
|---------|-------|
| Name | Platform Team Approval |
| Approvals required | 1 |
| Approvers | @platform-team |
| Target branch | All branches |

#### Rule 2: Governance Changes (Elevated)

For changes to critical governance files, require additional approval:

| Setting | Value |
|---------|-------|
| Name | Governance Changes |
| Approvals required | 2 |
| Approvers | @platform-team-leads |
| Target branch | main |
| File patterns | `governance/manifest.yaml`, `governance/spec-mappings/*`, `governance/baseline/*` |

### CODEOWNERS Integration

Enable CODEOWNERS approval:

1. Go to **Settings → Merge requests**
2. Under "Merge request approvals":
   - Enable "Require approval from code owners"
   - Set "Code owner approval" to "Required"

## Merge Request Settings

### General MR Settings

**Settings → Merge requests → Merge method**

| Setting | Value |
|---------|-------|
| Merge method | Merge commit |
| Squash commits | Encourage (allow but don't require) |
| Delete source branch | Enable by default |

### Merge Checks

**Settings → Merge requests → Merge checks**

| Check | Status |
|-------|--------|
| Pipelines must succeed | Enabled |
| All threads must be resolved | Enabled |
| Status checks must succeed | Enabled |

## Pipeline Configuration

### Required CI Jobs

The following jobs must pass before merge:

| Job | Stage | Purpose |
|-----|-------|---------|
| `validate` | validate | Validates OpenAPI specs |
| `aggregate` | aggregate | Runs aggregation pipeline |
| `breaking-change` | breaking-change | Checks for breaking changes |
| `lint` | lint | Runs Spectral linting |

### Branch-Specific Jobs

| Job | Runs On |
|-----|---------|
| `build` | main branch only |
| `test` | main branch only |

## Access Levels

### Role Definitions

| Role | Access Level | Permissions |
|------|--------------|-------------|
| Platform Team | Maintainer | Full access, approve MRs, merge to main |
| Internal Teams | Developer | Create branches, submit MRs |
| Viewers | Reporter | View repository, no changes |

### Group Membership

Create GitLab groups for access management:

| Group | Purpose |
|-------|---------|
| `@platform-team` | Platform engineers with full access |
| `@platform-team-leads` | Tech leads for elevated approvals |
| `@partner-api-contributors` | All teams contributing to Partner API |

## Automation Accounts

### CI/CD Service Account

For automated pipeline operations:

| Setting | Value |
|---------|-------|
| Username | `partner-api-bot` |
| Role | Maintainer |
| Purpose | Automated commits (if needed) |

**Note**: The bot account should only be used for automated operations triggered by the pipeline, never for bypassing review requirements.

## Checklist

Use this checklist when setting up a new Partner API repository:

- [ ] Create `main` branch protection
- [ ] Configure "No direct push" to main
- [ ] Enable CODEOWNERS file
- [ ] Set up approval rules
- [ ] Configure merge request settings
- [ ] Enable pipeline requirements
- [ ] Create access groups
- [ ] Add CODEOWNERS file to repository
- [ ] Test MR workflow with a sample change

## Verification

To verify branch protection is working:

1. Try to push directly to main (should fail)
2. Create an MR without approvals (should block merge)
3. Create an MR modifying governance files (should require 2 approvals)
4. Create an MR with failing pipeline (should block merge)
