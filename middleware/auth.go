package middleware

import (
	"context"
	"net/http"
	"strings"

	"prestasi-mahasiswa-api/utils" // Pastikan import ini benar

	"github.com/google/uuid"
)

// Kunci Context agar tidak bentrok
type contextKey string

const UserClaimsKey contextKey = "userClaims"

// AuthRequired memvalidasi JWT
func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing or invalid token")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// RBACRequired mengecek permission
func RBACRequired(requiredPermission string, next http.HandlerFunc) http.HandlerFunc {
	return AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserClaims(r.Context())

		hasPermission := false
		for _, p := range claims.Permissions {
			if p == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.RespondWithError(w, http.StatusForbidden, "Insufficient permissions")
			return
		}

		next(w, r)
	})
}

// Helper untuk mengambil claims dari context
func GetUserClaims(ctx context.Context) *utils.Claims {
	if claims, ok := ctx.Value(UserClaimsKey).(*utils.Claims); ok {
		return claims
	}
	return &utils.Claims{UserID: uuid.Nil}
}