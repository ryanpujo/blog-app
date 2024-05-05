package auth_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func Test_Save_Token(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				SQLMock.ExpectExec("INSERT INTO tokens").
					WithArgs(refreshToken.TokenHash, refreshToken.UserID, refreshToken.ExpiresAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				SQLMock.ExpectExec("INSERT INTO tokens").
					WithArgs(refreshToken.TokenHash, refreshToken.UserID, refreshToken.ExpiresAt).
					WillReturnError(errors.New("failed"))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "failed", err.Error())
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			err := refreshToken.SaveToken(testDB)

			tc.assert(t, err)
		})
	}
}
