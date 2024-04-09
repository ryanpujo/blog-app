package repositories_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/ryanpujo/blog-app/internal/user/repositories"
	"github.com/ryanpujo/blog-app/utils"

	"github.com/ryanpujo/blog-app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	host     = "localhost"
	port     = "5435"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=20"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDb *sql.DB
var userRepo repositories.UserRepository

var payload = models.UserPayload{
	FirstName: "michael",
	LastName:  "townley",
	Username:  "townley",
	Password:  "fucktrevor",
	Email:     "townley@gmail.com",
}
var payload1 = models.UserPayload{
	FirstName: "michael",
	LastName:  "townley",
	Username:  "townley1",
	Password:  "fucktrevor",
	Email:     "townley1@gmail.com",
}

func TestMain(m *testing.M) {
	// connect to docker; fail if docker not running
	p, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("could not connect to docker, is it running? %s", err)
	}

	pool = p

	// setup docker options, specifying the image and so forth
	opt := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.2-alpine",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource (docker image)
	resource, err = pool.RunWithOptions(&opt)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until its ready
	if err := pool.Retry(func() error {
		var err error
		testDb, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("error:", err)
		}
		return testDb.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	// populate database with empty table
	err = createTables()
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("cant create table: %s", err)
	}

	userRepo = repositories.NewUserRepository(testDb)
	code := m.Run()

	// clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("cant clean resources: %s", err)
	}

	os.Exit(code)
}

func createTables() (err error) {
	var tableSql []byte
	tableSql, err = os.ReadFile("../../../sql/query.sql")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = testDb.Exec(string(tableSql))
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func Test_userRepo_pingDB(t *testing.T) {
	err := testDb.Ping()
	require.Nil(t, err)
}

func Test_userRepo_Create(t *testing.T) {

	id, err := userRepo.Create(payload)
	require.Equal(t, uint(1), id)
	require.Nil(t, err)

	id, err = userRepo.Create(payload1)
	require.NoError(t, err)
	require.Equal(t, uint(2), id)

	id, err = userRepo.Create(payload)

	require.Equal(t, uint(0), id)
	require.Error(t, err)
	var pgErr *pgconn.PgError
	if assert.ErrorAs(t, err, &pgErr) {
		require.Equal(t, utils.ErrCodeUniqueViolation, pgErr.Code)
	}
}

func Test_userRepo_FindById(t *testing.T) {
	user, err := userRepo.FindById(uint(1))
	require.NoError(t, err)
	require.Equal(t, payload.FirstName, user.FirstName)

	user, err = userRepo.FindById(uint(3))
	require.Error(t, err)
	require.Nil(t, user)
	if assert.ErrorAs(t, err, &sql.ErrNoRows) {
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func Test_userRepo_FindUsers(t *testing.T) {
	users, err := userRepo.FindUsers()
	require.NoError(t, err)
	require.Equal(t, len(users), 2)
}

func Test_userRepo_DeleteById(t *testing.T) {
	err := userRepo.DeleteById(2)
	require.NoError(t, err)

	user, err := userRepo.FindById(2)
	require.Error(t, err)
	if assert.ErrorAs(t, err, &sql.ErrNoRows) {
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
	require.Nil(t, user)

	users, err := userRepo.FindUsers()
	require.NoError(t, err)
	require.Equal(t, 1, len(users))
}

func Test_userRepo_Update(t *testing.T) {
	user, err := userRepo.FindById(1)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, payload.LastName, user.LastName)

	var updated = models.UserPayload{
		ID:        1,
		FirstName: "michael",
		LastName:  "de santa",
		Username:  "townley",
		Password:  "fucktrevor",
		Email:     "townley@gmail.com",
	}
	err = userRepo.Update(updated.ID, &updated)
	require.NoError(t, err)

	user, err = userRepo.FindById(1)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, updated.LastName, user.LastName)

	err = userRepo.Update(2, &updated)
	require.Error(t, err)
	require.Equal(t, "no record found with id 2 to update", err.Error())
}
