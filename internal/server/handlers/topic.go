package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/throskam/kix/htmx"
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
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if topic == nil {
		RenderError(w, r, 404, fmt.Errorf("topic not found"))
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	ancestors, err := c.ts.ListAncestors(r.Context(), topic)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	descendants, err := c.ts.ListDescendants(r.Context(), topic)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	Render(w, r, pages.TopicPage(pages.TopicPageProps{Topic: topic, Ancestors: ancestors, Descendants: descendants}))
}

func (c *TopicController) InfoSave(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	form := htmx.NewFormFromRequest(r, &pages.TopicInfoEditForm{})

	if !form.OK() {
		w.WriteHeader(422)

		Render(w, r, pages.TopicInfoEdit(pages.TopicInfoEditProps{Topic: topic, Form: form}))

		return
	}

	topic.Title = form.Data.Title

	topic, err = c.ts.Update(r.Context(), topic)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	ancestors, err := c.ts.ListAncestors(r.Context(), topic)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	Render(w, r, pages.TopicInfo(pages.TopicInfoProps{Topic: topic}))

	Render(w, r, ui.OOB(
		"outerHTML:#breadcrumb",
		pages.TopicBreadcrumb(pages.TopicBreadcrumbProps{Topic: topic, Ancestors: ancestors}),
	))
}

func (c *TopicController) ContentSave(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	form := htmx.NewFormFromRequest(r, &pages.TopicContentForm{})

	if !form.OK() {
		w.WriteHeader(422)

		Render(w, r, pages.TopicContentData(pages.TopicContentDataProps{
			Topic: topic,
			Form:  form,
		}))

		return
	}

	topic.Content = form.Data.Content

	topic, err = c.ts.Update(r.Context(), topic)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	Render(w, r, pages.TopicContentData(pages.TopicContentDataProps{
		Topic: topic,
		Form:  form,
	}))
}

func (c *TopicController) InfoGet(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	Render(w, r, pages.TopicInfo(pages.TopicInfoProps{Topic: topic}))
}

func (c *TopicController) InfoEdit(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	Render(w, r, pages.TopicInfoEdit(pages.TopicInfoEditProps{
		Topic: topic,
		Form:  htmx.NewForm(&pages.TopicInfoEditForm{Title: topic.Title}),
	}))
}

func (c *TopicController) DescendantList(w http.ResponseWriter, r *http.Request) {
	parentID, err := uuid.Parse(r.FormValue("parent"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	parent, err := c.ts.Get(r.Context(), parentID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), parent); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	descendants, err := c.ts.ListDescendants(r.Context(), parent)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	Render(w, r, pages.TopicDescendantList(pages.TopicDescendantListProps{
		Topic:       parent,
		Descendants: descendants,
		Level:       0,
	}))
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
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), parent); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	project, err := c.ps.Get(r.Context(), parent.ProjectID)
	if err != nil {
		RenderError(w, r, 500, err)
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
		RenderError(w, r, 500, err)
		return
	}

	w.Header().Set("HX-Trigger", "topic-list-updated")

	Render(w, r, pages.TopicDescendantCreate(pages.TopicDescendantCreateProps{
		Form: htmx.NewForm(&pages.TopicDescendantCreateForm{
			ParentID: form.Data.ParentID,
		}),
	}))
}

