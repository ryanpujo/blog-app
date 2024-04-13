package repositories_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
	"github.com/stretchr/testify/require"
)

var payload = models.UserPayload{
	FirstName: "michael",
	LastName:  "townley",
	Username:  "townley",
	Password:  "fucktrevor",
	Email:     "townley@gmail.com",
}

func Test_userRepo_pingDB(t *testing.T) {
	err := testDB.Ping()
	require.Nil(t, err)
}

func Test_userRepo_Create(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualID *uint, err error)
	}{
		"success": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email).
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), *actualID)
			},
		},
		"failed": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email).
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.Error(t, err)
				require.Nil(t, actualID)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			user, err := userRepo.Create(payload)

			tc.assert(t, user, err)
		})
	}
}

func Test_userRepo_FindById(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualUser *models.User, err error)
	}{
		"success": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"}).
					AddRow(payload.ID, payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email, time.Now(), time.Now())

				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(1).WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualUser *models.User, err error) {
				require.NoError(t, err)
				require.Equal(t, payload.Username, actualUser.Username)
			},
		},
		"failed": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"})

				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(1).WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualUser *models.User, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			user, err := userRepo.FindById(1)

			tc.assert(t, user, err)
		})
	}
}

func Test_userRepo_FindUsers(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualUsers []*models.User, err error)
	}{
		"success": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"})
				for range 2 {
					rows.AddRow(payload.ID, payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email, time.Now(), time.Now())
				}

				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualUsers []*models.User, err error) {
				require.NoError(t, err)
				require.Equal(t, 2, len(actualUsers))
			},
		},
		"failed": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"})

				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(utils.ErrNoDataFound)
			},
			assert: func(t *testing.T, actualUsers []*models.User, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
				require.Nil(t, actualUsers)
			},
		},
		"scan error": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"})
				for range 2 {
					rows.AddRow(payload.ID, payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email, 1, time.Now())
				}

				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualUsers []*models.User, err error) {
				require.Error(t, err)
				require.Nil(t, actualUsers)
			},
		},
		"row error": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"}).
					AddRow(payload.ID, payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email, 1, time.Now()).
					RowError(0, utils.ErrNoDataFound)

				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualUsers []*models.User, err error) {
				require.Error(t, err)
				require.Nil(t, actualUsers)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			users, err := userRepo.FindUsers()

			tc.assert(t, users, err)
		})
	}
}

func Test_userRepo_DeleteById(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectExec("DELETE FROM users").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectExec("DELETE FROM users").WithArgs(1).WillReturnError(utils.ErrNoDataFound)
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"no record found": {
			arrange: func() {
				sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectExec("DELETE FROM users").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"result error": {
			arrange: func() {
				sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectExec("DELETE FROM users").WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(utils.ErrNoDataFound))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			err := userRepo.DeleteById(1)

			tc.assert(t, err)
		})
	}
}

func Test_userRepo_Update(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"})

				mock.ExpectExec("UPDATE users SET").WithArgs(
					payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email, 1,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"})

				mock.ExpectExec("UPDATE users SET").WithArgs(
					payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email, 1,
				).WillReturnError(utils.ErrNoDataFound)
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"no record found": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"})

				mock.ExpectExec("UPDATE users SET").WithArgs(
					payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email, 1,
				).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"result error": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "password", "email", "created_at", "updated_at"})

				mock.ExpectExec("UPDATE users SET").WithArgs(
					payload.FirstName, payload.LastName, payload.Username, payload.Password, payload.Email, 1,
				).WillReturnResult(sqlmock.NewErrorResult(utils.ErrNoDataFound))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			err := userRepo.Update(1, &payload)

			tc.assert(t, err)
		})
	}
}

func Test_userRepo_CheckIfUsernameOrEmailExists(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, isExists bool)
	}{
		"should be true": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

				mock.ExpectQuery("SELECT EXISTS").WithArgs("johndoe", "john").
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, isExists bool) {
				require.True(t, isExists)
			},
		},
		"should be false": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(time.Now())

				mock.ExpectQuery("SELECT EXISTS").WithArgs("johndoe", "john").
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, isExists bool) {
				require.False(t, isExists)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			isExists := userRepo.CheckIfEmailOrUsernameExist("johndoe", "john")

			tc.assert(t, isExists)
		})
	}
}
