package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTCustomClaims diubah namanya agar sinkron dengan middleware
type JWTCustomClaims struct {
	UserID      uuid.UUID `json:"userId"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateAuthTokens(userID string, role string, permissions []string) (string, string, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return "", "", fmt.Errorf("invalid user ID: %w", err)
	}

	claims := &JWTCustomClaims{
		UserID:      id,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, "refresh_token_placeholder", nil
}

// Nama fungsi diubah menjadi ValidateToken agar sinkron dengan pemanggilan middleware
func ValidateToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}