package middleware

import (
	"prestasi-mahasiswa-api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization header")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	claims, err := utils.ValidateToken(tokenString) // Sekarang sudah sinkron
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
	}

	c.Locals("user", claims)
	return c.Next()
}

// Mengembalikan pointer ke JWTCustomClaims agar sinkron
func GetUserClaims(c *fiber.Ctx) *utils.JWTCustomClaims {
	user := c.Locals("user")
	if user == nil {
		return nil
	}
	return user.(*utils.JWTCustomClaims)
}

func RBACRequired(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := GetUserClaims(c)
		if claims == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}
		
		hasPermission := false
		for _, p := range claims.Permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Forbidden: insufficient permissions")
		}

		return c.Next()
	}
}