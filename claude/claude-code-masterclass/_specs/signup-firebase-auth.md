# Spec for signup-firebase-auth

branch: claude/feature/signup-firebase-auth
figma_component (if used):| N/A

## Summary

Hook the existing signup form (`AuthForm` in `app/(public)/signup/`) to Firebase Authentication. On successful signup, generate a random codename for the user by combining words from three distinct word lists in PascalCase, set it as the Firebase `displayName`, and persist it (along with the user's UID) to a Firestore `users` collection document. Email must not be stored.

## Functional Requirements

- When the signup form is submitted, call `createUserWithEmailAndPassword` from the Firebase web SDK using the `auth` export from `lib/firebase/config.ts`
- After account creation, generate a random codename by:
  - Defining three separate word lists (e.g. adjectives, nouns, verbs or any three thematically distinct sets)
  - Randomly picking one word from each list
  - Concatenating them in PascalCase (e.g. `SwiftCrimsonFox`)
- Update the newly created user's Firebase profile by setting `displayName` to the generated codename
- Create a document in the Firestore `users` collection with:
  - `id`: the user's Firebase UID
  - `codename`: the generated codename
  - No email field
- Use only the Firebase web SDK (no Admin SDK, no server-side code)
- Use the existing `auth` export from `lib/firebase/config.ts`; add a `db` (Firestore) export to the same file if not already present
- On successful signup, redirect the user to the authenticated area (`/heists`)
- On failure, display an appropriate error message in the form

## Possible Edge Cases

- Firebase account creation fails (e.g. email already in use, weak password) — surface Firebase error message to the user
- Firestore write fails after account creation — user account exists but has no `users` document; consider retry or silent failure with logging
- Codename collision in Firestore — word lists should be large enough to make collisions negligible; no uniqueness enforcement required for now
- Network offline during signup — show a generic error

## Acceptance Criteria

- Submitting the signup form with a valid email and password creates a new Firebase Auth user
- The created user has a `displayName` set to a PascalCase codename composed of three words
- A document exists in `users/{uid}` (or as an auto-ID doc) containing `id` and `codename` fields, with no `email` field
- The user is redirected to `/heists` after successful signup
- Errors from Firebase are displayed to the user without crashing the app

## Open Questions

- Should the `users` document ID match the Firebase UID, or be auto-generated? use UID as document ID for easy lookups
- Are there thematic constraints on the word lists (e.g. heist/crime theme)? no constraints
- Should the codename be re-generated if a Firestore write fails, or is any valid codename acceptable? any valid is OK

## Testing Guidelines

Create a test file(s) in the `./tests` folder for the new feature, and create meaningful tests for the following cases, without going too heavy:

- The codename generator always returns a string in valid PascalCase composed of exactly three words
- The codename generator returns different values across multiple calls (probabilistic)
- The signup handler calls `createUserWithEmailAndPassword` with the provided credentials
- On success, `updateProfile` is called with a non-empty `displayName`
- On success, a Firestore document is written with `id` and `codename` but no `email`
- On Firebase error, the error message is displayed in the form
