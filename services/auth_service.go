package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/repositories"
	"prestasi-mahasiswa-api/utils"

	"github.com/google/uuid"
)

type AuthService interface {
	PerformLogin(ctx context.Context, username, password string) (*models.LoginResponse, int, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, int, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// PerformLogin mengimplementasikan flow login (FR-001)
func (s *authService) PerformLogin(ctx context.Context, username, password string) (*models.LoginResponse, int, error) {
	// 1. Dapatkan user dari database
	user, err := s.userRepo.FindUserByUsernameOrEmail(ctx, username)
	if err != nil {
		// Asumsi repository mengembalikan error saat user not found
		return nil, http.StatusUnauthorized, errors.New("Invalid credentials")
	}

	// 2. Sistem memvalidasi kredensial
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, http.StatusUnauthorized, errors.New("Invalid credentials")
	}

	// 3. Sistem mengecek status aktif user
	if !user.IsActive {
		return nil, http.StatusForbidden, errors.New("User account is inactive")
	}

	// 4. Sistem generate JWT token dengan role dan permissions
	accessToken, refreshToken, err := utils.GenerateAuthTokens(user.ID.String(), user.Role, user.Permissions)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to generate token: %w", err)
	}

	// 5. Return token dan user profile
	profile := models.UserProfile{
		ID:          user.ID.String(),
		Username:    user.Username,
		FullName:    user.FullName,
		Role:        user.Role,
		Permissions: user.Permissions,
	}

	resp := &models.LoginResponse{
		Status: "success",
		Data: models.LoginData{
			Token:        accessToken,
			RefreshToken: refreshToken,
			User:         profile,
		},
	}
	return resp, http.StatusOK, nil
}

// GetProfile (dipanggil dari AuthRequired middleware atau route profile)
func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, int, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("User profile not found")
	}

	profile := models.UserProfile{
		ID:          user.ID.String(),
		Username:    user.Username,
		FullName:    user.FullName,
		Role:        user.Role,
		Permissions: user.Permissions,
	}
	return &profile, http.StatusOK, nil
}