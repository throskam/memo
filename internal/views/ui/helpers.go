package ui

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/throskam/ki"
)

func HxAction(location ki.Location) templ.Attributes {
	attrs := templ.Attributes{}

	attrs["hx-"+strings.ToLower(location.Method())] = location.URL().String()
	attrs["aria-live"] = "polite"

	return attrs
}