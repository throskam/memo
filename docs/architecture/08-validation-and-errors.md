# Validation and Error Rendering

## Validation Model

Validation is represented as form state and messages, not as exceptional runtime errors.

Pattern:

1. Parse request into form object.
2. Validate and collect translated messages.
3. If invalid, return the same component/fragment populated with messages.
4. If valid, continue business action and return next UI state.

For HTMX, invalid form responses commonly return `422` with the updated fragment.

## Error Model

Reserve error rendering for true failures (unexpected exceptions, authorization failures, missing resources).

- Non-HTMX/full navigation: render an error page.
- HTMX fragment request: render a compact fragment notification.

This separates normal validation UX from operational failures.
