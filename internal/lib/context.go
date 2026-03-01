package lib

import (
	"context"

	"github.com/throskam/ki"
	"github.com/throskam/kix/auth"
)

func GetUser(ctx context.Context) (*User, error) {
	return auth.GetIdentity[*User](ctx)
}

func MustGetUser(ctx context.Context) *User {
	return auth.MustGetIdentity[*User](ctx)
}

func GetLocation(ctx context.Context, key string) ki.Location {
	return ki.GetLocation(ctx, key)
}
