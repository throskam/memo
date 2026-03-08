# Project Structure

A reusable, feature-oriented structure for Go + HTMX apps:

```text
cmd/server/                 # application entrypoint
internal/server/            # HTTP server, middleware, route registration
internal/server/handlers/   # request handlers per feature/page
internal/lib/               # domain services and reusable business logic
internal/views/pages/       # page and page-fragment templates (Templ)
internal/views/ui/          # reusable UI components (Templ)
internal/orm/               # generated sqlc code
internal/translations/      # i18n catalogs and runtime translation setup
db/migrations/              # schema migrations
db/queries/                 # SQL source files for sqlc generation
web/static/                 # CSS, JS, icons, and static assets
docs/                       # architecture and operational documentation
```

## Boundary Rules

- `handlers` should orchestrate HTTP concerns only.
- `lib` should own business/domain behavior.
- `views/pages` should map closely to route/page concerns.
- `views/ui` should be reusable and feature-agnostic.
- `db/queries` stays SQL-first; generated code is never edited manually.
