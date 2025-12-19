package utils

import "github.com/gofiber/fiber/v2"

// Response data structure
type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse mengirim respon sukses (200, 201, dll)
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(JSONResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirim respon error (400, 401, 403, 404, 500)
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(JSONResponse{
		Status:  "error",
		Message: message,
	})
}