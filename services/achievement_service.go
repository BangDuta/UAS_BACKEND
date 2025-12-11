package services

import (
	"context"
	"errors"
	"net/http"

	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/repositories"

	"github.com/google/uuid"
)

type AchievementService interface {
	CreateDraft(ctx context.Context, studentID uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error)
	DeleteDraft(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (int, error)
}

type achievementService struct {
	achieveRepo repositories.AchievementRepository
}

func NewAchievementService(achieveRepo repositories.AchievementRepository) AchievementService {
	return &achievementService{achieveRepo: achieveRepo}
}

// CreateDraft menangani pembuatan prestasi baru (FR-003)
func (s *achievementService) CreateDraft(ctx context.Context, studentID uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error) {
	// Validasi input di service layer sebelum masuk ke repository
	if req.Title == "" || req.AchievementType == "" {
		return nil, http.StatusBadRequest, errors.New("Title and AchievementType are required")
	}
	
	// Mapping request ke model MongoDB
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
		return nil, http.StatusInternalServerError, errors.New("Failed to create achievement draft")
	}

	return ref, http.StatusCreated, nil
}

// DeleteDraft menangani penghapusan prestasi berstatus 'draft' (FR-005)
func (s *achievementService) DeleteDraft(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (int, error) {
	err := s.achieveRepo.SoftDeleteAchievementAndReference(ctx, achievementRefID, studentID)
	if err != nil {
		if err.Error() == "achievement reference not found, status is not 'draft', or access denied" {
			return http.StatusForbidden, errors.New("Achievement not found or cannot be deleted (status not 'draft' or not owned by user)")
		}
		return http.StatusInternalServerError, errors.New("Failed to delete achievement")
	}

	return http.StatusOK, nil
}