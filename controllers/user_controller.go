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