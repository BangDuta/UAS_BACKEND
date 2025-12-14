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

func (ctrl *UserController) ListAllUsers(c *fiber.Ctx) error {
	users, status, err := ctrl.Service.ListAllUsers(c.Context())
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"status": "success", "data": users})
}

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