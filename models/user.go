package models

import (
	"time"
	"github.com/google/uuid"
)

// User merepresentasikan data dari tabel users (PostgreSQL)
type User struct {
    ID           uuid.UUID `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    FullName     string    `json:"fullName"`
    RoleID       uuid.UUID `json:"roleId"`
    Role         string    `json:"role"` 
    Permissions  []string  `json:"permissions"`
    IsActive     bool      `json:"isActive"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}

// LoginRequest untuk payload login
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
// ... tambahkan UserProfile dan LoginResponse (sama seperti di jawaban sebelumnya)