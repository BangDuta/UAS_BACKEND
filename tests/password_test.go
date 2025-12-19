package tests

import (
	"prestasi-mahasiswa-api/utils"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "rahasia123"

	// Test Hashing
	hash, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, password, hash)

	// Test Verification
	match := utils.CheckPasswordHash(password, hash)
	assert.True(t, match)

	// Test Wrong Password
	noMatch := utils.CheckPasswordHash("salah", hash)
	assert.False(t, noMatch)
}