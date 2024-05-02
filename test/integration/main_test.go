package integration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	lorem "github.com/derektata/lorem/ipsum"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/ryanpujo/blog-app/internal/registry"
	"github.com/ryanpujo/blog-app/internal/route"
)

var (
	host     = "localhost"
	port     = "5435"
	user     = "postgres"
	password = "postgres"
	DBName   = "db_test"
	DSN      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=20"
)

var (
	pool   *dockertest.Pool
	testDB *sql.DB
	mux    *gin.Engine
)

var loremGenerator *lorem.Generator

// TestMain sets up the test environment using Docker to run a PostgreSQL container.
// It ensures that Docker is running, starts a new PostgreSQL container, and establishes a connection to it.
// After running the tests, it cleans up the resources.
func TestMain(m *testing.M) {
	// Initialize the Docker pool.
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Failed to connect to Docker: %s", err)
	}

	// Define options for running the PostgreSQL container.
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.2-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", user),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
			fmt.Sprintf("POSTGRES_DB=%s", DBName),
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {{HostIP: "0.0.0.0", HostPort: port}},
		},
	}

	// Start the PostgreSQL container.
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Failed to run container: %s", err)
	}

	// Attempt to connect to the PostgreSQL container.
	if err := pool.Retry(func() error {
		var dbErr error
		testDB, dbErr = sql.Open("pgx", fmt.Sprintf(DSN, host, port, user, password, DBName))
		if dbErr != nil {
			return dbErr
		}
		return testDB.Ping()
	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Failed to connect to PostgreSQL: %s", err)
	}

	// Create tables in the test database.
	if err := createTables(); err != nil {
		pool.Purge(resource)
		log.Fatalf("Failed to create tables: %s", err)
	}

	// Initialize repositories.
	registry := registry.New(testDB)
	appController := registry.NewAppController()
	mux = route.Route(appController)
	loremGenerator = lorem.NewGenerator()
	// Run the tests.
	code := m.Run()

	// Clean up resources after tests have completed.
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Failed to clean resources: %s", err)
	}

	// Exit with the status code from the test run.
	os.Exit(code)
}

func createTables() (err error) {
	var tableSql []byte
	tableSql, err = os.ReadFile("../../sql/test.sql")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = testDB.Exec(string(tableSql))
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
