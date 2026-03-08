# Data Access with sqlc

## SQL-First Model

Write SQL queries first, then generate typed store methods.

1. Add/update SQL in `db/queries/*.sql`.
2. Run sqlc generation.
3. Use generated methods from service layer.

## Why sqlc

- Keeps queries explicit and reviewable.
- Produces compile-time-checked method signatures.
- Avoids runtime query-string assembly for core operations.

## Recommended Layering

- `db/queries`: source of truth for data access behavior.
- `internal/orm`: generated data-access package.
- `internal/lib`: domain services that orchestrate generated methods.

Never edit generated files directly.
