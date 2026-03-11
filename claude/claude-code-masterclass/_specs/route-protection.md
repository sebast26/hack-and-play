# Spec for Route Protection

branch: claude/feature/route-protection

## Summary

Add route protection to the application so that pages in the `(public)` route group are only accessible to unauthenticated users, and pages in the `(dashboard)` route group are only accessible to authenticated users. While Firebase resolves the current auth state, group layouts display a simple loading indicator to prevent flash of unauthorized content.

## Functional Requirements

- The `(public)` group layout checks auth state via `useUser`; if the user is authenticated, redirect them to `/heists`
- The `(dashboard)` group layout checks auth state via `useUser`; if the user is unauthenticated, redirect them to `/login`
- Both group layouts render a simple centered loader/spinner while `useUser` is still resolving (i.e. loading state is true)
- Once auth state is resolved, redirects happen before any protected page content is rendered
- Redirects use Next.js client-side navigation (e.g. `useRouter`)
- No protected page content should flash before the redirect occurs

## Possible Edge Cases

- Firebase may take a moment to resolve auth state on first load — the loader must cover this window to avoid a content flash
- A user landing directly on a dashboard URL while unauthenticated should be reliably redirected to login
- A user landing directly on a public URL (e.g. login, signup) while already authenticated should be reliably redirected to the dashboard
- The loader should not persist indefinitely — only while the auth state is genuinely unresolved

## Acceptance Criteria

- Visiting `/login` or `/signup` while authenticated redirects to `/heists`
- Visiting `/heists` (or any dashboard page) while unauthenticated redirects to `/login`
- A loader is shown in both group layouts while auth state is loading
- No protected content is visible before the redirect fires
- After auth resolves, the correct page content renders without a second flash

## Open Questions

- Should unauthenticated users attempting to reach a specific dashboard URL be redirected back to that URL after login (i.e. redirect intent)? No.
- Is a spinner the preferred loader style, or a full-page skeleton/blank screen? spinner using the clock icon from the title.

## Testing Guidelines

Create a test file(s) in the ./tests folder for the new feature, and create meaningful tests for the following cases, without going too heavy:

- `(public)` layout renders a loader when auth state is loading
- `(public)` layout redirects to `/heists` when user is authenticated
- `(public)` layout renders children when user is unauthenticated
- `(dashboard)` layout renders a loader when auth state is loading
- `(dashboard)` layout redirects to `/login` when user is unauthenticated
- `(dashboard)` layout renders children when user is authenticated
