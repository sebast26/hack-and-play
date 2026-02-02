# 08 - Partner Portal

## Purpose

The Partner Portal is the unified destination for partners—combining public API documentation with authenticated self-service capabilities. Partners can explore the API without logging in, then authenticate to register applications and manage credentials.

This plan defines what the portal must provide, how public and authenticated sections interact, and how content stays synchronized with the Partner API.

---

## Role in the Platform

The Partner Portal serves as the **discovery, learning, and self-service layer**:

| Partner Need | How Portal Addresses It |
|--------------|------------------------|
| What can this API do? | Public API reference with all endpoints |
| How do I authenticate? | Getting started guide (public) |
| How do I accomplish X? | Use-case recipes with examples (public) |
| What changed in new version? | Changelog and migration guides (public) |
| How do I get credentials? | Application registration (authenticated) |
| How do I manage my apps? | Application dashboard (authenticated) |
| How much am I using? | Usage statistics (authenticated) |

The Partner Portal does NOT:
- Provide interactive API testing (deferred to later phase)
- Handle partner billing (separate concern)
- Handle support tickets (links to support channels)
- Manage partner business relationships

---

## Unified Portal Concept

### Why Unified (Not Separate)

The previous approach considered separate portals: a public Developer Portal for documentation and a Self-Service Portal for credential management. The unified approach is preferred because:

| Separate Portals | Unified Portal |
|------------------|----------------|
| Partners navigate between two sites | Single destination for everything |
| Different URLs to remember | One URL: `developers.company.com` |
| Inconsistent navigation/branding | Consistent experience throughout |
| "Where do I go for X?" confusion | Clear information architecture |
| Two codebases to maintain | Shared infrastructure and styling |

### Public vs. Authenticated Sections

The portal has both public and authenticated areas:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           PARTNER PORTAL                                    │
│                        developers.company.com                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                     PUBLIC (No Login Required)                       │  │
│   │                                                                      │  │
│   │   • API Reference (all endpoints, schemas, examples)                 │  │
│   │   • Getting Started Guide                                            │  │
│   │   • Use-Case Recipes                                                 │  │
│   │   • Changelog and Migration Guides                                   │  │
│   │   • Search                                                           │  │
│   │                                                                      │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                   AUTHENTICATED (Login Required)                     │  │
│   │                                                                      │  │
│   │   • My Applications (list, create, delete)                           │  │
│   │   • Credentials (view client_id, rotate secret)                      │  │
│   │   • Usage Dashboard (requests, rate limits)                          │  │
│   │   • Settings (default API version)                                   │  │
│   │                                                                      │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Public section**: No login required. Partners can explore everything about the API before committing to integration.

**Authenticated section**: Requires login via Auth0 (through Keycloak). Contains partner-specific functionality.

### Transition Between Sections

When partners navigate to authenticated features without being logged in:

1. Partner clicks "My Applications" or "Sign In"
2. Redirected to Keycloak login page
3. Keycloak shows Auth0 as login option
4. Partner authenticates with OLX credentials
5. Redirected back to portal, now authenticated
6. Authenticated section accessible

**Seamless experience**: The navigation bar shows "Sign In" when unauthenticated, changing to partner name/avatar when authenticated.

---

## Public Content

### API Reference

Auto-generated from the Partner API OpenAPI specification. Every endpoint is documented with:

| Section | Content | Source |
|---------|---------|--------|
| Title | Operation summary | OpenAPI `summary` |
| Description | What this endpoint does | OpenAPI `description` |
| URL | Full endpoint path | OpenAPI `path` |
| Method | HTTP method | OpenAPI `method` |
| Authentication | Required auth type | OpenAPI `security` |
| Parameters | Path, query, header params | OpenAPI `parameters` |
| Request body | Schema with examples | OpenAPI `requestBody` |
| Responses | Status codes and schemas | OpenAPI `responses` |
| Code examples | Multi-language examples | Generated from spec |

### Code Examples

Each endpoint includes code examples in multiple languages:

| Language | Format |
|----------|--------|
| **curl** | Command-line HTTP request |
| **Python** | Using `requests` library |
| **JavaScript/Node** | Using `fetch` or `axios` |
| **Java/Kotlin** | Using standard HTTP client |

**Example (curl)**:
```bash
curl -X GET "https://api.partner.company.com/partner/orders/123" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "X-Partner-API-Version: 2024-06-01"
```

**Example (Python)**:
```python
import requests

response = requests.get(
    "https://api.partner.company.com/partner/orders/123",
    headers={
        "Authorization": "Bearer YOUR_ACCESS_TOKEN",
        "X-Partner-API-Version": "2024-06-01"
    }
)
order = response.json()
```

