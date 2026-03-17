package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/throskam/kix/htmx"
	"github.com/throskam/kix/i18n"
	"github.com/throskam/kix/sess"
	"github.com/throskam/memo/internal/lib"
	"github.com/throskam/memo/internal/views/pages"
	"github.com/throskam/memo/internal/views/ui"
)

type TopicController struct {
	ts *lib.TopicService
	ps *lib.ProjectService
}

func NewTopicController(ts *lib.TopicService, ps *lib.ProjectService) *TopicController {
	return &TopicController{
		ts: ts,
		ps: ps,
	}
}

func (c *TopicController) PageGet(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.PathValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		err = fmt.Errorf("topic not found")
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to access this topic.")),
		))
		return
	}

	ancestors, err := c.ts.ListAncestors(r.Context(), topic)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	descendants, err := c.ts.ListDescendants(r.Context(), topic)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	Render(w, r, pages.TopicPage(pages.TopicPageProps{Topic: topic, Ancestors: ancestors, Descendants: descendants}))
}

func (c *TopicController) OverviewGet(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to access this topic.")),
		))
		return
	}

	Render(w, r, pages.TopicOverview(pages.TopicOverviewProps{Topic: topic}))
}

func (c *TopicController) OverviewEdit(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to edit this topic.")),
		))
		return
	}

	sess.MustGetSession(r.Context()).Remove("collapsed-overview", topic.ID.String())

	Render(w, r, pages.TopicOverviewEdit(pages.TopicOverviewEditProps{
		Topic: topic,
		Form:  htmx.NewForm(&pages.TopicOverviewEditForm{Title: topic.Title, Content: topic.Content}),
	}))
}

func (c *TopicController) OverviewExpand(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	sess.MustGetSession(r.Context()).Remove("collapsed-overview", topic.ID.String())

	Render(w, r, pages.TopicOverview(pages.TopicOverviewProps{
		Topic: topic,
	}))
}

func (c *TopicController) OverviewCollapse(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	sess.MustGetSession(r.Context()).Add("collapsed-overview", topic.ID.String())

	Render(w, r, pages.TopicOverview(pages.TopicOverviewProps{
		Topic: topic,
	}))
}

func (c *TopicController) OverviewSave(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to edit this topic.")),
		))
		return
	}

	form := htmx.NewFormFromRequest(r, &pages.TopicOverviewEditForm{})

	if !form.OK() {
		w.WriteHeader(422)

		Render(w, r, pages.TopicOverviewEdit(pages.TopicOverviewEditProps{Topic: topic, Form: form}))

		return
	}

	topic.Title = form.Data.Title
	topic.Content = form.Data.Content

	topic, err = c.ts.Update(r.Context(), topic)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	ancestors, err := c.ts.ListAncestors(r.Context(), topic)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	Render(w, r, pages.TopicOverview(pages.TopicOverviewProps{Topic: topic}))

	Render(w, r, ui.OOB(
		"outerHTML:#breadcrumb",
		pages.TopicBreadcrumb(pages.TopicBreadcrumbProps{Topic: topic, Ancestors: ancestors}),
	))
}

func (c *TopicController) ToolbarCollapseRecursive(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	descendants, err := c.ts.ListDescendants(r.Context(), topic)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	session := sess.MustGetSession(r.Context())

	for _, descendant := range descendants {
		if descendant.ParentID.Valid {
			session.Remove("expanded-topics", descendant.ParentID.UUID.String())
		}
	}

	w.Header().Set("HX-Trigger", "topic-list-updated")
	w.WriteHeader(204)
}

func (c *TopicController) ToolbarExpandRecursive(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	descendants, err := c.ts.ListDescendants(r.Context(), topic)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	session := sess.MustGetSession(r.Context())

	for _, descendant := range descendants {
		if descendant.ParentID.Valid {
			session.Add("expanded-topics", descendant.ParentID.UUID.String())
		}
	}

	w.Header().Set("HX-Trigger", "topic-list-updated")
	w.WriteHeader(204)
}

