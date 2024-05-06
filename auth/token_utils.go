package auth

import (
	"context"
	"crypto"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token represents a token.
type Token struct {
	TokenHash string
	UserID    uint
	ExpiresAt time.Time
	Revoked   bool
}

// UserClaims struct represents the claims that are encoded into the JWT.
// UserID is the ID of the user. jwt.RegisteredClaims are a struct provided by the jwt package that includes standard claim fields such as issuer and expiration time.
type UserClaims struct {
	UserID uint `json:"id"`
	jwt.RegisteredClaims
}

// HMACMethod is the signing method used to create the JWT. It uses HMAC with SHA-256.
var HMACMethod = &jwt.SigningMethodHMAC{
	Name: "claim",
	Hash: crypto.SHA256,
}

// refreshTokenGenerator is a struct that holds the secret used to sign the JWT and a repository where the token can be saved.
type tokenGenerator struct {
	secret    string
	expiresAt time.Time
	repo      TokenSaver
}

// NewTokenGenerator is a constructor function that returns a new tokenGenerator.
// It takes a secret string and a TokenRepository as parameters.
func NewTokenGenerator(secret string, repo TokenSaver, expiresAt time.Time) tokenGenerator {
	return tokenGenerator{
		secret:    secret,
		repo:      repo,
		expiresAt: expiresAt,
	}
}

// GenerateToken is a method on tokenGenerator that generates a new JWT and saves it in the repository.
// It takes a userID as a parameter. The userID is included in the claims of the JWT.
// The method returns a pointer to the JWT string and an error. If there is an error at any point in the method, it will return the error.
func (r tokenGenerator) GenerateToken(userID uint) (*string, error) {

	// Create the claims for the JWT
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: r.expiresAt,
			},
		},
	}

	// Create a new JWT with the claims
	jwtToken := jwt.NewWithClaims(HMACMethod, claims)

	// Sign the JWT with the secret
	tokenString, err := jwtToken.SignedString([]byte(r.secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	// Hash the JWT string
	hash := sha256.Sum256([]byte(tokenString))

	// Create a new Token struct
	token := Token{
		TokenHash: string(hash[:]),
		UserID:    userID,
		ExpiresAt: r.expiresAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Save the Token struct in the repository
	err = r.repo.SaveToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	// Return a pointer to the JWT string
	return &tokenString, nil
}
