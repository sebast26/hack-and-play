# Plan: Create Heist Form

## Context

The `/heists/create` page currently renders a shell with only a title. This plan wires up a real form that lets authenticated users create a heist document in Firestore and redirects them to `/heists` on success. The `CreateHeistInput` type already defines the exact document shape; user identity comes from `AuthContext`; assignee data is fetched from the `users` Firestore collection (shape: `{ id, codename }`).

## Key findings

- `AuthContext` exposes `user` (Firebase `User`) — `user.uid` = `createdBy`, `user.displayName` = `createdByCodename`
- `signup.ts` confirms user docs are stored as `{ id: user.uid, codename }` in the `users` collection
- `CreateHeistInput.deadline` comment: **48 hours after `createdAt`**
- `db` is already exported from `@/lib/firebase/config`
- `COLLECTIONS` in `types/firestore/index.ts` only has `HEISTS` — needs `USERS` added
- Pattern: client component with `useState`, async `handleSubmit`, try/catch/finally, `router.push` (mirrors `AuthForm`)
- Component convention: `components/<Name>/<Name>.tsx` + `<Name>.module.css` + `index.ts` barrel

---

## Steps

### 1. Add `USERS` to `COLLECTIONS` and `UserDoc` type
**File:** `types/firestore/index.ts`

- Add `USERS: "users"` to the `COLLECTIONS` constant
- Export a `UserDoc` interface: `{ id: string; codename: string }`

### 2. Create `CreateHeistForm` component
**Files to create:**
- `components/CreateHeistForm/CreateHeistForm.tsx`
- `components/CreateHeistForm/CreateHeistForm.module.css`
- `components/CreateHeistForm/index.ts`

**Behaviour:**
- `"use client"` directive
- On mount: fetch all docs from `COLLECTIONS.USERS` via `getDocs(collection(db, COLLECTIONS.USERS))`, map to `UserDoc[]`, store in state
- Form fields (controlled inputs):
  - `title` — text input, required
  - `description` — textarea, required
  - `assignedTo` — `<select>` of users; value encodes `id|codename`, required
- On submit (`handleSubmit`):
  1. Clear error, set `isLoading = true`
  2. Parse selected assignee into `assignedTo` + `assignedToCodename`
  3. Build `CreateHeistInput`:
     - `createdBy`: `user.uid`
     - `createdByCodename`: `user.displayName`
     - `createdAt`: `serverTimestamp()`
     - `deadline`: `new Date(Date.now() + 48 * 60 * 60 * 1000)`
     - `finalStatus`: `null`
  4. `await addDoc(collection(db, COLLECTIONS.HEISTS), payload)`
  5. On success: `router.push('/heists')`
  6. On error: set error message, keep form editable
  7. `finally`: set `isLoading = false`
- Submit button disabled when `isLoading` or users list is empty
- Inline error `<p>` when error state is set
- All multi-class styling via CSS module with `@apply`

### 3. Update `CreateHeistPage`
**File:** `app/(dashboard)/heists/create/page.tsx`

Replace stub — render `<CreateHeistForm />` inside the existing page shell. Page stays a server component; client boundary lives inside the component.

---

## Files

| Action | Path |
|--------|------|
| Modify | `types/firestore/index.ts` |
| Create | `components/CreateHeistForm/CreateHeistForm.tsx` |
| Create | `components/CreateHeistForm/CreateHeistForm.module.css` |
| Create | `components/CreateHeistForm/index.ts` |
| Modify | `app/(dashboard)/heists/create/page.tsx` |

---

## Tests

**File:** `tests/components/CreateHeistForm.test.tsx`

Mocks needed:
- `next/navigation` → `useRouter` returning `{ push: mockPush }`
- `@/context/AuthContext` → `useUser` returning `{ uid: 'u1', displayName: 'Viper' }`
- `firebase/firestore` → `getDocs`, `addDoc`, `collection`, `serverTimestamp`
- `@/lib/firebase/config` → `db`

Test cases:
1. Renders title input, description textarea, and assignee dropdown
2. Dropdown is populated with users returned by `getDocs`
3. Submitting with valid inputs calls `addDoc` with a payload matching `CreateHeistInput` shape
4. After successful `addDoc`, router pushes to `/heists`
5. Submit button is disabled while submission is in progress
6. Displays an error message when `addDoc` rejects; form remains interactive

## Verification

1. `npm run dev` → navigate to `/heists/create`, verify form renders and users populate the dropdown
2. Submit the form → confirm a new document appears in Firestore console under `heists`
3. Confirm redirect to `/heists` after submit
4. `npm test` → all existing and new tests pass
