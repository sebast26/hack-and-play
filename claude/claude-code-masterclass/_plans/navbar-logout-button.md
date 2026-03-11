# Plan: Navbar Logout Button

## Context
The app has authenticated users via Firebase Auth, tracked globally via `useUser` (AuthContext). The Navbar already conditionally renders an Avatar when a user is logged in. We need to add a styled logout button next to the Avatar/Create Heist button that calls Firebase `signOut` when clicked. No redirect is needed — the UI reactively updates via `onAuthStateChanged`.

## Critical Files
- `components/Navbar/Navbar.tsx` — add logout button + handler
- `components/Navbar/Navbar.module.css` — add gradient pill button style
- `tests/components/Navbar.test.tsx` — extend with logout button tests

## Implementation Steps

### 1. Add logout button style to `Navbar.module.css`
Add a new `.logoutBtn` class using `@apply`:
- Gradient background left-to-right: `from-secondary to-primary` (pink → purple)
- Pill shape: `rounded-full`
- White text and icon
- Same padding as `.btn`: `px-4 py-2`
- `flex items-center gap-2` to align icon + label
- `font-semibold`

### 2. Update `Navbar.tsx`
- Import `signOut` from `firebase/auth`
- Import `auth` from `@/lib/firebase/config`
- Import `LogOut` icon from `lucide-react`
- Add `handleLogout` async function that calls `signOut(auth)` and catches errors (console.error only)
- Add a new `<li>` **before** the Create Heist `<li>` containing a `<button>` with:
  - `className={styles.logoutBtn}`
  - `onClick={handleLogout}`
  - `<LogOut size={16} />` icon + "Logout" text
  - Render condition: `!isLoading && user`

Final nav list order: `[Logout btn (auth only)] → [Create Heist] → [Avatar (auth only)]`

### 3. Extend `tests/components/Navbar.test.tsx`
Add three new test cases (keep existing two):
- **"renders logout button when user is authenticated"** — mock `useUser` returning a user, assert button with text `/logout/i` is in the document
- **"does not render logout button when not authenticated"** — mock `useUser` returning `null` user, assert button is absent
- **"calls signOut when logout button is clicked"** — mock `firebase/auth` `signOut`, simulate click, assert it was called

## Verification
1. Run `npm test` — all existing + new tests pass
2. Run `npm run dev` — log in, confirm logout button appears in Navbar to the left of "Create Heist"
3. Click Logout — confirm user state clears and button disappears
4. Run `npm run lint` — no lint errors
