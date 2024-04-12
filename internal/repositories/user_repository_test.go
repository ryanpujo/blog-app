package repositories_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ryanpujo/blog-app/utils"

	"github.com/ryanpujo/blog-app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func Test_userRepo_pingDB(t *testing.T) {
	err := testDB.Ping()
	require.Nil(t, err)
}

func Test_userRepo_Create(t *testing.T) {
	id, err := userRepo.Create(payload)
	log.Println("id sini", *id)
	require.Equal(t, uint(11), *id)
	require.Nil(t, err)

	id, err = userRepo.Create(payload1)
	require.NoError(t, err)
	require.Equal(t, uint(12), *id)

	id, err = userRepo.Create(payload)

	require.Nil(t, id)
	require.Error(t, err)
	var pgErr *pgconn.PgError
	if assert.ErrorAs(t, err, &pgErr) {
		require.Equal(t, utils.ErrCodeUniqueViolation, pgErr.Code)
	}
}

func Test_userRepo_FindById(t *testing.T) {
	user, err := userRepo.FindById(uint(1))
	require.NoError(t, err)
	require.Equal(t, "John", user.FirstName)

	user, err = userRepo.FindById(uint(13))
	require.Error(t, err)
	require.Nil(t, user)
	if assert.ErrorAs(t, err, &sql.ErrNoRows) {
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func Test_userRepo_FindUsers(t *testing.T) {
	users, err := userRepo.FindUsers()
	require.NoError(t, err)
	require.Equal(t, len(users), 12)
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
	require.Equal(t, 11, len(users))
}

func Test_userRepo_Update(t *testing.T) {
	user, err := userRepo.FindById(11)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, payload.LastName, user.LastName)

	var updated = models.UserPayload{
		ID:        11,
		FirstName: "michael",
		LastName:  "de santa",
		Username:  "townley",
		Password:  "fucktrevor",
		Email:     "townley@gmail.com",
	}
	err = userRepo.Update(updated.ID, &updated)
	require.NoError(t, err)

	user, err = userRepo.FindById(11)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, updated.LastName, user.LastName)

	err = userRepo.Update(2, &updated)
	require.Error(t, err)
	require.Equal(t, fmt.Errorf("no record found with id %d to update: %w", 2, sql.ErrNoRows), err)
}

func Test_userRepo_CheckIfUsernameOrEmailExists(t *testing.T) {
	isExists := userRepo.CheckIfEmailOrUsernameExist(payload.Email, payload.Username)
	require.True(t, isExists)

	isExists = userRepo.CheckIfEmailOrUsernameExist("desanta@gmail.com", "okeoke")
	require.False(t, isExists)
}
