package middleware

import (
	"strings"

	"prestasi-mahasiswa-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthRequired memvalidasi JWT untuk Fiber
func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing or invalid token",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
		})
	}

	// Simpan claims ke Locals (Context Fiber)
	c.Locals("userClaims", claims)
	return c.Next()
}

// RBACRequired mengecek permission (FR-002)
func RBACRequired(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Pastikan AuthRequired sudah dijalankan sebelumnya
		claims := GetUserClaims(c)

		hasPermission := false
		for _, p := range claims.Permissions {
			if p == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// Helper untuk mengambil claims dari Fiber Context
func GetUserClaims(c *fiber.Ctx) *utils.Claims {
	claims, ok := c.Locals("userClaims").(*utils.Claims)
	if !ok {
		return &utils.Claims{UserID: uuid.Nil}
	}
	return claims
}