package auth_test

import (
	"context"
	"crypto"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/blog-app/auth"
	"github.com/ryanpujo/blog-app/config"
	"github.com/ryanpujo/blog-app/database"
	"github.com/ryanpujo/blog-app/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type DBMock struct {
	mock.Mock
}

func (m *DBMock) ExecContext(context context.Context, query string, args ...any) (sql.Result, error) {
	arg := m.Called(context, query, args)
	return arg.Get(0).(sql.Result), arg.Error(1)
}

func (m *DBMock) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	arg := m.Called(ctx, query, args)
	return arg.Get(0).(*sql.Rows), arg.Error(1)
}

func (m *DBMock) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	arg := m.Called(ctx, query, args)
	return arg.Get(0).(*sql.Row)
}

type Result struct{}

func (r Result) LastInsertId() (int64, error) {
	return 0, nil
}

func (r Result) RowsAffected() (int64, error) {
	return 0, nil
}

var (
	HMAC      = auth.HMACMethod
	db        = database.EstablishDBConnectionWithRetry
	cfg       = config.Config
	mockDB    = new(DBMock)
	bcrypHash = utils.HashPassword
)

func arrange() {
	database.EstablishDBConnectionWithRetry = func() database.DatabaseOperations {
		return mockDB
	}
	config.Config = func() config.Configuration {
		return config.Configuration{
			JWT: config.JWTConfig{
				RefreshTokenSecret: `1f597b3488697e817abe7222b38f70d4f6392e4568fafd57e0b77de6a2092b48`,
			},
		}
	}
}

func tearDown() {
	database.EstablishDBConnectionWithRetry = db
	config.Config = cfg
}

func Test_Generate_RefreshToken(t *testing.T) {
	testTable := map[string]struct {
		arrange  func()
		assert   func(t *testing.T, actual *string, err error)
		tearDown func()
	}{
		"success": {
			arrange: func() {
				arrange()
				mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(Result{}, nil).Once()
			},
			assert: func(t *testing.T, actual *string, err error) {
				require.NoError(t, err)
				require.NotNil(t, actual)
			},
			tearDown: func() {
				tearDown()
			},
		},
		"fail sign": {
			arrange: func() {
				arrange()
				auth.HMACMethod = &jwt.SigningMethodHMAC{
					Name: "claims",
					Hash: crypto.BLAKE2b_384,
				}
			},
			assert: func(t *testing.T, actual *string, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
			},
			tearDown: func() {
				tearDown()
				auth.HMACMethod = HMAC
			},
		},
		"failed bcrypt": {
			arrange: func() {
				arrange()
				utils.HashPassword = func(plain string) (string, error) {
					return "", errors.New("failed")
				}
			},
			assert: func(t *testing.T, actual *string, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.Equal(t, "failed", err.Error())
			},
			tearDown: func() {
				tearDown()
				utils.HashPassword = bcrypHash
			},
		},
		"failed to save": {
			arrange: func() {
				arrange()
				mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(Result{}, errors.New("failed to save")).Once()
			},
			assert: func(t *testing.T, actual *string, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.Equal(t, "failed to save", err.Error())
			},
			tearDown: func() {
				tearDown()
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			token, err := auth.GenerateRefreshToken(1)

			tc.assert(t, token, err)

			tc.tearDown()
		})
	}
}
