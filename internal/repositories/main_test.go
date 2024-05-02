package repositories_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanpujo/blog-app/internal/repositories"
)

var (
	DBName = "db_test"
	DSN    = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=20"
)

var (
	testDB   *sql.DB
	blogRepo repositories.StoryRepository
	userRepo repositories.UserRepository
	mock     sqlmock.Sqlmock
)

// TestMain sets up the test environment using Docker to run a PostgreSQL container.
// It ensures that Docker is running, starts a new PostgreSQL container, and establishes a connection to it.
// After running the tests, it cleans up the resources.
func TestMain(m *testing.M) {
	var err error

	testDB, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("failed create mock db: %s", err)
	}
	defer testDB.Close()
	// Initialize the Docker pool.

	// Initialize repositories.
	userRepo = repositories.NewUserRepository(testDB)
	blogRepo = repositories.NewStoryRepository(testDB)

	// Run the tests.
	code := m.Run()

	// Exit with the status code from the test run.
	os.Exit(code)
}
