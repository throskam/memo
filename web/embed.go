// Package web provides static assets for the web application.
package web

import (
	"embed"
	"io/fs"
)

//go:embed all:static
var static embed.FS

func Static() fs.FS {
	f, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

	return f
}