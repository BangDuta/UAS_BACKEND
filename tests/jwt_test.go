package tests

import (
	"os"
	"prestasi-mahasiswa-api/utils" // Pastikan ini sesuai nama module di go.mod
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTFlow(t *testing.T) {
	// Setup environment variable untuk testing
	os.Setenv("JWT_SECRET", "supersecretkey")

	userID := "user-uuid-test"
	role := "Mahasiswa"
	permissions := []string{"achievement:create"}

	// 1. Test Generate Token
	token, _, err := utils.GenerateAuthTokens(userID, role, permissions)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 2. Test Validate Token
	// Sekarang utils.ValidateToken seharusnya sudah terbaca
	claims, err := utils.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
}