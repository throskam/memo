# Documentation

This folder contains reusable architecture guidance for Go + HTMX + Templ applications.

## Contents

- [Architecture Overview](./architecture/01-overview.md)
- [Project Structure](./architecture/02-project-structure.md)
- [Routing and Handlers](./architecture/03-routing-and-handlers.md)
- [Page-First Component Strategy](./architecture/04-page-component-strategy.md)
- [UI and `Props` Conventions](./architecture/05-ui-props-conventions.md)
- [Data Access with sqlc](./architecture/06-data-access-sqlc.md)
- [HTMX and Alpine.js Interaction Model](./architecture/07-htmx-alpine-interactivity.md)
- [Validation and Error Rendering](./architecture/08-validation-and-errors.md)
- [Internationalization Workflow](./architecture/09-i18n-workflow.md)
- [Passwordless Authentication Architecture](./architecture/10-authentication-passwordless.md)
- [Feature Delivery Playbook](./architecture/11-feature-playbook.md)

## Intended Usage

Use these docs as a baseline for HTMX-first Go apps.

- Keep the principles and flow.
- Adjust naming and package boundaries per project.
- Extend each document with domain-specific decisions (billing, ACL, audit, etc.) when needed.