func (c *TopicController) ToolbarEnableSelectionMode(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	sess.MustGetSession(r.Context()).Add("mode", "selection")

	w.Header().Set("HX-Trigger", "topic-list-updated")
	Render(w, r, pages.TopicToolbar(pages.TopicToolbarProps{
		Topic: topic,
	}))
}

func (c *TopicController) ToolbarDisableSelectionMode(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	sess.MustGetSession(r.Context()).Remove("mode", "selection")

	w.Header().Set("HX-Trigger", "topic-list-updated")
	Render(w, r, pages.TopicToolbar(pages.TopicToolbarProps{
		Topic: topic,
	}))
}

func (c *TopicController) DescendantList(w http.ResponseWriter, r *http.Request) {
	parentID, err := uuid.Parse(r.FormValue("parent"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	parent, err := c.ts.Get(r.Context(), parentID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if parent == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested parent topic could not be found.")),
		))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), parent); err != nil {
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to access this topic.")),
		))
		return
	}

	descendants, err := c.ts.ListDescendants(r.Context(), parent)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	Render(w, r, pages.TopicDescendantList(pages.TopicDescendantListProps{
		Topic:       parent,
		Descendants: descendants,
		Level:       0,
	}))
}

func (c *TopicController) DescendantListMove(w http.ResponseWriter, r *http.Request) {
	zone := r.FormValue("zone")
	isAbove := zone == "above"
	isBelow := zone == "below"

	destinationID, err := uuid.Parse(r.FormValue("destination"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	destination, err := c.ts.Get(r.Context(), destinationID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if destination == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The destination topic could not be found.")),
		))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), destination); err != nil {
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to move topics in this area.")),
		))
		return
	}

	parent := destination
	if isAbove || isBelow {
		parent, err = c.ts.Get(r.Context(), destination.ParentID.UUID)
		if err != nil {
			RenderProblem(w, r, NewProblem(err))
			return
		}

		if parent == nil {
			RenderProblem(w, r, NewProblem(
				fmt.Errorf("topic not found"),
				WithStatus(http.StatusNotFound),
				WithDetail(i18n.T(r.Context(), "The destination parent topic could not be found.")),
			))
			return
		}
	}

	ancestors, err := c.ts.ListAncestors(r.Context(), parent)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	sourceIDs := r.Form["sources"]
	sources := []*lib.Topic{}
	for _, sourceID := range sourceIDs {
		if destinationID.String() == sourceID {
			continue
		}

		for _, ancestor := range ancestors {
			if ancestor.ID.String() == sourceID {
				RenderProblem(w, r, NewProblem(
					fmt.Errorf("cannot create a cycle in the tree"),
					WithType("topic-cycle"),
					WithDetail(i18n.T(r.Context(), "Cannot move topics because it would create a cycle.")),
				),
				)
				return
			}
		}

		ID, err2 := uuid.Parse(sourceID)
		if err2 != nil {
			RenderProblem(w, r, NewProblem(err2))
			return
		}

		source, err2 := c.ts.Get(r.Context(), ID)
		if err2 != nil {
			RenderProblem(w, r, NewProblem(err2))
			return
		}

		if source == nil {
			RenderProblem(w, r, NewProblem(
				fmt.Errorf("topic not found"),
				WithStatus(http.StatusNotFound),
				WithDetail(i18n.T(r.Context(), "One or more selected topics could not be found.")),
			))
			return
		}

		if err = c.ts.Can(lib.MustGetUser(r.Context()), source); err != nil {
			RenderProblem(w, r, NewProblem(
				err,
				WithStatus(http.StatusForbidden),
				WithDetail(i18n.T(r.Context(), "You do not have permission to move one or more selected topics.")),
			))
			return
		}

		sources = append(sources, source)
	}

	sourceCount := len(sources)

	if isAbove || isBelow {
		start := destination.SortOrder

		if isBelow {
			start = start + 1
		}

		err2 := c.ts.Shift(r.Context(), parent, start, sourceCount)
		if err2 != nil {
			RenderProblem(w, r, NewProblem(err2))
			return
		}

		for index, source := range sources {
			err3 := c.ts.Move(r.Context(), source, parent, start+index)
			if err3 != nil {
				RenderProblem(w, r, NewProblem(err3))
				return
			}
		}
	} else {
		children, err2 := c.ts.ListChildren(r.Context(), parent)
		if err2 != nil {
			RenderProblem(w, r, NewProblem(err2))
			return
		}

		childrenCount := len(children)

		for index, source := range sources {
			err3 := c.ts.Move(r.Context(), source, parent, index+childrenCount+1)
			if err3 != nil {
				RenderProblem(w, r, NewProblem(err3))
				return
			}
		}

		sess.MustGetSession(r.Context()).Add("expanded-topics", parent.ID.String())
	}

	// Self healing.
	project, err := c.ps.Get(r.Context(), parent.ProjectID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if project == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("project not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The project for this topic could not be found.")),
		))
		return
	}

	err = c.ts.Reindex(r.Context(), project)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	w.Header().Set("HX-Trigger", "topic-list-updated")
	w.WriteHeader(204)
}

