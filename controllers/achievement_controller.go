package controllers

import (
	"fmt"
	"os"
	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/services"
	"prestasi-mahasiswa-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementController struct {
	Service services.AchievementService
}

func NewAchievementController(service services.AchievementService) *AchievementController {
	return &AchievementController{Service: service}
}

// Create godoc
// @Summary      Create Achievement Draft
// @Tags         Achievements
// @Security     BearerAuth
// @Param        request body models.CreateAchievementRequest true "Achievement Data"
// @Success      201  {object}  utils.JSONResponse
// @Router       /achievements [post]
func (ctrl *AchievementController) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	resp, status, err := ctrl.Service.CreateDraft(c.Context(), claims.UserID, &req)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievement draft created", resp)
}

// Update godoc
// @Summary      Update Achievement Draft
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id path string true "Achievement ID"
// @Router       /achievements/{id} [put]
func (ctrl *AchievementController) Update(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}
	
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	resp, status, err := ctrl.Service.UpdateDraft(c.Context(), claims.UserID, id, &req)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievement updated", resp)
}

// Delete (Soft) godoc
// @Summary      Soft Delete Achievement
// @Tags         Achievements
// @Security     BearerAuth
// @Router       /achievements/{id} [delete]
func (ctrl *AchievementController) Delete(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	status, err := ctrl.Service.DeleteDraft(c.Context(), claims.UserID, id)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievement moved to trash", nil)
}

// HardDelete godoc
// @Summary      Hard Delete Achievement
// @Description  Hapus permanen dari PostgreSQL & MongoDB (Admin Only)
// @Tags         Achievements
// @Security     BearerAuth
// @Router       /achievements/{id}/hard [delete]
func (ctrl *AchievementController) HardDelete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	status, err := ctrl.Service.HardDelete(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievement permanently deleted", nil)
}

// Submit godoc
// @Summary      Submit Achievement
// @Tags         Achievements
// @Security     BearerAuth
// @Router       /achievements/{id}/submit [post]
func (ctrl *AchievementController) Submit(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	resp, status, err := ctrl.Service.SubmitForVerification(c.Context(), claims.UserID, id)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievement submitted for verification", resp)
}

// Verify godoc
// @Summary      Verify Achievement (Dosen Wali)
// @Tags         Achievements
// @Security     BearerAuth
// @Router       /achievements/{id}/verify [post]
func (ctrl *AchievementController) Verify(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	resp, status, err := ctrl.Service.VerifyAchievement(c.Context(), claims.UserID, id)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievement verified", resp)
}

// Reject godoc
// @Summary      Reject Achievement (Dosen Wali)
// @Tags         Achievements
// @Security     BearerAuth
// @Router       /achievements/{id}/reject [post]
func (ctrl *AchievementController) Reject(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		RejectionNote string `json:"rejectionNote"`
	}
	_ = c.BodyParser(&req) // Ignore error, use default if empty
	
	resp, status, err := ctrl.Service.RejectAchievement(c.Context(), claims.UserID, id, req.RejectionNote)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievement rejected", resp)
}

// List godoc
// @Summary      List Achievements
// @Tags         Achievements
// @Security     BearerAuth
// @Router       /achievements [get]
func (ctrl *AchievementController) List(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	resp, status, err := ctrl.Service.ListFilteredAchievements(c.Context(), claims)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievements retrieved", resp)
}

// Detail godoc
// @Summary      Get Achievement Detail
// @Tags         Achievements
// @Security     BearerAuth
// @Router       /achievements/{id} [get]
func (ctrl *AchievementController) Detail(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}
	
	resp, status, err := ctrl.Service.GetDetailWithVerification(c.Context(), claims, id)
	if err != nil {
		return utils.ErrorResponse(c, status, err.Error())
	}
	return utils.SuccessResponse(c, status, "Achievement detail retrieved", resp)
}

// UploadAttachment godoc
// @Summary      Upload Attachment
// @Tags         Achievements
// @Accept       mpfd
// @Security     BearerAuth
// @Param        attachment formData file true "File"
// @Router       /achievements/{id}/attachments [post]
func (ctrl *AchievementController) UploadAttachment(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	file, err := c.FormFile("attachment")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "File is required")
	}

	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	filePath := fmt.Sprintf("%s/%s-%s", uploadDir, id.String(), file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save file")
	}

	attachment := models.AttachmentFile{
		FileName: file.Filename,
		FileUrl:  filePath,
		FileType: file.Header.Get("Content-Type"),
	}

	status, err := ctrl.Service.AddAttachment(c.Context(), claims.UserID, id, attachment)
	if err != nil {
		os.Remove(filePath)
		return utils.ErrorResponse(c, status, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "File uploaded successfully", nil)
}