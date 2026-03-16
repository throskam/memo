package lib

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/throskam/memo/internal/orm"
)

type Project struct {
	ID uuid.UUID

	OwnerID uuid.UUID

	CreatedAt time.Time

	Root  *Topic
	Owner *User
}

type ProjectService struct {
	queries *orm.Queries
}

func NewProjectService(queries *orm.Queries) *ProjectService {
	return &ProjectService{
		queries: queries,
	}
}

func (s *ProjectService) ListByOwnerWithRoot(ctx context.Context, owner *User) ([]*Project, error) {
	projects := []*Project{}

	rows, err := s.queries.ListProjectsByOwnerWithRoot(ctx, toPGUUID(owner.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to get project by owner [ID: %v]: %w", owner.ID, err)
	}

	for _, row := range rows {
		project, err := fromProjectRow(row.Project)
		if err != nil {
			return nil, fmt.Errorf("failed to parse project: %w", err)
		}

		project.Root, err = fromTopicRow(row.Topic)
		if err != nil {
			return nil, fmt.Errorf("failed to parse project root: %w", err)
		}

		project.Owner = owner

		projects = append(projects, project)
	}

	return projects, nil
}

func (s *ProjectService) Get(ctx context.Context, ID uuid.UUID) (*Project, error) {
	row, err := s.queries.GetProjectByID(ctx, toPGUUID(ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get project [ID: %v]: %w", ID, err)
	}

	project, err := fromProjectRow(row.Project)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}

	project.Root, err = fromTopicRow(row.Topic)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project root: %w", err)
	}

	project.Owner, err = fromUserRow(row.User)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project owner: %w", err)
	}

	return project, nil
}

func (s *ProjectService) Create(ctx context.Context, p *Project) (*Project, error) {
	row, err := s.queries.CreateProject(ctx, toPGUUID(p.OwnerID))
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	project, err := fromProjectRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) Remove(ctx context.Context, p *Project) error {
	err := s.queries.RemoveTopicsByProjectID(ctx, toPGUUID(p.ID))
	if err != nil {
		return fmt.Errorf("failed to remove topics by project [ID: %v]: %w", p.ID, err)
	}

	return s.queries.RemoveProjectByID(ctx, toPGUUID(p.ID))
}

func (s *ProjectService) Can(user *User, p *Project) error {
	if p.OwnerID != user.ID {
		return fmt.Errorf("cannot access the project")
	}

	return nil
}

func fromProjectRow(row orm.Project) (*Project, error) {
	return &Project{
		ID:        fromPGUUID(row.ID).UUID,
		OwnerID:   fromPGUUID(row.OwnerID).UUID,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}