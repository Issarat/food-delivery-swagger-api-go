package auth

import (
	"errors"
	"strings"
	"time"
	"crypto/rand"
	
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	key    []byte
	parser *jwt.Parser
}

func NewTokenManager(secret string) (*TokenManager, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 bytes")
	}

	return &TokenManager{
		key: []byte(secret),
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			jwt.WithIssuer("food-delivery-auth"),
			jwt.WithAudience("orders-api"),
			jwt.WithExpirationRequired(),
		),
	}, nil
}

func (m *TokenManager) Issue(userID string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", errors.New("user ID is required")
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		ID:        rand.Text(), 
		Subject:   userID,
		Issuer:    "food-delivery-auth",
		Audience:  jwt.ClaimStrings{"orders-api"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.key)
}

func (m *TokenManager) Verify(rawToken string) (*jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	_, err := m.parser.ParseWithClaims(rawToken, &claims, func(*jwt.Token) (any, error) {
		return m.key, nil
	})
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(claims.Subject) == "" {
		return nil, errors.New("token subject is required")
	}

	return &claims, nil
}