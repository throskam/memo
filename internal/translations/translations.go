// Package translations provides translations for the application.
package translations

//go:generate go tool gotext -srclang=en-US update -out=catalog.go -lang=en-US,fr-FR github.com/throskam/memo/internal/views/pages github.com/throskam/memo/internal/server/handlers
