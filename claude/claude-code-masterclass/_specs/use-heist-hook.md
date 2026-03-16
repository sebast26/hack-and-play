# Spec for useHeist Hook

branch: claude/feature/use-heist-hook

## Summary

Create a `useHeist` React hook in `hooks/useHeist.ts` that subscribes to real-time Firestore data from the `heists` collection. The hook accepts a mode argument (`'active'`, `'assigned'`, or `'expired'`) and returns a typed array of `Heist` objects filtered by the appropriate query. Wire the hook into `app/(dashboard)/heists/page.tsx` to display heist titles across the three sections.

## Functional Requirements

- The hook is named `useHeist` and accepts a single argument: `mode: 'active' | 'assigned' | 'expired'`
- The hook returns an array of `Heist` objects (from the existing `Heist` type in `types/firestore`)
- The hook uses Firestore's `onSnapshot` for real-time updates — data refreshes automatically when the underlying collection changes
- The authenticated user's id is sourced from `AuthContext` (`useUser`)
- Query behaviour per mode:
  - **`'active'`** — heists where `assignedTo === currentUser.uid` AND `deadline > now`
  - **`'assigned'`** — heists where `createdBy === currentUser.uid` AND `deadline > now`
  - **`'expired'`** — heists where `deadline <= now` AND `finalStatus !== null` (regardless of user)
- The hook cleans up the Firestore listener (unsubscribes) when the component unmounts or the mode changes
- The hook uses the `heistConverter` already defined in `types/firestore/index.ts` for typed deserialization
- The `heists` page (`app/(dashboard)/heists/page.tsx`) uses three instances of the hook (one per mode) and renders only the title of each heist under the appropriate section heading

## Possible Edge Cases

- User is not yet authenticated when the hook mounts — hook should not attempt a Firestore query until `user` is available
- A mode produces zero results — the section renders an empty state (e.g. "No heists yet")
- Firestore query fails — surface an error state; do not crash the page
- The `mode` argument changes at runtime — the previous listener must be unsubscribed before the new one is created

## Acceptance Criteria

- `useHeist('active')` returns only heists where the current user is the assignee and the deadline has not passed
- `useHeist('assigned')` returns only heists where the current user is the creator and the deadline has not passed
- `useHeist('expired')` returns heists where the deadline has passed and `finalStatus` is not null, regardless of which user created or was assigned to them
- Data updates in real time without a page refresh
- The listener is cleaned up on unmount / mode change
- `app/(dashboard)/heists/page.tsx` renders three sections, each showing heist titles from the corresponding mode
- Each section shows an empty state message when there are no results

## Open Questions

- Should the hook also return a `loading` boolean so the UI can show a skeleton/spinner while the initial snapshot arrives? yes
- Should the hook return an `error` value alongside the heists array? yes

## Testing Guidelines

Create a test file at `tests/hooks/useHeist.test.ts`. Write meaningful tests for the following cases:

- Returns an empty array initially before the snapshot fires
- Returns heists matching the `'active'` query (assignedTo current user, deadline in future)
- Returns heists matching the `'assigned'` query (createdBy current user, deadline in future)
- Returns heists matching the `'expired'` query (deadline in past, finalStatus not null)
- Does not include heists that do not match the active mode's criteria
- Unsubscribes from Firestore when the component unmounts
- Does not query Firestore when the user is not authenticated
