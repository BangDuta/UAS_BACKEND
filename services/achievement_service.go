package services

import (
	"context"
	"errors"
	"net/http"
	"time"
	"fmt"

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
	AddAttachment(ctx context.Context, studentID uuid.UUID, refID uuid.UUID, attachment models.AttachmentFile) (int, error) // Tambahan
	
	// Workflow
	SubmitForVerification(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (*models.AchievementReference, int, error)
	VerifyAchievement(ctx context.Context, advisorUserID uuid.UUID, achievementRefID uuid.UUID) (*models.AchievementReference, int, error)
	RejectAchievement(ctx context.Context, advisorUserID uuid.UUID, achievementRefID uuid.UUID, rejectionNote string) (*models.AchievementReference, int, error)
	
	// Read (FR-006, FR-010)
	ListFilteredAchievements(ctx context.Context, claims *utils.JWTCustomClaims) ([]models.AchievementDetailResponse, int, error)
	GetDetailWithVerification(ctx context.Context, claims *utils.JWTCustomClaims, refID uuid.UUID) (*models.AchievementDetailResponse, int, error)
	
	HardDelete(ctx context.Context, refID uuid.UUID) (int, error)
}

type achievementService struct {
	achieveRepo repositories.AchievementRepository
	userRepo repositories.UserRepository // Diperlukan untuk FR-006 (Dosen Wali)
}

func NewAchievementService(achieveRepo repositories.AchievementRepository, userRepo repositories.UserRepository) AchievementService {
	return &achievementService{achieveRepo: achieveRepo, userRepo: userRepo}
}

// CreateDraft (FR-003)
func (s *achievementService) CreateDraft(ctx context.Context, studentID uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error) {
	// Pastikan ID user adalah student ID
	// ... (Diabaikan untuk POC, diasumsi claims.UserID adalah StudentID)
	
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
		return nil, http.StatusInternalServerError, errors.New("failed to create achievement: " + err.Error())
	}
	return ref, http.StatusCreated, nil
}

// DeleteDraft (FR-005)
func (s *achievementService) DeleteDraft(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (int, error) {
	err := s.achieveRepo.SoftDeleteAchievementAndReference(ctx, achievementRefID, studentID)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

// UpdateDraft (FR-003 Update)
func (s *achievementService) UpdateDraft(ctx context.Context, studentID uuid.UUID, refID uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error) {
	// 1. Cek Reference status: harus 'draft' dan milik studentID
	ref, err := s.achieveRepo.GetReferenceByID(ctx, refID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("achievement not found")
	}
	if ref.Status != "draft" || ref.StudentID != studentID {
		return nil, http.StatusForbidden, errors.New("only draft achievements can be updated by the owner")
	}

	// 2. Update Mongo data
	update := bson.M{"$set": bson.M{
		"achievementType": req.AchievementType,
		"title":           req.Title,
		"description":     req.Description,
		"details":         req.Details,
		"tags":            req.Tags,
		"points":          req.Points,
		"updatedAt":       time.Now(),
	}}

	err = s.achieveRepo.UpdateAchievement(ctx, ref.MongoAchievementID, update)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to update achievement detail")
	}
	
	// 3. Update PG reference (updated_at)
	_, err = s.achieveRepo.UpdateReferenceUpdatedAt(ctx, refID)
	if err != nil {
		// Log error, tapi tidak perlu mengembalikan 500 karena data sudah terupdate di Mongo
	}
	
	return ref, http.StatusOK, nil
}

// AddAttachment
func (s *achievementService) AddAttachment(ctx context.Context, studentUserID uuid.UUID, refID uuid.UUID, attachment models.AttachmentFile) (int, error) {
	// 1. Ambil data referensi
	ref, err := s.achieveRepo.GetReferenceByID(ctx, refID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("achievement reference not found: %w", err)
	}

	// SAFETY CHECK: Pastikan ref tidak nil sebelum akses ref.StudentID (Baris 116)
	if ref == nil {
		return http.StatusNotFound, errors.New("achievement data is empty")
	}

	// 2. Validasi kepemilikan (Penyebab panic jika ref nil)
	if ref.StudentID != studentUserID {
		return http.StatusForbidden, errors.New("access denied: this is not your achievement")
	}

	// ... sisa kode ...
    return http.StatusCreated, nil
}

// SubmitForVerification (FR-004)
func (s *achievementService) SubmitForVerification(ctx context.Context, studentID uuid.UUID, achievementRefID uuid.UUID) (*models.AchievementReference, int, error) {
	// Status check (Optional: verify studentID ownership first)
	ref, err := s.achieveRepo.GetReferenceByID(ctx, achievementRefID)
	if err != nil || ref.StudentID != studentID {
		return nil, http.StatusNotFound, errors.New("achievement not found or forbidden")
	}
	
	if ref.Status != "draft" {
		return nil, http.StatusConflict, errors.New("achievement status must be 'draft' to be submitted")
	}
	
	ref, err = s.achieveRepo.UpdateReferenceStatus(ctx, achievementRefID, "draft", "submitted", "", uuid.Nil)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to update status: " + err.Error())
	}
	// TODO: Add notification logic (SRS P. 9, item 3)
	return ref, http.StatusOK, nil
}

