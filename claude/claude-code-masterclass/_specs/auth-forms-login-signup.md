# Spec for Authentication Forms for Login and Signup

branch: claude/feature/auth-forms-login-signup
figma_component (if used):| N/A

## Summary

Add functional authentication forms to the existing `/login` and `/signup` pages. Each page presents a form with email and password fields, a toggle to show/hide the password, and a submit button. On submission, the form data is logged to the console. Users can easily navigate between the two forms via a link.

## Functional Requirements

- The `/login` page renders a login form with an email field, a password field, and a "Login" submit button
- The `/signup` page renders a signup form with an email field, a password field, and a "Sign Up" submit button
- Both forms include a password visibility toggle icon (show/hide) next to the password field
- On form submission, the email and password values are logged to the browser console (no real auth)
- The login page includes a link to the signup page ("Don't have an account? Sign up")
- The signup page includes a link to the login page ("Already have an account? Log in")
- Both forms prevent default browser submission and handle state via React

## Figma Design Reference (only if referenced)

N/A

## Possible Edge Cases

- User submits the form with empty fields — the form should still log (or optionally rely on browser native required validation)
- User rapidly toggles password visibility — state should remain in sync with the input type
- User navigates between login and signup — form state should reset on each page

## Acceptance Criteria

- `/login` renders an email input, a password input with a show/hide toggle, and a "Login" button
- `/signup` renders an email input, a password input with a show/hide toggle, and a "Sign Up" button
- Clicking the password toggle switches the input between `type="password"` and `type="text"`
- Submitting either form logs `{ email, password }` to the console
- Each page contains a visible link that navigates to the other form page
- No third-party auth libraries are introduced

## Open Questions

- Should the password field enforce any minimum length or complexity at the UI level, or is plain logging sufficient for now? No minimum length.
- Should the email field validate format client-side, or is that out of scope for this spec? Light validation.

## Testing Guidelines

Create a test file(s) in the ./tests folder for the new feature, and create meaningful tests for the following cases, without going too heavy:

- Renders the email and password fields on both pages
- Password toggle changes the input type between "password" and "text"
- Submitting the form calls `console.log` with the email and password values
- The navigation link to the other form page is present and points to the correct route
