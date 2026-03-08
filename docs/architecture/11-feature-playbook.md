# Feature Delivery Playbook

Use this workflow for new features in HTMX-first Go projects.

## 1. Model and Data

- Add migration(s) as needed.
- Write SQL query files.
- Regenerate sqlc code.
- Add/extend service methods.

## 2. Page First

- Implement or extend the full page template.
- Add a page handler and route.
- Ensure baseline UX works without fragment splitting.

## 3. Introduce Interactivity

- Identify one interactive section.
- Extract section component.
- Add fragment route + handler.
- Connect with HTMX attributes/events.

## 4. Form and Validation

- Parse into typed form struct.
- Validate with translated messages.
- Return updated fragment/page state.

## 5. Client Enhancements

- Add Alpine.js only for local, client-only behavior.
- Avoid duplicating server state in browser memory.

## 6. Quality Checks

- Verify page load, fragment updates, and empty/error states.
- Verify auth/permission boundaries.
- Verify translations for all user-facing strings.
- Verify accessibility basics (labels, focus flow, ARIA state).
