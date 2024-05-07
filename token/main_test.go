package token_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/blog-app/token"
)

var (
	testDB                *sql.DB
	SQLMock               sqlmock.Sqlmock
	repoMock              = new(RepoMock)
	refreshTokenGenerator token.TokenGenerator
)

func TestMain(m *testing.M) {
	var err error
	testDB, SQLMock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("failed create mock db: %s", err)
	}
	refreshTokenGenerator = token.NewTokenGenerator("lejjfelkjrrlekjssrjlejr", repoMock, time.Now().Add(time.Minute*60))
	defer testDB.Close()

	os.Exit(m.Run())
}
