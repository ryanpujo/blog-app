package auth

import (
	"context"
	"database/sql"
	"time"
)

type Token struct {
	TokenHash string
	UserID    uint
	ExpiresAt time.Time
	Revoked   bool
}

type refreshToken struct {
	db *sql.DB
}

func NewRefreshToken(db *sql.DB) refreshToken {
	return refreshToken{
		db: db,
	}
}

func (r refreshToken) SaveToken(ctx context.Context, t Token) error {
	stmt := `
		INSERT INTO refresh_tokens (token_hase, user_id, expires_at)
			VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, stmt,
		t.TokenHash,
		t.UserID,
		t.ExpiresAt,
	)
	return err
}

func (t Token) SaveToken(tRepo TokenRepository) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err := tRepo.SaveToken(ctx, t)
	return err
}
