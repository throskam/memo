package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/throskam/ki"
	"github.com/throskam/kix/i18n"
	"github.com/throskam/memo/internal/views/ui"
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	err := c.Render(r.Context(), w)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
	}
}

func RenderProblem(w http.ResponseWriter, r *http.Request, problem Problem) {
	ki.MustGetLogger(r.Context()).LogAttrs(
		r.Context(),
		slog.LevelError,
		"error",
		slog.Any("err", problem.Err),
	)

	if problem.Instance == "" {
		problem.Instance = r.URL.Path
	}

	if problem.Title == "" {
		problem.Title = getDefaultProblemTitle(r.Context(), problem)
	}

	// Page level errors.

	if r.Header.Get("HX-Boosted") == "true" || r.Header.Get("HX-Request") != "true" {
		page := ui.ErrorPage(ui.ErrorPageProps{
			Code:  problem.Status,
			Title: problem.Title,
			Text:  problem.Detail,
		})

		w.WriteHeader(problem.Status)

		err2 := page.Render(r.Context(), w)
		if err2 != nil {
			panic(err2)
		}

		return
	}

	// Fragment level errors.

	title := problem.Title
	text := problem.Detail

	if text == "" {
		title = fmt.Sprintf("%d", problem.Status)
		text = getDefaultProblemTitle(r.Context(), problem)
	}

	notification := ui.Notification(ui.NotificationProps{
		Kind:  "alert",
		Title: title,
		Text:  text,
	})

	w.Header().Set("HX-Reswap", "none")
	w.WriteHeader(problem.Status)

	err2 := notification.Render(r.Context(), w)
	if err2 != nil {
		panic(err2)
	}
}

func getDefaultProblemTitle(ctx context.Context, problem Problem) string {
	switch problem.Type {
	case TypeForbidden:
		return i18n.T(ctx, "Forbidden")
	case TypeNotFound:
		return i18n.T(ctx, "Not Found")
	default:
		return i18n.T(ctx, "Internal Server Error")
	}
}
