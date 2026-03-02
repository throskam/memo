FROM golang:1.24-bullseye AS base

WORKDIR /opt/app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

FROM base AS dev

ENV APP_ENV=development

EXPOSE 8080

FROM base AS build

RUN useradd -u 1000 go

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o server \
  cmd/server/main.go

FROM scratch

ENV APP_ENV=production

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /opt/app/server server

USER go

EXPOSE 8080

CMD ["/server"]