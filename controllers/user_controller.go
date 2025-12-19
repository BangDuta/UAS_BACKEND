package controllers

import (
	"fmt"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/services"
	"prestasi-mahasiswa-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{Service: service}
}

// ListAllUsers godoc
// @Summary      List All Users
// @Description  Admin melihat semua daftar pengguna (Mahasiswa, Dosen, Admin)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.JSONResponse
// @Failure      500  {object}  utils.JSONResponse
// @Router       /users [get]
func (ctrl *UserController) ListAllUsers(c *fiber.Ctx) error {
	users, status, err := ctrl.Service.ListAllUsers(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Users retrieved successfully", users)
}

// GetUserByID godoc
// @Summary      Get User By ID
// @Description  Admin melihat detail data pengguna berdasarkan UUID
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Success      200  {object}  utils.JSONResponse
// @Failure      404  {object}  utils.JSONResponse
// @Router       /users/{id} [get]
func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format")
	}

	user, status, err := ctrl.Service.GetUserByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "User retrieved successfully", user)
}

// CreateUser godoc
// @Summary      Create New User
// @Description  Admin membuat akun pengguna baru
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateUserRequest true "User Data"
// @Success      201  {object}  utils.JSONResponse
// @Failure      400  {object}  utils.JSONResponse
// @Router       /users [post]
func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload")
	}
	
	user, status, err := ctrl.Service.CreateUser(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "User created successfully", user)
}

// UpdateUser godoc
// @Summary      Update User
// @Description  Admin memperbarui data pengguna
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Param        request body models.UpdateUserRequest true "Update Data"
// @Success      200  {object}  utils.JSONResponse
// @Failure      400  {object}  utils.JSONResponse
// @Router       /users/{id} [put]
func (ctrl *UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format")
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	user, status, err := ctrl.Service.UpdateUser(c.Context(), id, &req)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "User updated successfully", user)
}

// DeleteUser godoc
// @Summary      Deactivate User (Soft Delete)
// @Description  Admin menonaktifkan akun pengguna
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Success      200  {object}  utils.JSONResponse
// @Failure      400  {object}  utils.JSONResponse
// @Router       /users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format")
	}

	status, err := ctrl.Service.DeleteUser(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, fmt.Sprintf("User %s has been deactivated", id.String()), nil)
}

// SetStudentProfile godoc
// @Summary      Set Student Profile
// @Description  Melengkapi data akademik mahasiswa (NIM, Prodi, Tahun)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Param        request body models.StudentProfileRequest true "Student Profile Data"
// @Success      200  {object}  utils.JSONResponse
// @Router       /users/{id}/student-profile [post]
func (ctrl *UserController) SetStudentProfile(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format")
	}

	var req models.StudentProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	res, status, err := ctrl.Service.SetStudentProfile(c.Context(), id, &req)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Student profile updated", res)
}

// SetLecturerProfile godoc
// @Summary      Set Lecturer Profile
// @Description  Melengkapi data kepegawaian dosen (NIP, Departemen)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Param        request body models.LecturerProfileRequest true "Lecturer Profile Data"
// @Success      200  {object}  utils.JSONResponse
// @Router       /users/{id}/lecturer-profile [post]
func (ctrl *UserController) SetLecturerProfile(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format")
	}

	var req models.LecturerProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	res, status, err := ctrl.Service.SetLecturerProfile(c.Context(), id, &req)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Lecturer profile updated", res)
}

// AssignAdvisor godoc
// @Summary      Assign Advisor to Student
// @Description  Menghubungkan mahasiswa dengan dosen walinya
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID Mahasiswa (UUID)"
// @Param        request body models.AssignAdvisorRequest true "User ID Dosen (UUID)"
// @Success      200  {object}  utils.JSONResponse
// @Router       /users/{id}/advisor [put]
func (ctrl *UserController) AssignAdvisor(c *fiber.Ctx) error {
	studentUserID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Student User ID")
	}

	var req models.AssignAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	advisorUUID, err := uuid.Parse(req.AdvisorUserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Advisor User ID")
	}

	status, err := ctrl.Service.AssignAdvisor(c.Context(), studentUserID, advisorUUID)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Advisor assigned successfully", nil)
}