package services

import (
	"context"
	"net/http"

	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/repositories"
	"prestasi-mahasiswa-api/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type AchievementService interface {
	CreateDraft(ctx context.Context, studentID uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error)
	DeleteDraft(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (int, error)
	UpdateDraft(ctx context.Context, studentID uuid.UUID, refID uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error)
	
	// Workflow
	SubmitForVerification(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (*models.AchievementReference, int, error)
	VerifyAchievement(ctx context.Context, advisorUserID uuid.UUID, achievementRefID uuid.UUID) (*models.AchievementReference, int, error)
	RejectAchievement(ctx context.Context, advisorUserID uuid.UUID, achievementRefID uuid.UUID, rejectionNote string) (*models.AchievementReference, int, error)
	
	// Read
	ListFilteredAchievements(ctx context.Context, claims *utils.Claims) ([]models.AchievementDetailResponse, int, error)
	GetDetailWithVerification(ctx context.Context, claims *utils.Claims, refID uuid.UUID) (*models.AchievementDetailResponse, int, error)
}

type achievementService struct {
	achieveRepo repositories.AchievementRepository
}

func NewAchievementService(achieveRepo repositories.AchievementRepository) AchievementService {
	return &achievementService{achieveRepo: achieveRepo}
}

// CreateDraft
func (s *achievementService) CreateDraft(ctx context.Context, studentID uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error) {
	achievement := &models.Achievement{
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
	}

	ref, err := s.achieveRepo.CreateAchievementAndReference(ctx, achievement, studentID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return ref, http.StatusCreated, nil
}

// DeleteDraft
func (s *achievementService) DeleteDraft(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (int, error) {
	err := s.achieveRepo.SoftDeleteAchievementAndReference(ctx, achievementRefID, studentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// UpdateDraft (Placeholder Clean)
func (s *achievementService) UpdateDraft(ctx context.Context, studentID uuid.UUID, refID uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error) {
	// Variabel didefinisikan tapi di-ignore (_) agar compiler tidak error "unused variable"
	_ = bson.M{
		"achievementType": req.AchievementType,
		"title":           req.Title,
	}
	
	// Implementasi update sesungguhnya belum ada
	return nil, http.StatusOK, nil
}

// SubmitForVerification
func (s *achievementService) SubmitForVerification(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (*models.AchievementReference, int, error) {
	ref, err := s.achieveRepo.UpdateReferenceStatus(ctx, achievementRefID, "draft", "submitted", "", uuid.Nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return ref, http.StatusOK, nil
}

// VerifyAchievement
func (s *achievementService) VerifyAchievement(ctx context.Context, advisorUserID uuid.UUID, achievementRefID uuid.UUID) (*models.AchievementReference, int, error) {
	ref, err := s.achieveRepo.UpdateReferenceStatus(ctx, achievementRefID, "submitted", "verified", "", advisorUserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return ref, http.StatusOK, nil
}

// RejectAchievement
func (s *achievementService) RejectAchievement(ctx context.Context, advisorUserID uuid.UUID, achievementRefID uuid.UUID, rejectionNote string) (*models.AchievementReference, int, error) {
	ref, err := s.achieveRepo.UpdateReferenceStatus(ctx, achievementRefID, "submitted", "rejected", rejectionNote, advisorUserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return ref, http.StatusOK, nil
}

// ListFilteredAchievements (Placeholder Clean)
func (s *achievementService) ListFilteredAchievements(ctx context.Context, claims *utils.Claims) ([]models.AchievementDetailResponse, int, error) {
	// Hapus deklarasi variabel unused (filterUserID, filterRole)
	return []models.AchievementDetailResponse{}, http.StatusOK, nil
}

// GetDetailWithVerification (Placeholder Clean)
func (s *achievementService) GetDetailWithVerification(ctx context.Context, claims *utils.Claims, refID uuid.UUID) (*models.AchievementDetailResponse, int, error) {
	return &models.AchievementDetailResponse{}, http.StatusOK, nil
}