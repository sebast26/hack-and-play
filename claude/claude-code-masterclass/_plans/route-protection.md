# Plan: Route Protection with Auth-Aware Layouts

## Context

The app has two Next.js route groups — `(public)` for unauthenticated pages and `(dashboard)` for authenticated pages — but currently neither enforces access control. This plan adds route protection to both group layouts using the existing `useUser` hook, and shows a spinner (using the `Clock8` Lucide icon, matching the title) while Firebase resolves auth state.

## Files to Modify

| File | Change |
|---|---|
| `app/(public)/layout.tsx` | Convert to client component; add auth guard |
| `app/(dashboard)/layout.tsx` | Convert to client component; add auth guard |
| `app/globals.css` | Add `.auth-loader` and `@keyframes spin-clock` styles |

## Files to Create

| File | Purpose |
|---|---|
| `tests/app/(public)/layout.test.tsx` | Tests for public layout auth guard |
| `tests/app/(dashboard)/layout.test.tsx` | Tests for dashboard layout auth guard |

## Implementation Steps

### 1. Add spinner styles to `app/globals.css`

Add a `.auth-loader` utility class and a CSS `@keyframes` animation for rotating the clock icon. The loader should be a full-viewport centered container (reusing the `.center-content` pattern) with the `Clock8` icon spinning inside it.

### 2. Update `app/(public)/layout.tsx`

- Add `"use client"` directive
- Import `useUser` from `@/context/AuthContext`, `useRouter` from `next/navigation`, and `Clock8` from `lucide-react`
- Logic:
  - `isLoading === true` → render the `.auth-loader` spinner
  - `user !== null` → call `router.replace("/heists")` and render the spinner (no flash)
  - Otherwise → render `<main className="public">{children}</main>`

### 3. Update `app/(dashboard)/layout.tsx`

- Add `"use client"` directive
- Import `useUser`, `useRouter`, and `Clock8`
- Logic:
  - `isLoading === true` → render the `.auth-loader` spinner
  - `user === null` → call `router.replace("/login")` and render the spinner
  - Otherwise → render `<><Navbar /><main>{children}</main></>`

### 4. Write tests — `tests/app/(public)/layout.test.tsx`

Mock `useUser` (via `vi.mock("@/context/AuthContext")`) and `next/navigation` (mock `useRouter` / `replace`).

- renders loader when `isLoading: true`
- redirects to `/heists` when `user` is set (not loading)
- renders children when user is `null` and not loading

### 5. Write tests — `tests/app/(dashboard)/layout.test.tsx`

Same mocking pattern.

- renders loader when `isLoading: true`
- redirects to `/login` when `user` is `null` (not loading)
- renders children (and Navbar) when user is authenticated and not loading

## Key Reuse

- `useUser` hook → `context/AuthContext.tsx` — returns `{ user, isLoading }`
- `Clock8` icon → already used in `components/Navbar/Navbar.tsx`
- `.center-content` CSS class (or same pattern) → `app/globals.css`
- Mock patterns from `tests/components/Navbar.test.tsx` and `tests/context/AuthContext.test.tsx`

## Verification

1. Run `npm test` — all existing tests should still pass; new layout tests should pass
2. Run `npm run dev` and manually test:
   - Visit `/login` while logged in → should redirect to `/heists`
   - Visit `/heists` while logged out → should redirect to `/login`
   - On first load, a brief spinner should appear before redirect or page content
3. Run `npm run lint` — no lint errors
