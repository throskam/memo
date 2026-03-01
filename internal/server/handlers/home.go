package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/throskam/kix/htmx"
	"github.com/throskam/memo/internal/lib"
	"github.com/throskam/memo/internal/views/pages"
)

type HomeController struct {
	ps *lib.ProjectService
	ts *lib.TopicService
}

func NewHomeController(ps *lib.ProjectService, ts *lib.TopicService) *HomeController {
	c := &HomeController{
		ps: ps,
		ts: ts,
	}

	return c
}

func (c *HomeController) PageGet(w http.ResponseWriter, r *http.Request) {
	projects, err := c.ps.ListByOwnerWithRoot(r.Context(), lib.MustGetUser(r.Context()))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	Render(w, r, pages.HomePage(pages.HomePageProps{Projects: projects}))
}

func (c *HomeController) ProjectListGet(w http.ResponseWriter, r *http.Request) {
	projects, err := c.ps.ListByOwnerWithRoot(r.Context(), lib.MustGetUser(r.Context()))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	Render(w, r, pages.HomeProjectList(pages.HomeProjectListProps{Projects: projects}))
}

func (c *HomeController) ProjectCreateSubmit(w http.ResponseWriter, r *http.Request) {
	form := htmx.NewFormFromRequest(r, &pages.HomeProjectCreateForm{})

	if !form.OK() {
		w.WriteHeader(422)

		Render(w, r, pages.HomeProjectCreate(pages.HomeProjectCreateProps{Form: form}))

		return
	}

	user := lib.MustGetUser(r.Context())

	project := &lib.Project{
		OwnerID: user.ID,

		Owner: user,
	}

	project, err := c.ps.Create(r.Context(), project)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	root := &lib.Topic{
		Title:     form.Data.Name,
		SortOrder: 0,

		ProjectID: project.ID,

		Project: project,
	}

	_, err = c.ts.Create(r.Context(), root)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	w.Header().Set("HX-Trigger", "project-created")

	Render(w, r, pages.HomeProjectCreate(pages.HomeProjectCreateProps{Form: htmx.NewForm(&pages.HomeProjectCreateForm{})}))
}

func (c *HomeController) ProjectItemDelete(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(r.FormValue("project"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	project, err := c.ps.Get(r.Context(), projectID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ps.Can(lib.MustGetUser(r.Context()), project); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	err = c.ps.Remove(r.Context(), project)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	w.WriteHeader(204)
}
