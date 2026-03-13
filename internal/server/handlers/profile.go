package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/throskam/kix/htmx"
	"github.com/throskam/memo/internal/lib"
	"github.com/throskam/memo/internal/views/pages"
	"github.com/throskam/memo/internal/views/ui"
)

type ProfileController struct {
	us *lib.UserService
}

func NewProfileController(us *lib.UserService) *ProfileController {
	return &ProfileController{
		us: us,
	}
}

func (c *ProfileController) PageGet(w http.ResponseWriter, r *http.Request) {
	Render(w, r, pages.ProfilePage(pages.ProfilePageProps{User: lib.MustGetUser(r.Context())}))
}

func (c *ProfileController) InfoGet(w http.ResponseWriter, r *http.Request) {
	user := lib.MustGetUser(r.Context())

	Render(w, r, pages.ProfileInfo(pages.ProfileInfoProps{User: user}))
}

func (c *ProfileController) InfoEdit(w http.ResponseWriter, r *http.Request) {
	user := lib.MustGetUser(r.Context())

	form := htmx.NewForm(&pages.ProfileInfoEditForm{
		Username: user.Username,
	})

	Render(w, r, pages.ProfileInfoEdit(pages.ProfileInfoEditProps{User: user, Form: form}))
}

func (c *ProfileController) InfoSave(w http.ResponseWriter, r *http.Request) {
	user := lib.MustGetUser(r.Context())

	form := htmx.NewFormFromRequest(r, &pages.ProfileInfoEditForm{})

	if !form.OK() {
		w.WriteHeader(422)

		Render(w, r, pages.ProfileInfoEdit(pages.ProfileInfoEditProps{User: user, Form: form}))

		return
	}

	user.Username = form.Data.Username

	user, err := c.us.Update(r.Context(), user)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	Render(w, r, pages.ProfileInfo(pages.ProfileInfoProps{User: user}))
	Render(w, r, ui.OOB(
		"innerHTML:#navbar-username",
		templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, user.Username)
			return err
		}),
	))
}