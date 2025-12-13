package routes

import (
	"prestasi-mahasiswa-api/controllers"
	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/repositories"
	"prestasi-mahasiswa-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes: Hanya berisi definisi endpoint (tidak ada logic handler)
func SetupRoutes(app *fiber.App, pgDB *pgxpool.Pool, mongoClient *mongo.Client) {
	// --- Dependency Injection ---
	userRepo := repositories.NewUserRepository(pgDB)
	achieveRepo := repositories.NewAchievementRepository(pgDB, mongoClient)

	authService := services.NewAuthService(userRepo)
	achieveService := services.NewAchievementService(achieveRepo, userRepo) // Inject userRepo ke AchieveService

	authController := controllers.NewAuthController(authService)
	achieveController := controllers.NewAchievementController(achieveService)

	// --- Grouping Routes ---
	api := app.Group("/api/v1")

	// 1. Auth (FR-001)
	auth := api.Group("/auth")
	auth.Post("/login", authController.Login)
	auth.Get("/profile", middleware.AuthRequired, authController.GetProfile) // FR-002 applied

	// 2. Achievements
	ach := api.Group("/achievements", middleware.AuthRequired) // Semua butuh Auth

	// Mahasiswa Actions (FR-003, FR-005)
	ach.Post("/", middleware.RBACRequired("achievement:create"), achieveController.Create)
	ach.Put("/:id", middleware.RBACRequired("achievement:update"), achieveController.Update) // Update Draft
	ach.Delete("/:id", middleware.RBACRequired("achievement:delete"), achieveController.Delete)
	ach.Post("/:id/submit", middleware.RBACRequired("achievement:update"), achieveController.Submit) // FR-004
	
	// Dosen Wali/Admin Read & Workflow (FR-006, FR-010, FR-007, FR-008)
	ach.Get("/", achieveController.List) // Filtering logic is inside the service
	ach.Get("/:id", achieveController.Detail)
	ach.Post("/:id/verify", middleware.RBACRequired("achievement:verify"), achieveController.Verify)
	ach.Post("/:id/reject", middleware.RBACRequired("achievement:verify"), achieveController.Reject)
	
	// Upload Attachments
	ach.Post("/:id/attachments", middleware.RBACRequired("achievement:update"), achieveController.UploadAttachment)
}