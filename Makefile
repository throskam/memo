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

.PHONY: css
css:
	podman run --rm \
		-v $(PWD)/web/static/css:/css:rw \
		-w /home/node node:24-alpine \
		sh -c 'echo "module.exports = { plugins: { '\''postcss-import'\'': {}, '\''postcss-nesting'\'': {}, '\''postcss-flexbugs-fixes'\'': {}, '\''autoprefixer'\'': {} } };" > postcss.config.js && \
		cp /css/* . 2>/dev/null; \
		npm install postcss postcss-cli postcss-import postcss-nesting postcss-flexbugs-fixes autoprefixer && \
		./node_modules/.bin/postcss main.css -o main.compat.css && \
		cp main.compat.css /css/'

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