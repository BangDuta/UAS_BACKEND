package utils

import (
	"time"
	"os"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims berisi data yang disimpan dalam JWT
type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Ambil dari env

// GenerateAuthTokens membuat token akses dan refresh
func GenerateAuthTokens(userID string, role string, permissions []string) (string, string, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return "", "", fmt.Errorf("invalid user ID: %w", err)
	}

	// Token Akses (Contoh: kadaluarsa dalam 1 jam)
	accessClaims := &Claims{
		UserID: id,
		Role:   role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Token Refresh (Contoh: kadaluarsa dalam 7 hari)
	refreshClaims := &Claims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// ValidateJWT memvalidasi token dan mengembalikan claims
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}