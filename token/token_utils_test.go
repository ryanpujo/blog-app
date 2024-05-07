package token_test

import (
	"context"
	"crypto"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/blog-app/token"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type RepoMock struct {
	mock.Mock
}

func (m *RepoMock) SaveToken(ctx context.Context, t token.Token) error {
	arg := m.Called(ctx, t)
	return arg.Error(0)
}

var HMAC = token.HMACMethod

func Test_Generate_RefreshToken(t *testing.T) {
	testTable := map[string]struct {
		arrange  func()
		assert   func(t *testing.T, actual *string, err error)
		tearDown func()
	}{
		"success": {
			arrange: func() {
				repoMock.On("SaveToken", mock.Anything, mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, actual *string, err error) {
				require.NoError(t, err)
				require.NotNil(t, actual)
			},
			tearDown: func() {
			},
		},
		"fail sign": {
			arrange: func() {
				token.HMACMethod = &jwt.SigningMethodHMAC{
					Name: "claims",
					Hash: crypto.BLAKE2b_384,
				}
			},
			assert: func(t *testing.T, actual *string, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.Equal(t, "failed to sign token: the requested hash function is unavailable", err.Error())
			},
			tearDown: func() {
				token.HMACMethod = HMAC
			},
		},
		"failed to save": {
			arrange: func() {
				repoMock.On("SaveToken", mock.Anything, mock.Anything).Return(errors.New("failed to save")).Once()
			},
			assert: func(t *testing.T, actual *string, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.Equal(t, "failed to save token: failed to save", err.Error())
			},
			tearDown: func() {
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			token, err := refreshTokenGenerator.GenerateToken(1)

			tc.assert(t, token, err)

			tc.tearDown()
		})
	}
}
