package controllers

import (
	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	Service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{Service: service}
}

// Login Handler (FR-001)
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	resp, status, err := ctrl.Service.PerformLogin(c.Context(), req.Username, req.Password)
	if err != nil {
		[cite_start]// Error responses use the standard table in SRS [cite: 348]
		return c.Status(status).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(status).JSON(resp) // Status OK (200) atau Created (201)
}

// Profile Handler
func (ctrl *AuthController) GetProfile(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	resp, status, err := ctrl.Service.GetProfile(c.Context(), claims.UserID)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}