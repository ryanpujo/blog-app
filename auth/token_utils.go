package auth

import (
	"crypto"
	"crypto/sha256"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/blog-app/config"
	"github.com/ryanpujo/blog-app/utils"
)

type UserClaims struct {
	UserID uint `json:"id"`
	jwt.RegisteredClaims
}

var HMACMethod = &jwt.SigningMethodHMAC{
	Name: "claim",
	Hash: crypto.SHA256,
}

type refreshTokenGenerator struct {
	secret string
	repo   TokenRepository
}

func NewRefreshTokenGenerator(secret string, repo TokenRepository) refreshTokenGenerator {
	return refreshTokenGenerator{
		secret: secret,
		repo:   repo,
	}
}

func (r refreshTokenGenerator) GenerateToken(userID uint) (*string, error) {

	expiresAt := time.Now().Add(config.RefreshTokenExpiration)
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expiresAt,
			},
		},
	}

	token := jwt.NewWithClaims(HMACMethod, claims)

	tokenString, err := token.SignedString([]byte(r.secret))
	if err != nil {
		return nil, err
	}

	sha256 := sha256.Sum256([]byte(tokenString))

	hash, err := utils.HashPassword(string(sha256[:]))
	if err != nil {
		return nil, err
	}

	refreshToken := Token{
		TokenHash: hash,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	err = refreshToken.SaveToken(r.repo)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
