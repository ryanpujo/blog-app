package auth

import (
	"context"
	"database/sql"
	"fmt"
)

// RefreshToken is a type that implements the TokenSaver interface.
type RefreshToken struct {
	db *sql.DB
}

// NewRefreshToken creates a new RefreshToken.
func NewRefreshToken(db *sql.DB) RefreshToken {
	return RefreshToken{
		db: db,
	}
}

// SaveToken saves a Token to the database.
func (r RefreshToken) SaveToken(ctx context.Context, t Token) error {
	stmt := `
		INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
			VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, stmt,
		t.TokenHash,
		t.UserID,
		t.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}
