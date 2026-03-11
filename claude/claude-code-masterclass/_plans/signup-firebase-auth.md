# Plan: Signup Firebase Auth with Codename Generation

## Context
The signup form (`AuthForm` in `app/(public)/signup/`) currently only logs credentials to the console. This plan wires it up to Firebase Auth, generates a random PascalCase codename for each new user, sets it as their `displayName`, and creates a Firestore `users/{uid}` document storing `id` and `codename` (no email).

---

## Files to Modify / Create

| Action | Path |
|--------|------|
| Modify | `lib/firebase/config.ts` — add `db` (Firestore) export |
| Create | `lib/codename.ts` — pure codename generator utility |
| Create | `lib/firebase/signup.ts` — signup orchestration function |
| Modify | `components/AuthForm/AuthForm.tsx` — call signup, handle errors, redirect |
| Create | `tests/lib/codename.test.ts` — codename generator tests |
| Create | `tests/lib/firebase/signup.test.ts` — signup flow tests |

---

## Step-by-Step Implementation

### 1. Add Firestore export to `lib/firebase/config.ts`
- Import `getFirestore` from `firebase/firestore`
- Initialise and export `db` alongside the existing `auth` export

### 2. Create `lib/codename.ts`
- Define three distinct word arrays (e.g. `ADJECTIVES`, `NOUNS`, `VERBS`) — at least 20 words each
- Export `generateCodename(): string` that picks one word from each array at random and returns them concatenated in PascalCase (first letter uppercased, rest lowercase for each word)

### 3. Create `lib/firebase/signup.ts`
- Export `async function signUpUser(email: string, password: string): Promise<void>`
- Steps inside:
  1. `createUserWithEmailAndPassword(auth, email, password)` → get `UserCredential`
  2. `generateCodename()` → `codename`
  3. `updateProfile(user, { displayName: codename })`
  4. `setDoc(doc(db, "users", user.uid), { id: user.uid, codename })` — uses UID as document ID
- Let errors propagate to the caller (no swallowing)

### 4. Modify `components/AuthForm/AuthForm.tsx`
- Add `error` state (`string | null`, init `null`)
- Add `isLoading` state (`boolean`, init `false`)
- When `mode === "signup"`:
  - In `handleSubmit`: set `isLoading(true)`, call `signUpUser(email, password)`
  - On success: call `router.push("/heists")` (import `useRouter` from `next/navigation`)
  - On error: set `error` to the Firebase error message
  - Always set `isLoading(false)` in finally
- Render error message below the form inputs (use existing `styles` or a minimal inline class)
- Disable the submit button while `isLoading` is true
- Login mode (`mode === "login"`) is left unchanged — the submit handler for login is out of scope

### 5. Tests — `tests/lib/codename.test.ts`
- `generateCodename()` returns a non-empty string
- Result is valid PascalCase (each of 3 words starts with uppercase, rest lowercase, no separators)
- Calling it multiple times returns different results (run 10 times, expect > 1 unique value)

### 6. Tests — `tests/lib/firebase/signup.test.ts`
- Mock `firebase/auth`: `createUserWithEmailAndPassword`, `updateProfile`
- Mock `firebase/firestore`: `setDoc`, `doc`
- Mock `@/lib/firebase/config`: `auth: {}`, `db: {}`
- Mock `@/lib/codename`: `generateCodename` returns a fixed string (e.g. `"SwiftCrimsonFox"`)
- Test: `signUpUser` calls `createUserWithEmailAndPassword` with provided email/password
- Test: on success, `updateProfile` is called with `{ displayName: "SwiftCrimsonFox" }`
- Test: on success, `setDoc` is called with a doc ref containing uid, and `{ id: uid, codename: "SwiftCrimsonFox" }` — no email field
- Test: if `createUserWithEmailAndPassword` throws, the error propagates

---

## Verification
1. Run `npm test` — all existing tests and new tests pass
2. Run `npm run dev`, navigate to `/signup`, create a new account
3. Verify in Firebase Console → Authentication that the user exists with a `displayName` set
4. Verify in Firebase Console → Firestore → `users` collection that a document exists with `id` and `codename` fields and no `email` field
5. Verify browser redirects to `/heists` after signup
6. Verify that submitting with an already-used email shows an error in the form
