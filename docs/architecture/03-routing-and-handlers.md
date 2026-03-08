# Routing and Handlers

## Router Strategy

Use a router abstraction that supports named routes. Named routes reduce hard-coded URL strings in templates and handlers and make refactors safer.

This project uses:

- `ki`: wrapper around Go's native router with named-route ergonomics.
- `kix`: helper packages for common HTMX/web concerns.

## Route Organization

- Group by feature/page (`/home`, `/topic`, `/profile`, etc.).
- Split page route and fragment routes:
  - Page route: returns full page.
  - Fragment routes: return focused partials for HTMX swaps.

## Handler Responsibilities

1. Parse and validate request/form data.
2. Authorize the action.
3. Delegate to domain services.
4. Render a full page or fragment.

Avoid putting business rules in handlers. If logic is not HTTP-specific, move it to services.

## Naming

Adopt stable route names with a predictable shape, for example:

```text
feature:section:action
home:project-create:submit
topic:descendant-list:get
```

This keeps location construction explicit in templates and handlers.