**Example (JavaScript)**:
```javascript
const response = await fetch(
  "https://api.partner.company.com/partner/orders/123",
  {
    headers: {
      "Authorization": "Bearer YOUR_ACCESS_TOKEN",
      "X-Partner-API-Version": "2024-06-01"
    }
  }
);
const order = await response.json();
```

Code examples are auto-generated from the OpenAPI spec during documentation build.

### Getting Started Guide

Step-by-step guide for new partners:

1. **Overview**: What the Partner API offers
2. **Authentication**: How to obtain and use access tokens
3. **Making Your First Request**: Complete example from credential to response
4. **Error Handling**: Understanding error responses
5. **Rate Limits**: How limits work and what headers to expect
6. **Versioning**: How to specify API version

### Use-Case Recipes

Task-oriented guides showing how to accomplish specific goals:

| Recipe | Description |
|--------|-------------|
| Getting Your First Order | Authentication → fetch order → parse response |
| Listing Orders with Pagination | Pagination patterns, cursor handling |
| Creating an Order | Validation, required fields, error handling |
| Handling Webhooks | (if applicable) Webhook verification, processing |
| Error Recovery | Retry strategies, idempotency |

**Recipe structure**:
```markdown
# [Task Name]

## Overview
Brief description of what this recipe accomplishes.

## Prerequisites
- What you need before starting

## Steps

### Step 1: [First action]
Explanation and code example.

### Step 2: [Second action]
...

## Complete Example
Full working code combining all steps.

## Error Handling
Common errors and how to handle them.
```

### Changelog

Organized by version, newest first:

```markdown
# Changelog

## Version 2025-01-15 (Active)
Released: January 15, 2025

### New Features
- Added `GET /partner/invoices` endpoint
- Added `metadata` field to order responses

### Improvements
- Improved error messages for validation failures

### Breaking Changes
None in this version.

---

## Version 2024-06-01
Released: June 1, 2024

### Breaking Changes
- Restructured customer data from flat fields to nested object

[Migration Guide →](/docs/migration/2024-01-15-to-2024-06-01)
```

### Versioned Documentation

Partners can view documentation for any supported version:

- Version selector in navigation
- Each version has complete API reference
- Version lifecycle status shown (Active, Supported, Deprecated)
- Migration guides linked from deprecated versions

### Search

Partners can search across all public content:

- Endpoint names and descriptions
- Recipe titles and content
- Schema field names
- Changelog entries

Search powered by Docusaurus built-in search (Algolia DocSearch or local search).

### Documentation Feedback

Each documentation page includes a **"Report an issue"** link that allows partners to report problems:

| Element | Details |
|---------|---------|
| Link location | Footer of each documentation page |
| Link text | "Report an issue with this page" |
| Action | Opens pre-filled support ticket (JIRA) or email |
| Pre-filled context | Page URL, page title, current API version |

This provides low-friction feedback while integrating with existing support workflows.

---

## Authenticated Features

### My Applications

Partners see all their registered applications:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  My Applications                                            [+ New App]     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Production Integration                              Standard Tier   │   │
│  │  Client ID: acme-corp_a1b2c3d4                                      │   │
│  │  Created: Jan 15, 2024                                              │   │
│  │                                                                      │   │
│  │  [View Credentials]  [Rotate Secret]  [Delete]                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Staging Environment                                 Standard Tier   │   │
│  │  Client ID: acme-corp_e5f6g7h8                                      │   │
│  │  Created: Jan 20, 2024                                              │   │
│  │                                                                      │   │
│  │  [View Credentials]  [Rotate Secret]  [Delete]                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Application Registration

Creating a new application:

1. Click "+ New App"
2. Enter application name (e.g., "Production Integration")
3. Select default API version
4. Submit
5. Receive client_id and client_secret (secret shown once)
6. Confirmation with next steps

### Credentials View

Partners can view their credentials:

- client_id always visible
- client_secret shown only on creation or rotation
- Copy-to-clipboard buttons
- Rotation instructions

### Secret Rotation

Partners can rotate secrets:

1. Click "Rotate Secret"
2. Warning: "Current secret will be immediately invalidated"
3. Confirm
4. New secret displayed (shown once)
5. Partner must update their integration

### Usage Dashboard

Partners see their API usage:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  Usage: Production Integration                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Today's Usage                                                              │
│  ├── Requests: 12,450 / 50,000 (24.9%)                                     │
│  └── Rate: 3.2 req/s (limit: 10 req/s)                                     │
│                                                                             │
│  [==============================------------------------] 24.9%            │
│                                                                             │
│  Last 7 Days                                                                │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │    ▄                                                                │   │
│  │   ▄█▄    ▄                                                          │   │
│  │  ▄███▄  ▄█▄  ▄▄▄                                                    │   │
│  │ ▄█████▄▄███▄▄███▄                                                   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│    Mon   Tue   Wed   Thu   Fri   Sat   Sun                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Usage data source**: Data is sourced from **NewRelic with caching**. The portal backend queries NewRelic API and caches results (e.g., 5-minute TTL) to improve responsiveness and reduce API calls. See `05-security-auth.md` for details.

