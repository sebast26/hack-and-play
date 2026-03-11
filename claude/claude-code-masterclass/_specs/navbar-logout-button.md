# Spec for Navbar Logout Button

branch: claude/feature/navbar-logout-button

## Summary

Add a logout button to the Navbar component that signs the currently authenticated user out of the application when clicked. The button is only visible when a user is logged in.

## Functional Requirements

- The Navbar must display a logout button only when the user is authenticated
- Clicking the logout button signs the user out via Firebase Auth
- No redirect is required after logout — the UI simply reflects the signed-out state
- The button label reads "Logout" and includes an arrow-right (exit) icon to the left of the text

## Figma Design Reference (only if referenced)

- Component name: Logout Button (from user-provided image)
- Key visual constraints:
  - Gradient background: left-to-right pink-to-purple (matches `primary`/`secondary` theme colours)
  - Rounded pill shape
  - White text and icon
  - Icon: exit/arrow-right bracket icon (e.g. `LogOut` from Lucide React) placed to the left of the label
  - Consistent padding with other Navbar actions

## Possible Edge Cases

- User is not logged in — button must not render at all
- Firebase `signOut` call fails — handle the error gracefully (e.g. log to console), do not crash the UI
- `useUser` hook returns a loading/undefined state — treat as unauthenticated, hide the button

## Acceptance Criteria

- The logout button renders in the Navbar when `useUser` returns an authenticated user
- The logout button is absent from the Navbar when no user is authenticated
- Clicking the button calls Firebase Auth `signOut`
- After sign-out, `useUser` reflects the signed-out state and the button disappears
- Button styling matches the gradient pill design shown in the reference image

## Open Questions

- Should the button be positioned at the far right of the Navbar, or adjacent to other nav items? just left to the create button
- Is a loading/disabled state needed while the sign-out request is in flight? no.

## Testing Guidelines

Create a test file in `tests/components/Navbar.test.tsx` (or extend the existing one) covering:

- Renders the logout button when a user is authenticated
- Does not render the logout button when no user is authenticated
- Calls the Firebase `signOut` function when the button is clicked
- Does not throw when `signOut` rejects (error handling)
