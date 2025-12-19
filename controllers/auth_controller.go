package controllers

import (
	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/services"
	"prestasi-mahasiswa-api/utils"

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
// @Success      200  {object}  utils.JSONResponse
// @Failure      401  {object}  utils.JSONResponse
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	resp, status, err := ctrl.Service.PerformLogin(c.Context(), req.Username, req.Password)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}

	return utils.SuccessResponse(c, status, "Login successful", resp)
}

// Logout godoc
// @Summary      Logout User
// @Description  Logout dari sistem (Client harus menghapus token di sisi mereka)
// @Tags         Auth
// @Security     BearerAuth
// @Success      200  {object}  utils.JSONResponse
// @Router       /auth/logout [post]
func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	// Stateless logout: Memberitahu client untuk membersihkan token
	return utils.SuccessResponse(c, fiber.StatusOK, "Logged out successfully", nil)
}

// GetProfile godoc
// @Summary      Get User Profile
// @Description  Mendapatkan data profil user yang sedang login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.JSONResponse
// @Failure      401  {object}  utils.JSONResponse
// @Router       /auth/profile [get]
func (ctrl *AuthController) GetProfile(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	resp, status, err := ctrl.Service.GetProfile(c.Context(), claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}

	return utils.SuccessResponse(c, status, "Profile retrieved successfully", resp)
}