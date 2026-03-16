package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/throskam/kix/auth"
	"github.com/throskam/kix/htmx"
	"github.com/throskam/kix/i18n"
	"github.com/throskam/kix/sess"
	"github.com/throskam/memo/internal/lib"
	"github.com/throskam/memo/internal/server/handlers"
)

func Session(
	store sess.SessionStore,
) func(http.Handler) http.Handler {
	handleError := func(w http.ResponseWriter, r *http.Request, err error) {
		handlers.RenderProblem(w, r, handlers.NewProblem(err))
	}

	return sess.Sessionizer(store, handleError)
}

func Authenticate(
	us *lib.UserService,
	jwks auth.JWKS,
	store sess.SessionStore,
) func(http.Handler) http.Handler {
	identify := func(ctx context.Context, session *sess.Session) (any, error) {
		bearer := session.GetFirst("bearer")
		if bearer == "" {
			return nil, nil
		}

		claims, err := auth.ParseJWT(jwks, bearer, time.Minute*5)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JWT: %w", err)
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return nil, fmt.Errorf("missing sub claims")
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			return nil, fmt.Errorf("failed to parse UUID [sub: %v]: %w", sub, err)
		}

		user, err := us.Get(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user [ID: %v]: %w", userID, err)
		}

		if user == nil {
			return nil, fmt.Errorf("user not found [ID: %v]", userID)
		}

		return user, nil
	}

	handleError := func(w http.ResponseWriter, r *http.Request, err error) {
		// e.g expiration
		if errors.Is(err, auth.ErrJWTClaimsParseFailure) {
			err1 := store.Erase(r, w)
			if err1 != nil {
				handlers.RenderProblem(w, r, handlers.NewProblem(err1))
				return
			}

			redirectURL := fmt.Sprintf("/auth?redirect_url=%s", r.URL.String())

			http.Redirect(w, r, redirectURL, http.StatusSeeOther)

			return
		}

		handlers.RenderProblem(w, r, handlers.NewProblem(err))
	}

	return auth.Authenticate(identify, handleError)
}

func Anonymous() func(http.Handler) http.Handler {
	handleError := func(w http.ResponseWriter, r *http.Request, err error) {
		htmx.Redirect(w, r, "/")
	}

	return auth.Anonymous[*lib.User](handleError)
}

func Authenticated() func(http.Handler) http.Handler {
	handleError := func(w http.ResponseWriter, r *http.Request, err error) {
		htmx.Redirect(w, r, fmt.Sprintf("/auth?redirect_url=%s", r.URL.String()))
	}

	return auth.Authenticated[*lib.User](handleError)
}

func CSRF() func(http.Handler) http.Handler {
	handleError := func(w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, auth.ErrCSRFTokenMismatch) {
			handlers.RenderProblem(w, r, handlers.NewProblem(
				err,
				handlers.WithStatus(http.StatusForbidden),
				handlers.WithDetail(i18n.T(r.Context(), "Your request could not be verified. Refresh the page and try again.")),
			))
			return
		}

		handlers.RenderProblem(w, r, handlers.NewProblem(err))
	}

	return auth.CSRF(handleError)
}
