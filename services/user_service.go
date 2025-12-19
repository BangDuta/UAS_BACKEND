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

type UserService interface {
	ListAllUsers(ctx context.Context) ([]models.User, int, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, int, error)
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, int, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req *models.UpdateUserRequest) (*models.User, int, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (int, error) // Deactivate

	SetStudentProfile(ctx context.Context, userID uuid.UUID, req *models.StudentProfileRequest) (*models.Student, int, error)
	SetLecturerProfile(ctx context.Context, userID uuid.UUID, req *models.LecturerProfileRequest) (*models.Lecturer, int, error)
	AssignAdvisor(ctx context.Context, studentUserID uuid.UUID, advisorUserID uuid.UUID) (int, error)

}

type userService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
	profileRepo repositories.ProfileRepository
    // ... (opsional: student/lecturer repo jika logic set profile ada di sini)
}

func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, profileRepo repositories.ProfileRepository) UserService {
	return &userService{
			userRepo: userRepo, 
			roleRepo: roleRepo, 
			profileRepo: profileRepo,}
}

// ListAllUsers
func (s *userService) ListAllUsers(ctx context.Context) ([]models.User, int, error) {
	users, err := s.userRepo.ListAllUsers(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return users, http.StatusOK, nil
}

// GetUserByID
func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, int, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	return user, http.StatusOK, nil
}

// CreateUser
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, int, error) {
	// 1. Cek duplikasi username/email
    // ... (logic untuk pengecekan duplikasi) ...

	// 2. Dapatkan RoleID
	role, err := s.roleRepo.GetRoleByName(ctx, req.RoleName)
	if err != nil || role == nil {
		return nil, http.StatusBadRequest, errors.New("invalid role name specified")
	}

	// 3. Hash Password dan Buat User Object (menggunakan utils/password.go)
    // ... (logic hashing password) ...
    hashedPassword, _ := utils.HashPassword(req.Password)
	
	newUser := &models.User{
		ID: uuid.New(), Username: req.Username, Email: req.Email, 
        PasswordHash: hashedPassword, FullName: req.FullName, 
        RoleID: role.ID, IsActive: true,
	}

	// 4. Simpan ke database
	createdUser, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, http.StatusCreated, nil
}

// UpdateUser
func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, req *models.UpdateUserRequest) (*models.User, int, error) {
	// 1. Cek User
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	// 2. Dapatkan RoleID baru jika RoleName diubah
	var newRoleID *uuid.UUID
	if req.RoleName != nil && *req.RoleName != "" {
		role, err := s.roleRepo.GetRoleByName(ctx, *req.RoleName)
		if err != nil || role == nil {
			return nil, http.StatusBadRequest, errors.New("invalid role name specified")
		}
		newRoleID = &role.ID
	}

	// 3. Update User
	updatedUser, err := s.userRepo.UpdateUser(ctx, userID, req, newRoleID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, http.StatusOK, nil
}

// DeleteUser (Soft Delete/Deactivate)
func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) (int, error) {
	// 1. Cek User
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}
	
	// 2. Nonaktifkan user
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to deactivate user: %w", err)
	}
	return http.StatusOK, nil
}

func (s *userService) SetStudentProfile(ctx context.Context, userID uuid.UUID, req *models.StudentProfileRequest) (*models.Student, int, error) {
	// Cek user exists
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	student := &models.Student{
		UserID:       userID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
	}

	updated, err := s.profileRepo.UpsertStudent(ctx, student)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return updated, http.StatusOK, nil
}

func (s *userService) SetLecturerProfile(ctx context.Context, userID uuid.UUID, req *models.LecturerProfileRequest) (*models.Lecturer, int, error) {
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	lecturer := &models.Lecturer{
		UserID:     userID,
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	updated, err := s.profileRepo.UpsertLecturer(ctx, lecturer)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return updated, http.StatusOK, nil
}

func (s *userService) AssignAdvisor(ctx context.Context, studentUserID uuid.UUID, advisorUserID uuid.UUID) (int, error) {
	// 1. Cari Profil Dosen berdasarkan User ID Dosen
	lecturer, err := s.profileRepo.GetLecturerByUserID(ctx, advisorUserID)
	if err != nil {
		return http.StatusNotFound, errors.New("advisor profile not found (user must have lecturer profile first)")
	}

	// 2. Update Student Profile dengan ID Dosen (Lecturer ID)
	err = s.profileRepo.AssignAdvisor(ctx, studentUserID, lecturer.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}