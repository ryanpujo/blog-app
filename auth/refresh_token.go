package auth

import (
	"context"
	"time"

	"github.com/ryanpujo/blog-app/database"
)

type RefreshToken struct {
	TokenHash string
	UserID    uint
	ExpiresAt time.Time
	Revoked   bool
}

func (r RefreshToken) SaveToken(db database.ExecContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	stmt := `INSERT INTO tokens (token_hash, user_id, expires_at)
						values ($1, $2, $3)
					`
	_, err := db.ExecContext(ctx, stmt,
		r.TokenHash,
		r.UserID,
		r.ExpiresAt,
	)
	return err
}
