# Spec for Login Form Functionality

branch: claude/feature/login-form-functionality

## Summary

Wire up the login form in `app/(public)/login/` so that users can authenticate with their email and password via Firebase Auth. On success, display an inline success message (no redirect). On failure, display a relevant error message.

## Functional Requirements

- The `AuthForm` component in `login` mode should call the Firebase Auth sign-in method when the form is submitted
- On successful login, show an inline success message (e.g. "You're logged in!")
- On failed login (wrong credentials, user not found, etc.), display a human-readable error message below the form fields
- While the request is in flight, the submit button should be disabled and show a loading state
- The form should clear any previous error/success state when the user resubmits

## Possible Edge Cases

- Invalid email format (already handled by `type="email"` input, but Firebase may also return an error)
- Wrong password — Firebase returns a specific error code
- User does not exist — Firebase returns a specific error code
- Network failure or unexpected Firebase error
- User submits the form multiple times quickly (debounce / disable submit during loading)

## Acceptance Criteria

- Submitting valid credentials logs the user in and shows a success message
- Submitting invalid credentials shows a readable error message
- The submit button is disabled while the login request is pending
- No page redirect occurs after successful login
- Success and error states are mutually exclusive — only one is shown at a time
- An existing `loginUser` Firebase helper is used (or created in `lib/firebase/`) consistent with the existing `signUpUser` pattern

## Open Questions

- Should the success message auto-dismiss after a timeout, or stay until the user navigates away? auto-dismiss
- Should the logged-in user's codename be shown in the success message (e.g. "Welcome back, Shadow Fox!")? no.

## Testing Guidelines

Create a test file in `tests/components/AuthForm.login.test.tsx` for the login mode, covering:

- Renders email, password inputs and a submit button
- Shows a success message after submitting valid credentials (mock Firebase)
- Shows an error message when Firebase rejects the credentials
- Submit button is disabled while loading
- Error is cleared and replaced when the user resubmits
