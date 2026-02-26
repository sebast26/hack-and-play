# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
npm run dev      # Start dev server at http://localhost:3000
npm run build    # Production build
npm run lint     # ESLint
npm test         # Run all tests (Vitest)
```

To run a single test file:
```bash
npx vitest tests/components/Navbar.test.tsx
```

## Frameworks & Libraries

- **Next.js 16** (App Router) — React framework
- **React 19** — UI library
- **Tailwind CSS 4** — utility-first styling via PostCSS
- **Lucide React** — icon library
- **Vitest 4** — test runner (configured in `vitest.config.mts`)
- **React Testing Library** — component testing with jsdom environment
- **TypeScript 5** — strict mode enabled

## Architecture

**Next.js App Router** with two route groups that share no layout:

- `app/(public)/` — Unauthenticated pages (home, login, signup, preview). No Navbar.
- `app/(dashboard)/` — Authenticated pages under `/heists`. Includes the Navbar via its own layout.

The home page (`(public)/page.tsx`) is a splash/redirect page only — it does not contain app logic.

**Components** live in `components/<ComponentName>/` with a barrel `index.ts` export. Scoped styles use CSS Modules (`*.module.css`).

**Styling** is Tailwind CSS 4 with a custom theme defined in `app/globals.css` using the `@theme` directive. Theme colours: `primary` (#C27AFF), `secondary` (#FB64B6), `dark` (#030712), `success` (#05DF72), `error` (#FF6467). Global utility classes (`.page-content`, `.center-content`, `.form-title`) are also defined there.

**Tests** live in `tests/` mirroring the source structure (e.g. `tests/components/`). Uses Vitest + React Testing Library with jsdom. Query by semantic roles; avoid implementation-detail queries.

**Path alias** `@/*` resolves to the project root — use it for all imports.

## Additional coding preferences

- Do NOT use semicolons for JavaScript or TypeScript code.
- Do NOT apply tailwind classes directly in component template unless essential or just 1 at most. If an element needs more than a single tailwind class, combine them in custom class using `@apply` directive.
- Use minimal project dependency possible.
- Use the `git switch -c` command to switch to new branch, not `git checkout`.