package controllers

import (
	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/services"
	"prestasi-mahasiswa-api/utils"

	"github.com/gofiber/fiber/v2"
)

type ReportController struct {
	Service services.ReportService
}

func NewReportController(service services.ReportService) *ReportController {
	return &ReportController{Service: service}
}

// GetDashboardStats godoc
// @Summary      Get Dashboard Statistics
// @Description  Melihat ringkasan data statistik prestasi (Total, by Status, by Type)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.JSONResponse
// @Failure      500  {object}  utils.JSONResponse
// @Router       /reports/statistics [get]
func (ctrl *ReportController) GetDashboardStats(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)

	stats, status, err := ctrl.Service.GetDashboardStats(c.Context(), claims.Role, claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}

	return utils.SuccessResponse(c, status, "Dashboard statistics retrieved", stats)
}