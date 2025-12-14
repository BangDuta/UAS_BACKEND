package services

import (
	"context"
	"net/http"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/repositories"

	"github.com/google/uuid"
)

type ReportService interface {
	GetDashboardStats(ctx context.Context, userRole string, userID uuid.UUID) (*models.DashboardStats, int, error)
}

type reportService struct {
	achieveRepo repositories.AchievementRepository
}

func NewReportService(achieveRepo repositories.AchievementRepository) ReportService {
	return &reportService{achieveRepo: achieveRepo}
}

func (s *reportService) GetDashboardStats(ctx context.Context, userRole string, userID uuid.UUID) (*models.DashboardStats, int, error) {
	var filterStudentID *uuid.UUID

	// Jika Mahasiswa, hanya lihat data sendiri
	if userRole == "Mahasiswa" {
		filterStudentID = &userID
	}
	// Jika Admin/Dosen Wali, lihat semua (untuk penyederhanaan FR-011 saat ini)
	// Note: Dosen Wali idealnya difilter by advisee, tapi logic grouping-nya kompleks,
	// kita implementasi global stats dulu untuk Dosen/Admin.

	// 1. Ambil Stats Status (PostgreSQL)
	statusStats, err := s.achieveRepo.GetStatsByStatus(ctx, filterStudentID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// 2. Ambil Stats Tipe (MongoDB)
	typeStats, err := s.achieveRepo.GetStatsByType(ctx, filterStudentID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// 3. Hitung Total
	total := 0
	for _, count := range statusStats {
		total += count
	}

	stats := &models.DashboardStats{
		TotalAchievements: total,
		ByStatus:          statusStats,
		ByType:            typeStats,
	}

	return stats, http.StatusOK, nil
}