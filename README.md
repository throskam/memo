# memo

A small Go web app for organizing projects and nested topics (notes), with passwordless auth and PostgreSQL storage.

## Stack

- Go `1.24.6`
- PostgreSQL `16`
- `templ` for server-side views
- `sqlc` + `pgx` for database access

## Quick Start (Docker)

1. Copy env file:
   - `cp .env.sample .env`
1. Fill required values in `.env`:
   - `APP_BASE_URL` (for example `http://localhost:1080`)
   - `COOKIE_HASH_KEY`, `COOKIE_BLOCK_KEY`, `JWT_SECRET` (base64-encoded keys)
   - `DATABASE_CONNECTION_STRING` (for example `postgres://user:password@db:5432/app?sslmode=disable`)
   - `MAIL_SMTP` (for example `localhost:1025`)
1. Start services:
   - `docker compose up --build`
1. Run DB migrations:
   - `make db`
1. App is available at `http://localhost:1080`.

Auxiliary services:

- MailHog UI: `http://localhost:8025`
- pgAdmin: `http://localhost:8888`

## Common Commands

- `make live` - run with live reload (`air`)
- `make build` - generate code + build server binary
- `make run` - build and run server
- `make db` - run migrations
- `make orm` - regenerate `sqlc` code
- `make templ` - regenerate templ views
- `make lint` - run linter