package lib

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/throskam/memo/internal/orm"
)

type AuthenticationMethod struct {
	Provider orm.AuthenticationProvider
	Sub      string

	UserID uuid.UUID
}

func NewAuthenticationMethod(provider orm.AuthenticationProvider, sub string) *AuthenticationMethod {
	return &AuthenticationMethod{
		Provider: provider,
		Sub:      sub,
	}
}

type AuthenticationMethodService struct {
	queries *orm.Queries
}

func NewAuthenticationMethodService(queries *orm.Queries) *AuthenticationMethodService {
	return &AuthenticationMethodService{
		queries: queries,
	}
}

func (s *AuthenticationMethodService) Create(ctx context.Context, m *AuthenticationMethod) (*AuthenticationMethod, error) {
	row, err := s.queries.CreateAuthenticationMethod(ctx, orm.CreateAuthenticationMethodParams{
		Provider: m.Provider,
		Sub:      m.Sub,
		UserID:   toPGUUID(m.UserID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication method: %w", err)
	}

	authenticationMethod, err := fromAuthenticationMethodRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authentication method: %w", err)
	}

	return authenticationMethod, nil
}

func fromAuthenticationMethodRow(row orm.AuthenticationMethod) (*AuthenticationMethod, error) {
	return &AuthenticationMethod{
		Provider: row.Provider,
		Sub:      row.Sub,
		UserID:   fromPGUUID(row.UserID).UUID,
	}, nil
}
