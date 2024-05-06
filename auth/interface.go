package auth

import (
	"context"
)

type TokenSaver interface {
	SaveToken(ctx context.Context, t Token) error
}

type TokenGenerator interface {
	GenerateToken(userID uint) (*string, error)
}
