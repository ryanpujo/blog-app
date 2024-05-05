package auth_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/blog-app/auth"
	"github.com/ryanpujo/blog-app/config"
)

var (
	refreshToken = auth.RefreshToken{
		TokenHash: "slkfekfne",
		UserID:    1,
		ExpiresAt: time.Now().Add(config.RefreshTokenExpiration),
	}

	testDB  *sql.DB
	SQLMock sqlmock.Sqlmock
)

func TestMain(m *testing.M) {
	var err error
	testDB, SQLMock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("failed create mock db: %s", err)
	}
	defer testDB.Close()

	os.Exit(m.Run())
}
