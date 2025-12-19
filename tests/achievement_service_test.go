package tests

import (
	"context"
	"net/http"
	"testing"

	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/services"


	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK REPOSITORY ACHIEVEMENT ---
type MockAchieveRepo struct {
	mock.Mock
}

func (m *MockAchieveRepo) CreateAchievementAndReference(ctx context.Context, a *models.Achievement, id uuid.UUID) (*models.AchievementReference, error) {
	args := m.Called(ctx, a, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.AchievementReference), args.Error(1)
}

func (m *MockAchieveRepo) GetReferenceByID(ctx context.Context, refID uuid.UUID) (*models.AchievementReference, error) {
	args := m.Called(ctx, refID)
	if args.Get(0) == nil { return nil, args.Error(1) } // Mencegah nil pointer
	return args.Get(0).(*models.AchievementReference), args.Error(1)
}

func (m *MockAchieveRepo) UpdateAchievement(ctx context.Context, mid string, u interface{}) error {
	args := m.Called(ctx, mid, u)
	return args.Error(0)
}

func (m *MockAchieveRepo) UpdateReferenceUpdatedAt(ctx context.Context, rid uuid.UUID) (*models.AchievementReference, error) {
	args := m.Called(ctx, rid)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.AchievementReference), args.Error(1)
}

// Placeholder untuk method lain agar memenuhi interface AchievementRepository
func (m *MockAchieveRepo) SoftDeleteAchievementAndReference(ctx context.Context, aid uuid.UUID, sid uuid.UUID) error { return nil }
func (m *MockAchieveRepo) UpdateReferenceStatus(ctx context.Context, rid uuid.UUID, cs string, ns string, rn string, vb uuid.UUID) (*models.AchievementReference, error) { return nil, nil }
func (m *MockAchieveRepo) GetAchievementDetail(ctx context.Context, mid string) (*models.Achievement, error) { return nil, nil }
func (m *MockAchieveRepo) ListAchievementReferences(ctx context.Context, sids []uuid.UUID) ([]models.AchievementReference, error) { return nil, nil }
func (m *MockAchieveRepo) GetStatsByStatus(ctx context.Context, sid *uuid.UUID) (map[string]int, error) { return nil, nil }
func (m *MockAchieveRepo) GetStatsByType(ctx context.Context, sid *uuid.UUID) (map[string]int, error) { return nil, nil }
func (m *MockAchieveRepo) HardDeleteAchievement(ctx context.Context, rid uuid.UUID) error { return nil }

// --- MOCK REPOSITORY USER ---
type MockUserRepoForService struct {
	mock.Mock
}

func (m *MockUserRepoForService) GetAdviseeStudentUserIDsByAdvisorUserID(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// Placeholder UserRepo
func (m *MockUserRepoForService) FindUserByUsernameOrEmail(ctx context.Context, ident string) (*models.User, error) { return nil, nil }
func (m *MockUserRepoForService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) { return nil, nil }
func (m *MockUserRepoForService) CreateUser(ctx context.Context, u *models.User) (*models.User, error) { return nil, nil }
func (m *MockUserRepoForService) UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest, rid *uuid.UUID) (*models.User, error) { return nil, nil }
func (m *MockUserRepoForService) DeleteUser(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockUserRepoForService) ListAllUsers(ctx context.Context) ([]models.User, error) { return nil, nil }

// --- TEST CASES ---

func TestCreateAchievement(t *testing.T) {
	mockRepo := new(MockAchieveRepo)
	mockUser := new(MockUserRepoForService)
	service := services.NewAchievementService(mockRepo, mockUser)

	t.Run("Create Draft Success", func(t *testing.T) {
		studentID := uuid.New()
		req := &models.CreateAchievementRequest{Title: "Lomba Web"}
		expectedRef := &models.AchievementReference{Status: "draft"}

		mockRepo.On("CreateAchievementAndReference", mock.Anything, mock.Anything, studentID).Return(expectedRef, nil)

		res, status, err := service.CreateDraft(context.Background(), studentID, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
		assert.Equal(t, "draft", res.Status)
	})
}

func TestAddAttachment(t *testing.T) {
	mockRepo := new(MockAchieveRepo)
	mockUser := new(MockUserRepoForService)
	service := services.NewAchievementService(mockRepo, mockUser)

	t.Run("Add Attachment Success", func(t *testing.T) {
		studentID := uuid.New()
		refID := uuid.New()
		attachment := models.AttachmentFile{FileName: "bukti.jpg"}

		// Mocking GetReferenceByID mengembalikan data valid (bukan nil)
		mockRepo.On("GetReferenceByID", mock.Anything, refID).Return(&models.AchievementReference{
			ID:        refID,
			StudentID: studentID,
		}, nil)

		mockRepo.On("UpdateAchievement", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("UpdateReferenceUpdatedAt", mock.Anything, refID).Return(&models.AchievementReference{}, nil)

		status, err := service.AddAttachment(context.Background(), studentID, refID, attachment)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
	})

	t.Run("Add Attachment Forbidden - Wrong Student", func(t *testing.T) {
		studentID := uuid.New()
		otherStudentID := uuid.New()
		refID := uuid.New()

		mockRepo.On("GetReferenceByID", mock.Anything, refID).Return(&models.AchievementReference{
			StudentID: otherStudentID,
		}, nil)

		status, err := service.AddAttachment(context.Background(), studentID, refID, models.AttachmentFile{})

		assert.Error(t, err)
		assert.Equal(t, http.StatusForbidden, status)
		assert.Contains(t, err.Error(), "access denied")
	})
}