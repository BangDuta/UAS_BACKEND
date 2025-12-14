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

// GetDashboardStats Handler
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