### Settings

Partners can configure:

- Default API version for each application
- Notification preferences (if implemented)

---

## Authentication Flow

### Keycloak Integration

Portal uses Keycloak for authentication, which federates to Auth0:

```
Partner              Portal                  Keycloak                 Auth0
   │                   │                        │                       │
   ├── Sign In ───────►│                        │                       │
   │                   ├── Redirect ───────────►│                       │
   │                   │                        ├── Redirect ──────────►│
   │                   │                        │                       │
   │◄──────────────────────────────────── Auth0 Login ──────────────────┤
   │                   │                        │                       │
   ├── Credentials ────────────────────────────────────────────────────►│
   │                   │                        │                       │
   │                   │                        │◄── Token ─────────────┤
   │                   │                        │                       │
   │                   │◄── Session ────────────┤                       │
   │◄── Authenticated ─┤                        │                       │
```

### Session Management

| Aspect | Behavior |
|--------|----------|
| Session duration | 24 hours (aligned with Keycloak configuration) |
| Refresh | Silent refresh before expiration |
| Logout | Clears portal session and Keycloak session |
| Expired session | Redirect to login, then back to intended page |

### Protected Routes

Authenticated sections are protected:

- `/dashboard/*` - Application management
- `/settings` - Partner settings

Accessing protected route while unauthenticated triggers login flow with return URL preserved.

---

## Technical Architecture

### Repository Structure

The portal lives in a **separate repository** from the governance repo:

| Repository | Contents |
|------------|----------|
| `partner-api-portal` | Portal source code (Docusaurus + React dashboard) |
| `partner-api-governance` | Specs, transformations, CI templates, changelogs |

**Why separate**:
- Independent release cycles
- Different concerns (documentation vs. spec management)
- Portal team can deploy without affecting governance pipeline
- Cleaner ownership boundaries

The portal repository pulls OpenAPI specs and changelogs from the governance repo during build.

### Hybrid Approach

The portal combines two technologies:

| Component | Technology | Purpose |
|-----------|------------|---------|
| Documentation | Docusaurus | Static docs, versioning, search |
| Dashboard | React + Keycloak JS | Authenticated features |

**Why hybrid (not fully custom)**:
- Docusaurus excels at documentation (versioning, search, Markdown)
- React provides flexibility for dashboard features
- Shared styling ensures consistent look
- Both deployed together under same domain

### Site Structure

```
developers.company.com/
├── /                           # Landing page
├── /docs/                      # Documentation (Docusaurus)
│   ├── /getting-started/
│   ├── /api-reference/
│   ├── /recipes/
│   └── /changelog/
├── /dashboard/                 # Authenticated (React)
│   ├── /applications/
│   ├── /usage/
│   └── /settings/
└── /login                      # Keycloak redirect
```

### Navigation Integration

