package lib

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/throskam/memo/internal/orm"
)

type Topic struct {
	ID        uuid.UUID
	Title     string
	Content   string
	SortOrder int

	ParentID  uuid.NullUUID
	ProjectID uuid.UUID

	CreatedAt time.Time

	Project *Project
}

type TopicService struct {
	queries *orm.Queries
}

func NewTopicService(queries *orm.Queries) *TopicService {
	return &TopicService{
		queries: queries,
	}
}

func (s *TopicService) ListChildren(ctx context.Context, t *Topic) ([]*Topic, error) {
	children := []*Topic{}

	rows, err := s.queries.ListTopicChildren(ctx, toPGUUID(t.ID))
	if err != nil {
		return children, fmt.Errorf("failed to list topic children [ID: %v]: %w", t.ID, err)
	}

	for _, row := range rows {
		child, err := fromTopicRow(row)
		if err != nil {
			return children, fmt.Errorf("failed to parse topic child [ID: %v]: %w", t.ID, err)
		}

		children = append(children, child)
	}

	return children, nil
}

func (s *TopicService) ListAncestors(ctx context.Context, t *Topic) ([]*Topic, error) {
	ancestors := []*Topic{}

	rows, err := s.queries.ListTopicAncestors(ctx, toPGUUID(t.ID))
	if err != nil {
		return ancestors, fmt.Errorf("failed to list topic ancestors [ID: %v]: %w", t.ID, err)
	}

	for _, row := range rows {
		ancestor, err := fromTopicRow(orm.Topic(row))
		if err != nil {
			return ancestors, fmt.Errorf("failed to parse topic row: %w", err)
		}

		ancestors = append(ancestors, ancestor)
	}

	ancestors = ancestors[1:]

	slices.Reverse(ancestors)

	return ancestors, nil
}

func (s *TopicService) ListDescendants(ctx context.Context, t *Topic) ([]*Topic, error) {
	descendants := []*Topic{}

	rows, err := s.queries.ListTopicDescendants(ctx, toPGUUID(t.ID))
	if err != nil {
		return descendants, fmt.Errorf("failed to list topic descendants [ID: %v]: %w", t.ID, err)
	}

	for _, row := range rows {
		if !row.ParentID.Valid {
			continue
		}

		descendant, err := fromTopicRow(orm.Topic(row))
		if err != nil {
			return descendants, fmt.Errorf("failed to parse topic row: %w", err)
		}

		descendants = append(descendants, descendant)
	}

	return descendants, nil
}

func (s *TopicService) Get(ctx context.Context, ID uuid.UUID) (*Topic, error) {
	row, err := s.queries.GetTopicByID(ctx, toPGUUID(ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get topic [ID: %v]: %w", ID, err)
	}

	topic, err := fromTopicRow(row.Topic)
	if err != nil {
		return nil, fmt.Errorf("failed to parse topic row [ID: %v]: %w", ID, err)
	}

	topic.Project, err = fromProjectRow(row.Project)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project row [ID: %v]: %w", topic.ProjectID, err)
	}

	return topic, nil
}

func (s *TopicService) Move(ctx context.Context, topic *Topic, parent *Topic, sortOrder int) error {
	err := s.queries.MoveTopic(ctx, orm.MoveTopicParams{
		ID:        toPGUUID(topic.ID),
		ParentID:  toPGUUID(parent.ID),
		SortOrder: int32(sortOrder),
	})
	if err != nil {
		return fmt.Errorf("failed to move topic: %w", err)
	}

	return nil
}

func (s *TopicService) Shift(ctx context.Context, parent *Topic, start, amount int) error {
	err := s.queries.ShiftTopics(ctx, orm.ShiftTopicsParams{
		ParentID: toPGUUID(parent.ID),
		Start:    int32(start),
		Amount:   int32(amount),
	})
	if err != nil {
		return fmt.Errorf("failed to shift topics: %w", err)
	}

	return nil
}

func (s *TopicService) Reindex(ctx context.Context, p *Project) error {
	err := s.queries.ReindexTopics(ctx, toPGUUID(p.ID))
	if err != nil {
		return fmt.Errorf("failed to reindex topic [ID: %v]: %w", p.ID, err)
	}

	return nil
}

func (s *TopicService) Create(ctx context.Context, t *Topic) (*Topic, error) {
	if t.ParentID.Valid {
		lastSortOrder, err := s.queries.GetLastSortOrder(ctx, toPGUUID(t.ParentID.UUID))
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("failed to get last sort order: %w", err)
			}
		}

		t.SortOrder = int(lastSortOrder + 1)

	}

	row, err := s.queries.CreateTopic(ctx, orm.CreateTopicParams{
		Title:     t.Title,
		Content:   t.Content,
		SortOrder: int32(t.SortOrder),
		ParentID:  toNullablePGUUID(t.ParentID),
		ProjectID: toPGUUID(t.ProjectID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	newTopic, err := fromTopicRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to parse row: %w", err)
	}

	return newTopic, nil
}

func (s *TopicService) Update(ctx context.Context, t *Topic) (*Topic, error) {
	row, err := s.queries.UpdateTopic(ctx, orm.UpdateTopicParams{
		ID:      toPGUUID(t.ID),
		Title:   t.Title,
		Content: t.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update topic [ID: %v]: %w", t.ID, err)
	}

	updatedTopic, err := fromTopicRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to parse row: %w", err)
	}

	return updatedTopic, nil
}

func (s *TopicService) Remove(ctx context.Context, t *Topic) error {
	if !t.ParentID.Valid {
		return fmt.Errorf("cannot remove root topic")
	}

	rows, err := s.queries.ListTopicChildren(ctx, toPGUUID(t.ID))
	if err != nil {
		return fmt.Errorf("failed to get topic children [ID: %v]: %w", t.ID, err)
	}

	children := []*Topic{}

	for _, row := range rows {
		child, err2 := fromTopicRow(row)
		if err2 != nil {
			return fmt.Errorf("failed to parse topic child [ID: %v]: %w", t.ID, err2)
		}

		children = append(children, child)
	}

	for _, child := range children {
		err2 := s.Remove(ctx, child)
		if err2 != nil {
			return fmt.Errorf("failed to remove topic child: %w", err2)
		}
	}

	err = s.queries.RemoveTopicByID(ctx, toPGUUID(t.ID))
	if err != nil {
		return fmt.Errorf("failed to remove topic [ID: %v]: %w", t.ID, err)
	}

	return nil
}

func (s *TopicService) Can(user *User, t *Topic) error {
	if t.Project.OwnerID != user.ID {
		return fmt.Errorf("cannot access the topic")
	}

	return nil
}

func fromTopicRow(row orm.Topic) (*Topic, error) {
	return &Topic{
		ID:        fromPGUUID(row.ID).UUID,
		Title:     row.Title,
		Content:   row.Content,
		SortOrder: int(row.SortOrder),
		ParentID:  fromPGUUID(row.ParentID),
		ProjectID: fromPGUUID(row.ProjectID).UUID,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}