func (c *TopicController) DescendantListMove(w http.ResponseWriter, r *http.Request) {
	zone := r.FormValue("zone")
	isAbove := zone == "above"
	isBelow := zone == "below"

	destinationID, err := uuid.Parse(r.FormValue("destination"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	destination, err := c.ts.Get(r.Context(), destinationID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), destination); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	parent := destination
	if isAbove || isBelow {
		parent, err = c.ts.Get(r.Context(), destination.ParentID.UUID)
		if err != nil {
			RenderError(w, r, 500, err)
			return
		}
	}

	ancestors, err := c.ts.ListAncestors(r.Context(), parent)
	if err != nil {
		RenderError(w, r, 500, err)
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
				RenderError(w, r, 500, fmt.Errorf("cannot create a cycle in the tree"))
				return
			}
		}

		ID, err2 := uuid.Parse(sourceID)
		if err2 != nil {
			RenderError(w, r, 500, err2)
			return
		}

		source, err2 := c.ts.Get(r.Context(), ID)
		if err2 != nil {
			RenderError(w, r, 500, err2)
			return
		}

		if err = c.ts.Can(lib.MustGetUser(r.Context()), source); err != nil {
			RenderError(w, r, 403, err)
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
			RenderError(w, r, 500, err2)
			return
		}

		for index, source := range sources {
			err3 := c.ts.Move(r.Context(), source, parent, start+index)
			if err3 != nil {
				RenderError(w, r, 500, err3)
				return
			}
		}
	} else {
		children, err2 := c.ts.ListChildren(r.Context(), parent)
		if err2 != nil {
			RenderError(w, r, 500, err2)
			return
		}

		childrenCount := len(children)

		for index, source := range sources {
			err3 := c.ts.Move(r.Context(), source, parent, index+childrenCount+1)
			if err3 != nil {
				RenderError(w, r, 500, err3)
				return
			}
		}

		sess.MustGetSession(r.Context()).Add("expanded-topics", parent.ID.String())
	}

	// Self healing.
	project, err := c.ps.Get(r.Context(), parent.ProjectID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	err = c.ts.Reindex(r.Context(), project)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	w.Header().Set("HX-Trigger", "topic-list-updated")
	w.WriteHeader(204)
}

func (c *TopicController) DescendantDelete(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	if err = c.ts.Can(lib.MustGetUser(r.Context()), topic); err != nil {
		RenderError(w, r, 403, err)
		return
	}

	err = c.ts.Remove(r.Context(), topic)
	if err != nil {
		RenderError(w, r, 500, err)
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

func (c *TopicController) ContentCollapse(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	sess.MustGetSession(r.Context()).Remove("expanded-content", topic.ID.String())

	Render(w, r, pages.TopicContent(pages.TopicContentProps{
		Topic: topic,
		Form: htmx.NewForm(&pages.TopicContentForm{
			Content: topic.Content,
		}),
	}))
}

func (c *TopicController) ToolbarCollapseRecursive(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	descendants, err := c.ts.ListDescendants(r.Context(), topic)
	if err != nil {
		RenderError(w, r, 500, err)
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

func (c *TopicController) DescendantExpand(w http.ResponseWriter, r *http.Request) {
	sess.MustGetSession(r.Context()).Add("expanded-topics", r.FormValue("topic"))

	w.Header().Set("HX-Trigger", "topic-list-updated")
	w.WriteHeader(204)
}

func (c *TopicController) ToolbarExpandRecursive(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	descendants, err := c.ts.ListDescendants(r.Context(), topic)
	if err != nil {
		RenderError(w, r, 500, err)
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

func (c *TopicController) ContentExpand(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	sess.MustGetSession(r.Context()).Add("expanded-content", topic.ID.String())

	Render(w, r, pages.TopicContent(pages.TopicContentProps{
		Topic: topic,
		Form: htmx.NewForm(&pages.TopicContentForm{
			Content: topic.Content,
		}),
	}))
}

func (c *TopicController) ToolbarEnableSelectionMode(w http.ResponseWriter, r *http.Request) {
	topicID, err := uuid.Parse(r.FormValue("topic"))
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
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
		RenderError(w, r, 500, err)
		return
	}

	topic, err := c.ts.Get(r.Context(), topicID)
	if err != nil {
		RenderError(w, r, 500, err)
		return
	}

	sess.MustGetSession(r.Context()).Remove("mode", "selection")

	w.Header().Set("HX-Trigger", "topic-list-updated")
	Render(w, r, pages.TopicToolbar(pages.TopicToolbarProps{
		Topic: topic,
	}))
}
