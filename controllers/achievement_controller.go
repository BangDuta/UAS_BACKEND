package controllers

import (
	"fmt"
	"net/http"
	"os"
	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementController struct {
	Service services.AchievementService
}

func NewAchievementController(service services.AchievementService) *AchievementController {
	return &AchievementController{Service: service}
}

// Create Achievement (FR-003)
func (ctrl *AchievementController) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	var req models.CreateAchievementRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	resp, status, err := ctrl.Service.CreateDraft(c.Context(), claims.UserID, &req)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}

// Update Achievement Draft
func (ctrl *AchievementController) Update(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid ID format"})
	}
	
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	resp, status, err := ctrl.Service.UpdateDraft(c.Context(), claims.UserID, id, &req)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}

// Delete Achievement (FR-005)
func (ctrl *AchievementController) Delete(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid ID format"})
	}

	status, err := ctrl.Service.DeleteDraft(c.Context(), claims.UserID, id)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{"status": "success", "message": "Achievement deleted"})
}

// Submit for Verification (FR-004)
func (ctrl *AchievementController) Submit(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid ID format"})
	}

	resp, status, err := ctrl.Service.SubmitForVerification(c.Context(), claims.UserID, id)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{"status": "success", "data": resp})
}

// Verify Achievement (FR-007)
func (ctrl *AchievementController) Verify(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid ID format"})
	}

	resp, status, err := ctrl.Service.VerifyAchievement(c.Context(), claims.UserID, id)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{"status": "success", "data": resp})
}

// Reject Achievement (FR-008)
func (ctrl *AchievementController) Reject(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid ID format"})
	}

	var req struct {
		RejectionNote string `json:"rejectionNote"`
	}
	// Mengambil rejection note dari body
	if err := c.BodyParser(&req); err != nil {
		// Jika body kosong atau tidak valid, gunakan default message
		req.RejectionNote = "Rejected without specific note."
	}
	
	resp, status, err := ctrl.Service.RejectAchievement(c.Context(), claims.UserID, id, req.RejectionNote)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{"status": "success", "data": resp})
}

// List Achievements (FR-006, FR-010)
func (ctrl *AchievementController) List(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	
	// TODO: Add pagination parameters from query string here
	
	resp, status, err := ctrl.Service.ListFilteredAchievements(c.Context(), claims)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	
	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data": resp,
		// TODO: Add pagination info here
	})
}

// Detail Achievement
func (ctrl *AchievementController) Detail(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid ID format"})
	}
	
	resp, status, err := ctrl.Service.GetDetailWithVerification(c.Context(), claims, id)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	
	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data": resp,
	})
}

// Upload Attachment (Tambahan)
func (ctrl *AchievementController) UploadAttachment(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid ID format"})
	}

	file, err := c.FormFile("attachment")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Attachment file is required"})
	}

	// Pastikan direktori 'uploads' ada
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Simpan file ke server
	filePath := fmt.Sprintf("%s/%s-%s", uploadDir, id.String(), file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to save file"})
	}

	// Panggil Service untuk mencatat attachment di MongoDB
	attachment := models.AttachmentFile{
		FileName: file.Filename,
		FileUrl:  filePath, // Di production, ini harusnya S3/Blob URL
		FileType: file.Header.Get("Content-Type"),
		// UploadedAt di-set di service
	}

	status, err := ctrl.Service.AddAttachment(c.Context(), claims.UserID, id, attachment)
	if err != nil {
		// Jika gagal update Mongo, lakukan cleanup file yang baru diupload
		os.Remove(filePath)
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"status": "success", "message": "File uploaded and linked"})
}