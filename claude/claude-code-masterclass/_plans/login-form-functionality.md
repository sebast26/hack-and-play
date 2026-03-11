# Plan: Login Form Functionality

## Context

The login form at `app/(public)/login/` renders `<AuthForm mode="login" />`, but the submit handler is currently a stub — it only `console.log`s credentials and returns early. This plan wires up real Firebase Auth sign-in, shows an auto-dismissing success message on success, and a readable error message on failure (no redirect).

---

## What needs to change

### 1. Create `lib/firebase/login.ts`

New file mirroring `lib/firebase/signup.ts`. Calls `signInWithEmailAndPassword(auth, email, password)` from `firebase/auth` and returns `void`. No Firestore interaction needed (login doesn't change user data).

### 2. Update `components/AuthForm/AuthForm.tsx`

- Import `loginUser` from `@/lib/firebase/login`
- Add `success` state: `const [success, setSuccess] = useState(false)`
- Replace the stub `if (isLogin)` block with a real async try/catch:
  - `setError(null)`, `setSuccess(false)`, `setIsLoading(true)`
  - `await loginUser(email, password)`
  - On success: `setSuccess(true)` then schedule auto-dismiss with `setTimeout(() => setSuccess(false), 3000)`
  - On error: set human-readable error message from Firebase error code
- Render success message in JSX (styled via CSS module), mutually exclusive with the error paragraph

### 3. Update `components/AuthForm/AuthForm.module.css`

Add a `.success` class for the success message styling (mirrors the existing `.error` style, using the `success` theme colour `#05DF72`).

### 4. Update `tests/components/AuthForm.test.tsx`

- Remove the existing `"logs email and password on submit"` test (behaviour is being replaced)
- Add a `vi.mock` for `@/lib/firebase/login` (mirroring the existing `signup` mock at the top)
- Add login-mode tests:
  - Shows a success message after valid credentials are submitted (mock `loginUser` resolves)
  - Does not show the success message when there is an error
  - Shows an error message when `loginUser` rejects
  - Submit button is disabled while loading
  - Previous error is cleared on resubmit

---

## Critical files

| File | Action |
|---|---|
| `lib/firebase/login.ts` | Create |
| `components/AuthForm/AuthForm.tsx` | Modify (lines 24–29 stub + JSX) |
| `components/AuthForm/AuthForm.module.css` | Modify (add `.success`) |
| `tests/components/AuthForm.test.tsx` | Modify (replace stub test, add login tests) |

### Reference

- `lib/firebase/signup.ts` — pattern to mirror for `loginUser`
- `lib/firebase/config.ts` — exports `auth` to pass to Firebase
- `firebase/auth` — use `signInWithEmailAndPassword`

---

## Verification

1. Run `npm test` — all existing tests pass, new login tests pass
2. Run `npm run dev`, go to `/login`, submit valid credentials → success message appears then fades after ~3 s
3. Submit invalid credentials → readable error message shown
4. Submit button disabled during in-flight request
