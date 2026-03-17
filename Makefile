ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: db
db:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -database "${DATABASE_CONNECTION_STRING}" -path db/migrations up

.PHONY: translations
translations:
	go generate ./internal/translations/translations.go
	./scripts/prune-translations.sh

.PHONY: orm
orm:
	go tool sqlc generate

.PHONY: templ
templ:
	go tool templ generate

.PHONY: build
build: orm templ
	go build -o .build/bin/server cmd/server/main.go

.PHONY: run
run: build
	go tool godotenv -o .build/bin/server | go tool gojq -R '. as $$line | try (fromjson) catch $$line'

.PHONY: debug
debug: build
	 go tool dlv exec --continue --accept-multiclient --listen=0.0.0.0:2345 --headless=true --api-version=2 --log .build/bin/server

.PHONY: live
live:
	go tool air

.PHONY: lint
lint:
	go tool golangci-lint run ./...