func (c *TopicController) DescendantDelete(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if topic == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested topic could not be found.")),
		))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to delete this topic.")),
		))
		return
	}

	err = c.ts.Remove(r.Context(), topic)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	w.Header().Set("HX-Trigger", "topic-list-updated")
	w.WriteHeader(204)
}

func (c *TopicController) DescendantCollapse(w http.ResponseWriter, r *http.Request) {
	sess.MustGetSession(r.Context()).Remove("expanded-topics", r.FormValue("topic"))

	w.Header().Set("HX-Trigger", "topic-list-updated")
	w.WriteHeader(204)
}

func (c *TopicController) DescendantExpand(w http.ResponseWriter, r *http.Request) {
	sess.MustGetSession(r.Context()).Add("expanded-topics", r.FormValue("topic"))

	w.Header().Set("HX-Trigger", "topic-list-updated")
	w.WriteHeader(204)
}

func (c *TopicController) DescendantCreateSubmit(w http.ResponseWriter, r *http.Request) {
	form := htmx.NewFormFromRequest(r, &pages.TopicDescendantCreateForm{})

	if !form.OK() {
		w.WriteHeader(422)
		Render(w, r, pages.TopicDescendantCreate(pages.TopicDescendantCreateProps{
			Form: form,
		}))

		return
	}

	parent, err := c.ts.Get(r.Context(), form.Data.ParentID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if parent == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("topic not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The requested parent topic could not be found.")),
		))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), parent); err != nil {
		RenderProblem(w, r, NewProblem(
			err,
			WithStatus(http.StatusForbidden),
			WithDetail(i18n.T(r.Context(), "You do not have permission to create topics here.")),
		))
		return
	}

	project, err := c.ps.Get(r.Context(), parent.ProjectID)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	if project == nil {
		RenderProblem(w, r, NewProblem(
			fmt.Errorf("project not found"),
			WithStatus(http.StatusNotFound),
			WithDetail(i18n.T(r.Context(), "The project for this topic could not be found.")),
		))
		return
	}

	topic := &lib.Topic{
		Title: form.Data.Title,

		ParentID:  uuid.NullUUID{UUID: form.Data.ParentID, Valid: true},
		ProjectID: project.ID,

		Project: project,
	}

	_, err = c.ts.Create(r.Context(), topic)
	if err != nil {
		RenderProblem(w, r, NewProblem(err))
		return
	}

	w.Header().Set("HX-Trigger", "topic-list-updated")

	Render(w, r, pages.TopicDescendantCreate(pages.TopicDescendantCreateProps{
		Form: htmx.NewForm(&pages.TopicDescendantCreateForm{
			ParentID: form.Data.ParentID,
		}),
	}))
}
