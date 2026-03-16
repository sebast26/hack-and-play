# Plan: useHeist Hook

## Context

The heists dashboard page (`app/(dashboard)/heists/page.tsx`) is currently a static stub with three empty sections. This plan adds a `useHeist` hook that subscribes to real-time Firestore data and filters heists by mode (`'active'`, `'assigned'`, `'expired'`), then wires it into the page to display heist titles per section.

## Key findings

- **No `hooks/` directory** exists — must create it
- **`heistConverter`** already defined in `types/firestore/index.ts` — reuse for typed deserialization
- **`COLLECTIONS.HEISTS`** already defined in `types/firestore/index.ts`
- **`db`** exported from `@/lib/firebase/config`
- **`useUser()`** from `@/context/AuthContext` returns `{ user: User | null, isLoading: boolean }`
- **Heists page** is currently a plain server component — must become `"use client"` to use hooks
- **Firestore constraint**: `expired` filters by `where('deadline', '<=', now)` in Firestore only, then filters `finalStatus !== null` client-side — avoids needing a composite index for two inequality fields

---

## Steps

### 1. Create `hooks/useHeist.ts`

Return type: `{ heists: Heist[], loading: boolean, error: string | null }`

Logic:
- Call `useUser()` to get `user`
- `useEffect` depends on `[user, mode]`
- If `user` is null: set `loading = false`, return early — no Firestore call
- Build a typed query using `collection(db, COLLECTIONS.HEISTS).withConverter(heistConverter)`:
  - `'active'`: `where('assignedTo', '==', user.uid)` + `where('deadline', '>', new Date())`
  - `'assigned'`: `where('createdBy', '==', user.uid)` + `where('deadline', '>', new Date())`
  - `'expired'`: `where('deadline', '<=', new Date())` only; filter `finalStatus !== null` client-side
- Subscribe with `onSnapshot(q, successCb, errorCb)`
- Success: map `snapshot.docs.map(d => d.data())`, apply expired client-side filter, `setHeists`, `setLoading(false)`
- Error: `setError(err.message)`, `setLoading(false)`
- Return the unsubscribe function from `useEffect` for cleanup

### 2. Update `app/(dashboard)/heists/page.tsx`

- Add `"use client"` directive
- Import and call `useHeist` three times (once per mode)
- Under each section heading, render a list of `heist.title` values
- Show `"No heists yet"` when the array is empty
- Show a brief loading state while `loading` is true

---

## Files

| Action | Path |
|--------|------|
| Create | `hooks/useHeist.ts` |
| Modify | `app/(dashboard)/heists/page.tsx` |
| Create | `tests/hooks/useHeist.test.tsx` |

---

## Tests

**File:** `tests/hooks/useHeist.test.tsx`
Use `renderHook` from `@testing-library/react`.

Mocks:
- `@/lib/firebase/config` → `{ db: {} }`
- `firebase/firestore` → `onSnapshot`, `query`, `collection`, `where`
- `@/context/AuthContext` → `useUser` returning a fake user or null

`onSnapshot` mock pattern: capture the success callback, invoke it synchronously with a fake snapshot, return a mock unsubscribe fn.

Test cases:
1. Returns `loading: true` and empty array before snapshot fires
2. `'active'` — returns heists matching `assignedTo === user.uid` with future deadline
3. `'assigned'` — returns heists matching `createdBy === user.uid` with future deadline
4. `'expired'` — client-side filter: only heists with `finalStatus !== null` are returned
5. Does not call `onSnapshot` when `user` is null
6. Calls the unsubscribe fn on unmount
7. Sets `error` when `onSnapshot` fires its error callback

## Verification

1. `npm run dev` → `/heists` shows titles in all three sections in real time
2. Create a heist at `/heists/create` — appears live without refresh
3. `npm test` — all tests pass