// VerifyAchievement (FR-007)
func (s *achievementService) VerifyAchievement(ctx context.Context, advisorUserID uuid.UUID, achievementRefID uuid.UUID) (*models.AchievementReference, int, error) {
	// Validation: Check if the achievement belongs to an advisee of this advisorUserID (FR-006 prerequisite)
	
	// 1. Ambil reference
	ref, err := s.achieveRepo.GetReferenceByID(ctx, achievementRefID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("achievement not found")
	}
	
	// 2. Cek apakah advisorUserID adalah dosen wali dari ref.StudentID
	// TODO: Get student advisorID based on ref.StudentID and compare with advisorUserID
	
	if ref.Status != "submitted" {
		return nil, http.StatusConflict, errors.New("achievement status must be 'submitted' to be verified")
	}
	
	ref, err = s.achieveRepo.UpdateReferenceStatus(ctx, achievementRefID, "submitted", "verified", "", advisorUserID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to update status: " + err.Error())
	}
	
	// Set verified_at in PG is handled inside UpdateReferenceStatus
	return ref, http.StatusOK, nil
}

// RejectAchievement (FR-008)
func (s *achievementService) RejectAchievement(ctx context.Context, advisorUserID uuid.UUID, achievementRefID uuid.UUID, rejectionNote string) (*models.AchievementReference, int, error) {
	// Validation: Check if the achievement belongs to an advisee of this advisorUserID (FR-006 prerequisite)
	
	ref, err := s.achieveRepo.GetReferenceByID(ctx, achievementRefID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("achievement not found")
	}
	
	// Cek dosen wali ownership
	// TODO: Check Dosen Wali ownership here
	
	if ref.Status != "submitted" {
		return nil, http.StatusConflict, errors.New("achievement status must be 'submitted' to be rejected")
	}
	
	ref, err = s.achieveRepo.UpdateReferenceStatus(ctx, achievementRefID, "submitted", "rejected", rejectionNote, advisorUserID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to update status: " + err.Error())
	}
	
	// TODO: Add notification logic (SRS P. 10, item 4)
	return ref, http.StatusOK, nil
}

// ListFilteredAchievements (FR-006, FR-010)
func (s *achievementService) ListFilteredAchievements(ctx context.Context, claims *utils.JWTCustomClaims) ([]models.AchievementDetailResponse, int, error) {
	var filterStudentIDs []uuid.UUID
	
	// Tentukan filter berdasarkan Role
	switch claims.Role {
	case "Mahasiswa":
		filterStudentIDs = append(filterStudentIDs, claims.UserID)
	case "Dosen Wali":
		// FR-006: Get list student IDs dari tabel students where advisor id
		adviseeIDs, err := s.userRepo.GetAdviseeStudentUserIDsByAdvisorUserID(ctx, claims.UserID)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New("failed to retrieve advisees")
		}
		filterStudentIDs = adviseeIDs
	case "Admin":
		// FR-010: No student ID filter (ambil semua)
		filterStudentIDs = nil // nil means fetch all
	default:
		return nil, http.StatusForbidden, errors.New("user role not authorized to view achievements")
	}
	
	// 1. Get all relevant references from PG
	references, err := s.achieveRepo.ListAchievementReferences(ctx, filterStudentIDs)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to list references: " + err.Error())
	}
	
	if len(references) == 0 {
		return []models.AchievementDetailResponse{}, http.StatusOK, nil
	}
	
	var finalResponse []models.AchievementDetailResponse
	
	// 2. Fetch details from MongoDB for each reference and combine
	for _, ref := range references {
		detail, err := s.achieveRepo.GetAchievementDetail(ctx, ref.MongoAchievementID)
		if err != nil {
			// Jika Mongo detail tidak ditemukan, log error dan skip/lanjutkan
			continue
		}
		
		response := models.AchievementDetailResponse{
			RefID:           ref.ID,
			Status:          ref.Status,
			AchievementType: detail.AchievementType,
			Title:           detail.Title,
			Description:     detail.Description,
			Details:         detail.Details,
			Points:          detail.Points,
			RejectionNote:   ref.RejectionNote,
			SubmittedAt:     ref.SubmittedAt,
			VerifiedAt:      ref.VerifiedAt,
		}
		finalResponse = append(finalResponse, response)
	}
	
	return finalResponse, http.StatusOK, nil
}

// GetDetailWithVerification (Read)
func (s *achievementService) GetDetailWithVerification(ctx context.Context, claims *utils.JWTCustomClaims, refID uuid.UUID) (*models.AchievementDetailResponse, int, error) {
	// 1. Get reference
	ref, err := s.achieveRepo.GetReferenceByID(ctx, refID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("achievement not found")
	}
	
	// 2. Authorization Check (Must be owner, advisee, or Admin)
	isAuthorized := false
	if claims.Role == "Admin" {
		isAuthorized = true
	} else if claims.UserID == ref.StudentID {
		isAuthorized = true
	} else if claims.Role == "Dosen Wali" {
		// TODO: Check if ref.StudentID is an advisee of claims.UserID (FR-006 prerequisite)
		isAuthorized = true // Assume true for now
	}
	
	if !isAuthorized {
		return nil, http.StatusForbidden, errors.New("not authorized to view this achievement")
	}
	
	// 3. Get Mongo detail
	detail, err := s.achieveRepo.GetAchievementDetail(ctx, ref.MongoAchievementID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve achievement details")
	}
	
	// 4. Combine and return
	response := &models.AchievementDetailResponse{
		RefID:           ref.ID,
		Status:          ref.Status,
		AchievementType: detail.AchievementType,
		Title:           detail.Title,
		Description:     detail.Description,
		Details:         detail.Details,
		Points:          detail.Points,
		RejectionNote:   ref.RejectionNote,
		SubmittedAt:     ref.SubmittedAt,
		VerifiedAt:      ref.VerifiedAt,
	}
	
	return response, http.StatusOK, nil
}

func (s *achievementService) HardDelete(ctx context.Context, refID uuid.UUID) (int, error) {
	err := s.achieveRepo.HardDeleteAchievement(ctx, refID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}