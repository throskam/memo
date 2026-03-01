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

type User struct {
	ID       uuid.UUID
	Username string

	CreatedAt time.Time

	AuthenticationMethods []*AuthenticationMethod
}

type UserService struct {
	queries *orm.Queries

	ams *AuthenticationMethodService
}

func NewUserService(queries *orm.Queries, ams *AuthenticationMethodService) *UserService {
	return &UserService{
		queries: queries,
		ams:     ams,
	}
}

func (s *UserService) Get(ctx context.Context, ID uuid.UUID) (*User, error) {
	row, err := s.queries.GetUserByID(ctx, toPGUUID(ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get user [ID: %v]: %w", ID, err)
	}

	user, err := fromUserRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user row: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByAuthenticationMethodOrCreate(ctx context.Context, m *AuthenticationMethod) (*User, error) {
	row, err := s.queries.GetUserByAuthenticationMethod(ctx, orm.GetUserByAuthenticationMethodParams{
		Provider: m.Provider,
		Sub:      m.Sub,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user, err2 := s.Create(ctx, &User{
				Username: "anon",
			})
			if err2 != nil {
				return nil, fmt.Errorf("failed to create new user: %w", err2)
			}

			m.UserID = user.ID

			authenticationMethod, err2 := s.ams.Create(ctx, m)
			if err2 != nil {
				return nil, fmt.Errorf("failed to create authentication method for the new user: %w", err2)
			}

			user.AuthenticationMethods = []*AuthenticationMethod{authenticationMethod}

			return user, nil
		}

		return nil, fmt.Errorf("failed to get user by authentication method: %w", err)
	}

	user, err := fromUserRow(row.User)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user row: %w", err)
	}

	authenticationMethod, err := fromAuthenticationMethodRow(row.AuthenticationMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authentication method row: %w", err)
	}

	user.AuthenticationMethods = []*AuthenticationMethod{authenticationMethod}

	return user, nil
}

func (s *UserService) Create(ctx context.Context, u *User) (*User, error) {
	row, err := s.queries.CreateUser(ctx, u.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user, err := fromUserRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user row: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, u *User) (*User, error) {
	row, err := s.queries.UpdateUser(ctx, orm.UpdateUserParams{
		ID:       toPGUUID(u.ID),
		Username: u.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user [ID: %v]: %w", u.ID, err)
	}

	updatedTopic, err := fromUserRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to parse row: %w", err)
	}

	return updatedTopic, nil
}

func fromUserRow(row orm.User) (*User, error) {
	return &User{
		ID:                    fromPGUUID(row.ID).UUID,
		AuthenticationMethods: []*AuthenticationMethod{},
		Username:              row.Username,
		CreatedAt:             row.CreatedAt.Time,
	}, nil
}
