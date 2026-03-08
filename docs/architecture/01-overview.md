# Architecture Overview

## Goal

Build server-rendered web applications that stay simple by default and become interactive incrementally.

## Core Stack

- Go for application and HTTP layer.
- Templ for server-rendered views/components.
- HTMX for server-driven interactivity.
- Alpine.js for client-side interactivity that is not a good fit for HTMX.
- PostgreSQL for persistence.
- sqlc for typed data-access code generation.

## Guiding Principles

1. Start from full page delivery.
2. Move only truly interactive sections into dedicated components, routes, and handlers.
3. Keep UI elements reusable and explicit through `Props` structs.
4. Keep business logic in services, not in handlers or templates.
5. Treat validation failures as normal UI responses, not exceptional transport errors.
6. Keep translation integrated in template strings and form validation messages.

## High-Level Request Flow

1. Router resolves a named route.
2. Middleware enriches request context (session, auth, CSRF, i18n, logger).
3. Handler parses request/form and delegates to domain services.
4. Handler renders:
   - Full page for normal navigation.
   - Fragment for HTMX requests.
5. UI is progressively enhanced with Alpine.js where needed.
