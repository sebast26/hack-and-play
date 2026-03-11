# Spec for auth-state-management

branch: claude/feature/auth-state-management
figma_component (if used):| N/A

## Summary

Add a global auth state listener that tracks the Firebase Authentication state in real time and exposes it through a `useUser` hook. Any component or page can call `useUser()` to get the current user object (or `null` if logged out) without prop drilling or manual listener setup. This involves a React context provider wrapping the app, a `useUser` hook, and updates to existing components that need the current user.

## Functional Requirements

- Initialise Firebase Auth in `lib/firebase/config.ts` and export the `auth` instance
- Create an `AuthProvider` React context provider that:
  - Sets up a single `onAuthStateChanged` listener on mount
  - Stores the current Firebase `User` object (or `null`) in state
  - Provides the user value to all descendant components
  - Cleans up the listener on unmount
- Create a `useUser` hook that reads from the auth context and returns the current user (type: `User | null`)
- Wrap the root layout (`app/layout.tsx`) with `AuthProvider` so the hook is available in both `(public)` and `(dashboard)` route groups
- Update any existing components (e.g. `Navbar`, `Avatar`) that currently hard-code or mock the user to use `useUser` instead

## Figma Design Reference (only if referenced)

N/A

## Possible Edge Cases

- Hook called outside `AuthProvider` — should throw a clear error rather than returning `undefined`
- Firebase initialisation running on the server (SSR) — the provider must be a client component; the listener should only be set up client-side
- Brief loading window between mount and first `onAuthStateChanged` callback — the initial state should be `null` (not `undefined`) so consumers can safely check `if (user)` without an extra loading state for now
- Multiple components calling `useUser` simultaneously — all should read the same shared state from the single context

## Acceptance Criteria

- `useUser()` returns `null` when no user is signed in
- `useUser()` returns a Firebase `User` object when a user is signed in
- The auth listener is registered only once for the entire app (single provider at root)
- The Navbar (and any other component referencing the user) reads the user from `useUser` instead of any previous hard-coded value
- No prop drilling of the user object anywhere in the component tree
- The hook throws a descriptive error if used outside `AuthProvider`

## Open Questions

- Should a loading/pending state (e.g. `isLoading: boolean`) be included in the hook return value, or keep it minimal for now? include in hook return value
- Are there other components beyond `Navbar` and `Avatar` that will need to consume `useUser` immediately? no

## Testing Guidelines

Create a test file(s) in the ./tests folder for the new feature, and create meaningful tests for the following cases, without going too heavy:

- `useUser` returns `null` when rendered inside a provider with no authenticated user
- `useUser` returns the user object when the provider is given an authenticated user
- `useUser` throws when called outside of `AuthProvider`
- `AuthProvider` renders its children without errors
