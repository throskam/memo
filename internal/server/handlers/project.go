package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/throskam/kix/i18n"
	"github.com/throskam/memo/internal/lib"
	"github.com/throskam/memo/internal/views/pages"
)

type ProjectController struct {
	ts *lib.TopicService
	ps *lib.ProjectService
}

func NewProjectController(ps *lib.ProjectService, ts *lib.TopicService) *ProjectController {
	return &ProjectController{
		ts: ts,
		ps: ps,
	}
}

func (c *ProjectController) PageGet(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(r.PathValue("project"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	project, err := c.ps.Get(r.Context(), projectID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if project == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("project not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested project could not be found.")),
		))
		return
	}

	if err2 := c.ps.Can(lib.MustGetUser(r.Context()), project); err2 != nil {
		RenderProblem(w, r, NewProblem(
			err2,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to access this project.")),
		))
		return
	}

	topics, err := c.ts.ListDescendants(r.Context(), project.Root)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	Render(w, r, pages.ProjectPage(pages.ProjectPageProps{Project: project, TopicCount: len(topics)}))
}
