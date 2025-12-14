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

// SetupRoutes: Hanya berisi definisi endpoint (Tidak ada logic handler/func)
func SetupRoutes(app *fiber.App, pgDB *pgxpool.Pool, mongoClient *mongo.Client) {
	// 1. Dependency Injection (Layering: Repo -> Service -> Controller)
	
	// Repositories
	userRepo := repositories.NewUserRepository(pgDB)
	achieveRepo := repositories.NewAchievementRepository(pgDB, mongoClient)
	roleRepo := repositories.NewRoleRepository(pgDB) 
	profileRepo := repositories.NewProfileRepository(pgDB)
	
	// Services
	authService := services.NewAuthService(userRepo)
	achieveService := services.NewAchievementService(achieveRepo, userRepo)
	userService := services.NewUserService(userRepo, roleRepo, profileRepo) // NEW: User Service
	reportService := services.NewReportService(achieveRepo)


	// Controllers
	authController := controllers.NewAuthController(authService)
	achieveController := controllers.NewAchievementController(achieveService)
	userController := controllers.NewUserController(userService)
	reportController := controllers.NewReportController(reportService) // NEW: User Controller

	// 2. Grouping Routes
	api := app.Group("/api/v1")

	// --- Auth Routes ---
	auth := api.Group("/auth")
	auth.Post("/login", authController.Login)
	auth.Get("/profile", middleware.AuthRequired, authController.GetProfile)

	// --- Achievements Routes ---
	ach := api.Group("/achievements", middleware.AuthRequired)
	
	// Mahasiswa Actions
	ach.Post("/", middleware.RBACRequired("achievement:create"), achieveController.Create)
	ach.Put("/:id", middleware.RBACRequired("achievement:update"), achieveController.Update)
	ach.Delete("/:id", middleware.RBACRequired("achievement:delete"), achieveController.Delete)
	ach.Post("/:id/submit", middleware.RBACRequired("achievement:update"), achieveController.Submit) 
	ach.Post("/:id/attachments", middleware.RBACRequired("achievement:update"), achieveController.UploadAttachment)
	
	// Dosen Wali/Admin Actions (Read & Workflow)
	ach.Get("/", achieveController.List) 
	ach.Get("/:id", achieveController.Detail)
	ach.Post("/:id/verify", middleware.RBACRequired("achievement:verify"), achieveController.Verify)
	ach.Post("/:id/reject", middleware.RBACRequired("achievement:verify"), achieveController.Reject)
	
	
	users := api.Group("/users", middleware.AuthRequired, middleware.RBACRequired("user:manage"))
	
	users.Get("/", userController.ListAllUsers)      // GET /api/v1/users
	users.Post("/", userController.CreateUser)       // POST /api/v1/users
	users.Get("/:id", userController.GetUserByID)  // GET /api/v1/users/:id
	users.Put("/:id", userController.UpdateUser)   // PUT /api/v1/users/:id
	users.Delete("/:id", userController.DeleteUser)// DELETE /api/v1/users/:id (Deactivate)

	users.Post("/:id/student-profile", userController.SetStudentProfile)
	users.Post("/:id/lecturer-profile", userController.SetLecturerProfile)
	users.Put("/:id/advisor", userController.AssignAdvisor) // Set Dosen Wali untuk Mahasiswa

	reports := api.Group("/reports", middleware.AuthRequired)
	reports.Get("/statistics", reportController.GetDashboardStats)
}