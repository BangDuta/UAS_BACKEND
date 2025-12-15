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

// Create Achievement godoc
// @Summary      Create Achievement Draft
// @Description  Mahasiswa membuat draft prestasi baru
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateAchievementRequest true "Achievement Data"
// @Success      201  {object}  models.AchievementReference
// @Failure      400  {object}  map[string]string
// @Router       /achievements [post]
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

// Update Achievement godoc
// @Summary      Update Achievement Draft
// @Description  Mahasiswa mengedit data prestasi (hanya jika status draft)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Achievement ID (UUID)"
// @Param        request body models.CreateAchievementRequest true "Achievement Data"
// @Success      200  {object}  models.AchievementReference
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /achievements/{id} [put]
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

// Delete Achievement godoc
// @Summary      Delete Achievement
// @Description  Mahasiswa menghapus draft prestasi (Soft Delete)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Achievement ID (UUID)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /achievements/{id} [delete]
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

// Submit Achievement godoc
// @Summary      Submit Achievement
// @Description  Mahasiswa mengajukan prestasi draft untuk diverifikasi
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Achievement ID (UUID)"
// @Success      200  {object}  models.AchievementReference
// @Failure      400  {object}  map[string]string
// @Router       /achievements/{id}/submit [post]
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

// Verify Achievement godoc
// @Summary      Verify Achievement
// @Description  Dosen Wali menyetujui prestasi mahasiswa
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Achievement ID (UUID)"
// @Success      200  {object}  models.AchievementReference
// @Failure      400  {object}  map[string]string
// @Router       /achievements/{id}/verify [post]
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

// Reject Achievement godoc
// @Summary      Reject Achievement
// @Description  Dosen Wali menolak prestasi mahasiswa dengan catatan
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Achievement ID (UUID)"
// @Param        request body map[string]string true "Rejection Note (key: rejectionNote)"
// @Success      200  {object}  models.AchievementReference
// @Failure      400  {object}  map[string]string
// @Router       /achievements/{id}/reject [post]
func (ctrl *AchievementController) Reject(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid ID format"})
	}

	var req struct {
		RejectionNote string `json:"rejectionNote"`
	}
	if err := c.BodyParser(&req); err != nil {
		req.RejectionNote = "Rejected without specific note."
	}
	
	resp, status, err := ctrl.Service.RejectAchievement(c.Context(), claims.UserID, id, req.RejectionNote)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{"status": "success", "data": resp})
}

// List Achievements godoc
// @Summary      List Achievements
// @Description  Melihat daftar prestasi (Mahasiswa: milik sendiri, Dosen: milik bimbingan, Admin: semua)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.AchievementDetailResponse
// @Failure      400  {object}  map[string]string
// @Router       /achievements [get]
func (ctrl *AchievementController) List(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	
	resp, status, err := ctrl.Service.ListFilteredAchievements(c.Context(), claims)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	
	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data": resp,
	})
}

// Detail Achievement godoc
// @Summary      Get Achievement Detail
// @Description  Melihat detail lengkap prestasi
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Achievement ID (UUID)"
// @Success      200  {object}  models.AchievementDetailResponse
// @Failure      400  {object}  map[string]string
// @Router       /achievements/{id} [get]
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

// Upload Attachment godoc
// @Summary      Upload Attachment
// @Description  Upload file bukti prestasi (PDF/Image)
// @Tags         Achievements
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Achievement ID (UUID)"
// @Param        attachment formData file true "File Attachment"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /achievements/{id}/attachments [post]
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

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath := fmt.Sprintf("%s/%s-%s", uploadDir, id.String(), file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to save file"})
	}

	attachment := models.AttachmentFile{
		FileName: file.Filename,
		FileUrl:  filePath,
		FileType: file.Header.Get("Content-Type"),
	}

	status, err := ctrl.Service.AddAttachment(c.Context(), claims.UserID, id, attachment)
	if err != nil {
		os.Remove(filePath)
		return c.Status(status).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"status": "success", "message": "File uploaded and linked"})
}