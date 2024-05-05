package auth

import (
	"crypto"
	"crypto/sha256"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/blog-app/config"
	"github.com/ryanpujo/blog-app/database"
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

func GenerateRefreshToken(userID uint) (*string, error) {
	cfg := config.Config()
	expiresAt := time.Now().Add(config.RefreshTokenExpiration)
	claim := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expiresAt,
			},
		},
	}
	token := jwt.NewWithClaims(HMACMethod, claim)

	tokenString, err := token.SignedString([]byte(cfg.JWT.RefreshTokenSecret))
	if err != nil {
		return nil, err
	}

	sha256 := sha256.Sum256([]byte(tokenString))

	hash, err := utils.HashPassword(string(sha256[:]))
	if err != nil {
		log.Println("disini woyyyyy", err.Error())
		return nil, err
	}

	refreshToken := RefreshToken{
		TokenHash: hash,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	err = refreshToken.SaveToken(database.EstablishDBConnectionWithRetry())
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
