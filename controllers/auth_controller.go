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

// Login godoc
// @Summary      Login User
// @Description  Masuk menggunakan username dan password untuk mendapatkan token JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Credentials"
// @Success      200  {object}  models.LoginResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
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
		return c.Status(status).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(status).JSON(resp)
}

// GetProfile godoc
// @Summary      Get User Profile
// @Description  Mendapatkan data user yang sedang login berdasarkan Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.UserProfile
// @Failure      401  {object}  map[string]string
// @Router       /auth/profile [get]
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