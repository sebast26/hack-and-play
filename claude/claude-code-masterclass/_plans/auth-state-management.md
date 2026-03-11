# Plan: Auth State Management (`useUser` hook)

## Context

The app has Firebase configured (`lib/firebase/config.ts`) but no auth initialisation or global user state. Components that need the current user (Navbar, Avatar) have no way to access it. This plan adds a single `onAuthStateChanged` listener at the app root, exposed via a `useUser` hook returning `{ user, isLoading }`.

---

## Files to create / modify

### 1. `lib/firebase/config.ts` — add auth export
- Import `getAuth` from `firebase/auth`
- Call `getAuth(app)` and export as `auth`

### 2. `context/AuthContext.tsx` — new file (client component)
- Create `AuthContext` with `React.createContext`
- Context value type: `{ user: User | null, isLoading: boolean }`
- `AuthProvider` component:
  - Marks `"use client"`
  - Holds `user` (default `null`) and `isLoading` (default `true`) in state
  - Runs `onAuthStateChanged(auth, ...)` in a `useEffect` on mount; sets `user` and flips `isLoading` to `false` on first callback
  - Unsubscribes on unmount (effect cleanup)
  - Returns `<AuthContext.Provider value={{ user, isLoading }}>`
- `useUser` hook:
  - Reads context; throws `Error("useUser must be used within AuthProvider")` if context is `undefined`
  - Returns `{ user, isLoading }`

### 3. `app/layout.tsx` — wrap children with `AuthProvider`
- Import `AuthProvider` from `@/context/AuthContext`
- Wrap `{children}` inside `<AuthProvider>`
- Root layout stays a Server Component; `AuthProvider` is a client boundary

### 4. `components/Navbar/Navbar.tsx` — consume `useUser`
- Add `"use client"` directive (needs hook access)
- Import `useUser` from `@/context/AuthContext`
- Import `Avatar` from `@/components/Avatar`
- When `user` is present, render `<Avatar name={user.displayName ?? user.email ?? 'User'} />` in the nav
- When `isLoading`, render nothing / placeholder in that slot (keeps layout stable)

### 5. `tests/context/AuthContext.test.tsx` — new test file
Tests (using Vitest + React Testing Library, mock `firebase/auth`):
- `useUser` returns `{ user: null, isLoading: false }` when provider fires null
- `useUser` returns the user object when provider fires a user
- `useUser` throws when called outside `AuthProvider`
- `AuthProvider` renders children without errors

---

## Key conventions to follow (from CLAUDE.md)
- No semicolons
- No inline Tailwind (more than 1 class → `@apply` in CSS module)
- Path alias `@/*` for all imports
- New branch already created: `claude/feature/auth-state-management`

---

## Verification
1. `npm test` — all existing tests pass, new AuthContext tests pass
2. `npm run build` — no TypeScript errors
3. Manual: sign in via Firebase console token / DevTools → Navbar shows Avatar; sign out → Avatar disappears
