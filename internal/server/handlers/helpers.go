package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/throskam/ki"
	"github.com/throskam/memo/internal/views/ui"
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	err := c.Render(r.Context(), w)
	if err != nil {
		RenderError(w, r, 500, err)
	}
}

func RenderError(w http.ResponseWriter, r *http.Request, status int, err error) {
	ki.MustGetLogger(r.Context()).LogAttrs(
		r.Context(),
		slog.LevelError,
		"error",
		slog.Any("err", err),
	)

	// Page level errors.
	if r.Header.Get("HX-Boosted") == "true" || r.Header.Get("HX-Request") != "true" {
		page := ui.ErrorPage(ui.ErrorPageProps{
			Code:  status,
			Title: http.StatusText(status),
		})

		w.WriteHeader(status)

		err2 := page.Render(r.Context(), w)
		if err2 != nil {
			panic(err2)
		}

		return
	}

	// Fragment level errors.
	notification := ui.Notification(ui.NotificationProps{
		Kind:  "alert",
		Title: fmt.Sprintf("%d", status),
		Text:  http.StatusText(status),
	})

	w.Header().Set("HX-Reswap", "none")
	w.WriteHeader(status)

	err2 := notification.Render(r.Context(), w)
	if err2 != nil {
		panic(err2)
	}
}
