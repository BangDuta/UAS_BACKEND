package controllers

import (
	"fmt"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{Service: service}
}

// List Users godoc
// @Summary      List All Users
// @Description  Admin melihat semua user yang terdaftar
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (ctrl *UserController) ListAllUsers(c *fiber.Ctx) error {
	users, status, err := ctrl.Service.ListAllUsers(c.Context())
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "data": users})
}

// Get User godoc
// @Summary      Get User By ID
// @Description  Admin melihat detail user tertentu
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid User ID format"})
	}

	user, status, err := ctrl.Service.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "data": user})
}

// Create User godoc
// @Summary      Create New User
// @Description  Admin membuat user baru (Mahasiswa/Dosen/Admin)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateUserRequest true "User Data"
// @Success      201  {object}  models.User
// @Failure      400  {object}  map[string]string
// @Router       /users [post]
func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request payload"})
	}
	
	user, status, err := ctrl.Service.CreateUser(c.Context(), &req)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "data": user})
}

// Update User godoc
// @Summary      Update User
// @Description  Admin mengubah data user
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Param        request body models.UpdateUserRequest true "Update Data"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]string
// @Router       /users/{id} [put]
func (ctrl *UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid User ID format"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request payload"})
	}

	user, status, err := ctrl.Service.UpdateUser(c.Context(), id, &req)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "data": user})
}

// Delete User godoc
// @Summary      Deactivate User
// @Description  Admin menonaktifkan user (Soft Delete)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid User ID format"})
	}

	status, err := ctrl.Service.DeleteUser(c.Context(), id)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("User %s has been deactivated", id.String())})
}

// Set Student Profile godoc
// @Summary      Set Student Profile
// @Description  Melengkapi data profil Mahasiswa (NIM, Prodi)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Param        request body models.StudentProfileRequest true "Student Profile"
// @Success      200  {object}  models.Student
// @Failure      400  {object}  map[string]string
// @Router       /users/{id}/student-profile [post]
func (ctrl *UserController) SetStudentProfile(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid User ID"})
	}

	var req models.StudentProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid body"})
	}

	res, status, err := ctrl.Service.SetStudentProfile(c.Context(), id, &req)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "data": res})
}

// Set Lecturer Profile godoc
// @Summary      Set Lecturer Profile
// @Description  Melengkapi data profil Dosen (NIP, Department)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Param        request body models.LecturerProfileRequest true "Lecturer Profile"
// @Success      200  {object}  models.Lecturer
// @Failure      400  {object}  map[string]string
// @Router       /users/{id}/lecturer-profile [post]
func (ctrl *UserController) SetLecturerProfile(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid User ID"})
	}

	var req models.LecturerProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid body"})
	}

	res, status, err := ctrl.Service.SetLecturerProfile(c.Context(), id, &req)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "data": res})
}

// Assign Advisor godoc
// @Summary      Assign Advisor to Student
// @Description  Menunjuk Dosen Wali untuk Mahasiswa
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID Mahasiswa (UUID)"
// @Param        request body models.AssignAdvisorRequest true "User ID Dosen"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /users/{id}/advisor [put]
func (ctrl *UserController) AssignAdvisor(c *fiber.Ctx) error {
	studentUserID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid Student User ID"})
	}

	var req models.AssignAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid body"})
	}

	advisorUUID, err := uuid.Parse(req.AdvisorUserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid Advisor User ID"})
	}

	status, err := ctrl.Service.AssignAdvisor(c.Context(), studentUserID, advisorUUID)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "message": "Advisor assigned successfully"})
}