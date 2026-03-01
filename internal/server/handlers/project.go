package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
		RenderError(w, r, 500, err)
		return
	}

	project, err := c.ps.Get(r.Context(), projectID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if project == nil {
		RenderError(w, r, 404, fmt.Errorf("project not found"))
		return
	}

	if err2 := c.ps.Can(lib.MustGetUser(r.Context()), project); err2 != nil {
		RenderError(w, r, 403, err2)
		return
	}

	topics, err := c.ts.ListDescendants(r.Context(), project.Root)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	Render(w, r, pages.ProjectPage(pages.ProjectPageProps{Project: project, TopicCount: len(topics)}))
}
