package tests

import (
	"context"
	"errors"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/services"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	args := m.Called(ctx, identifier)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}

// Implementasikan method interface lainnya (kosongkan saja)
func (m *MockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) { return nil, nil }
func (m *MockUserRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) { return nil, nil }
func (m *MockUserRepo) UpdateUser(ctx context.Context, u uuid.UUID, r *models.UpdateUserRequest, ri *uuid.UUID) (*models.User, error) { return nil, nil }
func (m *MockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockUserRepo) ListAllUsers(ctx context.Context) ([]models.User, error) { return nil, nil }
func (m *MockUserRepo) GetAdviseeStudentUserIDsByAdvisorUserID(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) { return nil, nil }

func TestLoginService(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := services.NewAuthService(mockRepo)

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("FindUserByUsernameOrEmail", mock.Anything, "unknown").Return(nil, errors.New("not found"))

		_, status, err := service.PerformLogin(context.Background(), "unknown", "password")
		assert.Error(t, err)
		assert.Equal(t, 401, status)
	})
}