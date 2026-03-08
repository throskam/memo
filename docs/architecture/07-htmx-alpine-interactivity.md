# HTMX and Alpine.js Interaction Model

## Primary Model: HTMX

Use HTMX for interactions where the server should remain the state authority:

- form submit and re-render
- list refresh
- partial updates
- multi-step server workflows

## Secondary Model: Alpine.js

Use Alpine.js for light client-side behavior that does not justify a server roundtrip:

- local toggles and disclosure state
- keyboard shortcuts/focus helpers
- transient client-only UI state

## Rule of Thumb

- If the interaction changes persistent/domain state: HTMX.
- If the interaction is purely local presentation behavior: Alpine.js.

This avoids a heavy front-end state layer while preserving rich UX.
