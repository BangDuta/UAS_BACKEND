package controllers

import (
	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/services"

	"github.com/gofiber/fiber/v2"
)

type ReportController struct {
	Service services.ReportService
}

func NewReportController(service services.ReportService) *ReportController {
	return &ReportController{Service: service}
}

// Get Dashboard Stats godoc
// @Summary      Get Dashboard Statistics
// @Description  Melihat statistik prestasi (Total, Status, Tipe)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.DashboardStats
// @Failure      500  {object}  map[string]string
// @Router       /reports/statistics [get]
func (ctrl *ReportController) GetDashboardStats(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)

	stats, status, err := ctrl.Service.GetDashboardStats(c.Context(), claims.Role, claims.UserID)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}