Unified navigation bar across both sections:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  [Logo] Partner API    Docs  API Reference  Recipes    Dashboard    [Sign In]│
└─────────────────────────────────────────────────────────────────────────────┘
```

- "Docs", "API Reference", "Recipes" → Docusaurus pages
- "Dashboard" → React app (triggers login if needed)
- "Sign In" → Keycloak login (changes to user menu when authenticated)

### Mobile Responsiveness

The portal has differentiated mobile support:

| Section | Mobile Support |
|---------|----------------|
| Documentation (Docusaurus) | **Fully responsive** - readable and navigable on mobile |
| Dashboard (React) | **Desktop only** - requires larger screen |

**Rationale**: Partners commonly read documentation on mobile (quick reference while debugging), but credential management is typically done from a workstation. This focuses development effort where it provides the most value.

Dashboard pages on mobile display a message suggesting desktop access for the best experience.

### Build and Deployment

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           BUILD PIPELINE                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐      │
│   │  OpenAPI Spec   │────►│    Docusaurus   │────►│  Static Files   │      │
│   │  (per version)  │     │      Build      │     │  /docs/*        │      │
│   └─────────────────┘     └─────────────────┘     └─────────────────┘      │
│           │                       ▲                       │                 │
│           │               ┌───────┴───────┐               │                 │
│           │               │   Markdown    │               │                 │
│           │               │   Content     │               │                 │
│           │               └───────────────┘               │                 │
│           │                                               ▼                 │
│           │                                       ┌─────────────────┐       │
│           │                                       │    Combined     │       │
│   ┌─────────────────┐     ┌─────────────────┐    │    Deployment   │       │
│   │  Dashboard      │────►│   React Build   │───►│                 │       │
│   │  Source         │     │                 │    │   Kubernetes    │       │
│   └─────────────────┘     └─────────────────┘    │   (nginx)       │       │
│                                                   └─────────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Hosting

Portal deployed as static files + API proxy:

| Component | Hosting |
|-----------|---------|
| Static files (docs, dashboard) | nginx in Kubernetes |
| Dashboard API calls | Proxy to middleware service |
| Keycloak authentication | Redirect to Keycloak service |

---

## Analytics

The portal tracks anonymous usage analytics to understand documentation effectiveness:

### What's Tracked

| Metric | Purpose |
|--------|---------|
| Page views | Which documentation pages are popular |
| Search queries | What partners are looking for |
| Version usage | Which API versions are being referenced |
| Navigation patterns | How partners move through docs |

### What's NOT Tracked

| Data | Reason |
|------|--------|
| User identity | Privacy - no user identification |
| Partner organization | Privacy - no business identification |
| Dashboard activity | Contains sensitive credential operations |
| IP addresses | Privacy - not stored |

### Implementation

| Aspect | Approach |
|--------|----------|
| Tool | Privacy-focused analytics (e.g., Plausible, or server-side logging) |
| Cookie consent | Not required (no cookies, no user tracking) |
| Data retention | Aggregated metrics only; no individual session data |
| Access | Platform team dashboard for documentation improvements |

**Why anonymous**: Partners trust the portal with their credentials. Tracking their documentation usage without identification respects that trust while still providing insights to improve the documentation.

---

## Content Management

### Documentation Updates

| Content Type | Update Process |
|--------------|----------------|
| API Reference | Auto-generated when spec changes |
| Getting Started | Manual edit, PR review |
| Recipes | Manual edit, PR review |
| Changelog | Manual edit for context; auto-detect for change list |

### Documentation Ownership

| Content | Owner |
|---------|-------|
| API Reference | Aggregation pipeline (auto-generated) |
| Getting Started | Platform team |
| Recipes | Platform team (with domain team input) |
| Changelog | Platform team |
| Dashboard | Platform team |

### Version Documentation

When new Partner API version is released:

1. New spec added to docs
2. Version appears in version switcher
3. Previous version documentation preserved
4. Changelog updated with new version entry

---

## Future Enhancements

The following features are explicitly **deferred to later phases** based on partner feedback:

### Interactive API Testing ("Try It")

A "try it" feature allowing partners to make API calls directly from the documentation is not included in the initial release. This feature adds significant complexity (sandbox environment, credential handling in browser) and will be evaluated based on partner demand.

**If implemented later**:
- Would use a dedicated sandbox environment
- Partners would authenticate through the portal
- Requests would go to sandbox, not production
- Results displayed inline in documentation

### Additional Language Examples

If partner demographics indicate demand, additional language examples may be added:
- Go
- Ruby
- PHP
- C#/.NET

---

## Integration Points

### With Spec Aggregation Pipeline

- Pipeline produces OpenAPI spec per version
- Spec triggers documentation rebuild
- Docusaurus generates API reference from spec

### With Authentication Service

- Portal uses Keycloak for authentication
- Dashboard calls middleware for application management
- Keycloak session determines authenticated state

### With Middleware Service

Dashboard communicates with middleware for:
- Listing partner's applications
- Creating new applications
- Rotating credentials
- Deleting applications
- Fetching usage data (via NewRelic)

### With Governance Repository

- OpenAPI specs pulled during build
- Changelog content from governance repo
- Migration guides from governance repo

---

## Success Criteria

Partner Portal is complete when:

| Criteria | Validation |
|----------|------------|
| API reference accessible | All endpoints documented, searchable |
| Multi-language examples | curl, Python, JavaScript, Java examples present |
| Recipes available | At least 5 use-case recipes published |
| Version switching works | Can view docs for different versions |
| Search works | Can find endpoints, recipes, guides |
| Login works | Partner can authenticate via Auth0 |
| App registration works | Partner can create application, get credentials |
| Credential rotation works | Partner can rotate secret |
| App deletion works | Partner can delete application |
| Navigation unified | Seamless movement between docs and dashboard |
| Branding consistent | Logo, colors match company identity |
| Docs responsive | Documentation readable on mobile devices |
| Feedback mechanism works | "Report issue" link opens pre-filled ticket |
| Analytics collecting | Page views and search queries tracked anonymously |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1 | 2026-01-28 | Initial draft (as Developer Portal) |
| 0.2 | 2026-01-29 | Revised as unified Partner Portal combining documentation and self-service dashboard; renamed from 08-developer-portal.md |
| 0.3 | 2026-01-30 | Resolved open questions: separate repository; curl + Python/JS/Java examples; NewRelic for usage data; "Report issue" feedback link; interactive testing deferred; responsive docs with desktop-only dashboard; basic anonymous analytics |
