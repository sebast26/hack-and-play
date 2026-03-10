# Plan: Authentication Forms for Login & Signup

## Context

The `/login` and `/signup` pages are currently placeholders with just a title. This change adds functional auth forms with email/password fields, password visibility toggle, console-log-on-submit, and cross-navigation links per the spec in `_specs/auth-forms-login-signup.md`.

## Approach

Create a shared `AuthForm` component since both forms are nearly identical (same fields, same toggle, same submit behavior — only the title, button text, and nav link differ).

### 1. Create `components/AuthForm/` component

**`components/AuthForm/AuthForm.tsx`** — Client component (`"use client"`) with:
- Props: `mode: "login" | "signup"`
- State: `email`, `password`, `showPassword`
- Email input (`type="email"`, `required`)
- Password input with Eye/EyeOff toggle button (lucide-react)
- Submit button using existing `.btn` class
- Navigation link (Next.js `Link`) to the other page
- `onSubmit` handler: `preventDefault` + `console.log({ email, password })`

**`components/AuthForm/AuthForm.module.css`** — Form-specific styles using `@apply` (input styles, toggle button, form layout, link styles).

**`components/AuthForm/index.ts`** — Barrel export.

### 2. Update page files

**`app/(public)/login/page.tsx`** — Import `AuthForm`, render with `mode="login"`.

**`app/(public)/signup/page.tsx`** — Import `AuthForm`, render with `mode="signup"`.

Both pages keep the existing `center-content` / `page-content` / `form-title` wrapper structure.

### 3. Add tests

**`tests/components/AuthForm.test.tsx`** — Test both modes:
- Renders email and password inputs (login & signup)
- Password toggle switches input type between "password" and "text"
- Form submit calls `console.log` with `{ email, password }`
- Navigation link points to correct route (`/signup` from login, `/login` from signup)

## Key files

| File | Action |
|---|---|
| `components/AuthForm/AuthForm.tsx` | Create |
| `components/AuthForm/AuthForm.module.css` | Create |
| `components/AuthForm/index.ts` | Create |
| `app/(public)/login/page.tsx` | Edit |
| `app/(public)/signup/page.tsx` | Edit |
| `tests/components/AuthForm.test.tsx` | Create |

## Reusable pieces

- Global classes: `.form-title`, `.btn`, `.center-content`, `.page-content` (from `app/globals.css`)
- Icons: `Eye`, `EyeOff` from `lucide-react`
- `Link` from `next/link`
- Component pattern: follow `components/Avatar/` structure (CSS Modules with `@reference`, barrel export, no semicolons)

## Verification

1. `npm run lint` — no lint errors
2. `npm test` — all tests pass (existing + new)
3. `npm run dev` — visually check `/login` and `/signup` pages
4. Verify: fill form, submit, check console for `{ email, password }` log
5. Verify: toggle password visibility, navigate between pages via links
