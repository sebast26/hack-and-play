# Spec for Create Heist Form

branch: claude/feature/create-heist-form

## Summary

Add a form to `app/(dashboard)/heists/create/page.tsx` that allows authenticated users to create a new heist. The form collects the user-facing fields from `CreateHeistInput`. The `createdAt` timestamp and `deadline` are added programmatically (not entered by the user). On successful submission, a new document is written to the `heists` Firestore collection and the user is redirected to `/heists`.

## Functional Requirements

- The form must include the following user-facing fields derived from `CreateHeistInput`:
  - **Title** — short text input, required
  - **Description** — multiline text input, required
  - **Assign To** — dropdown populated from the `users` Firestore collection; displays each user's codename, stores their user id. The selected value provides both `assignedTo` (user id) and `assignedToCodename`
- The following fields are set programmatically and are not shown as form inputs:
  - `createdBy` — derived from the currently authenticated user's id (via `AuthContext`)
  - `createdByCodename` — derived from the currently authenticated user's codename (via `AuthContext`)
  - `createdAt` — set using Firestore `serverTimestamp()` at submission time
  - `deadline` — calculated programmatically at submission time (exact rule TBD — see Open Questions)
  - `finalStatus` — always initialised to `null`
- On valid submission, write a document to the `heists` collection shaped exactly as `CreateHeistInput`
- On success, redirect the user to `/heists`
- While submission is in flight, disable the submit button to prevent double-submission
- Display an inline error message if the Firestore write fails

## Possible Edge Cases

- The `users` collection is empty — dropdown shows a "no users available" state and submission is blocked
- Authenticated user assigns the heist to themselves — allowed
- Firestore write fails (network error, permissions) — show inline error, keep form data intact so the user can retry
- User navigates away mid-form — no autosave required

## Acceptance Criteria

- Form renders at `/heists/create` with title, description, and assign-to fields
- Assignee dropdown is populated from the `users` Firestore collection (codename shown, user id stored)
- Submitting with all valid fields creates a new document in the `heists` collection matching `CreateHeistInput`
- `createdAt`, `createdBy`, `createdByCodename`, `deadline`, and `finalStatus` are set programmatically — not editable by the user
- After a successful write, the user is redirected to `/heists`
- Submit button is disabled while submission is in progress
- An error message is shown if the write fails; the form remains editable

## Open Questions

- What fields does a user document in the `users` collection expose? At minimum `id` and `codename` are needed for the assignee dropdown. id and codename
- How should `deadline` be calculated programmatically? (e.g. fixed offset from `createdAt`, a default period, or a constant date?). fixed offset from createdAt
- Should there be a `USERS` constant added to `COLLECTIONS` in `types/firestore/index.ts`? if required.

## Testing Guidelines

Create a test file at `tests/app/heists/create/page.test.tsx`. Write meaningful tests for the following cases:

- Renders all user-facing form fields (title, description, assign-to dropdown)
- Assignee dropdown is populated with users fetched from Firestore
- Submitting the form with valid data calls the Firestore `addDoc` function with a payload matching `CreateHeistInput`, including programmatically set fields
- After successful submission, the router navigates to `/heists`
- Submit button is disabled while submission is in progress
- An error message is displayed when the Firestore write rejects
