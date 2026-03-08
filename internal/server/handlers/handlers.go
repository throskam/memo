package handlers

import (
	"fmt"
	"net/http"

	"github.com/throskam/kix/sess"
)

func NewNotFoundHanlder() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("route not found: %s", r.RequestURI),
			WithStatus(http.StatusNotFound),
		))
	})
}

func NewLogoutHandler(store sess.SessionStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := store.Erase(r, w)
		if err != nil {
			RenderProblem(w, r, NewProblem(err))
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func NewRedirectHandler(redirectPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
	}
}
