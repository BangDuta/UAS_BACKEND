package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/repositories"
	"prestasi-mahasiswa-api/services"
	"prestasi-mahasiswa-api/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

// RouteManager menyimpan semua dependencies service
type RouteManager struct {
	AuthService        services.AuthService
	AchievementService services.AchievementService // Untuk Commit #6 dst.
	// ... services lainnya
}

// NewRouteManager melakukan Dependency Injection
func NewRouteManager(pgDB *pgxpool.Pool, mongoClient *mongo.Client) *RouteManager {
	// Repositories
	userRepo := repositories.NewUserRepository(pgDB)
	// achieveRepo := repositories.NewAchievementRepository(pgDB, mongoClient)

	return &RouteManager{
		AuthService: services.NewAuthService(userRepo),
		// AchievementService: services.NewAchievementService(achieveRepo),
	}
}

// SetupRoutes mendaftarkan semua endpoint API
func SetupRoutes(r *mux.Router, pgDB *pgxpool.Pool, mongoClient *mongo.Client) {
	rm := NewRouteManager(pgDB, mongoClient)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// 5.1 Authentication (FR-001)
	auth := v1.PathPrefix("/auth").Subrouter()
	
	// POST /auth/login
	auth.HandleFunc("/login", rm.Login).Methods("POST")
	
	// GET /auth/profile
	// Route ini dilindungi oleh AuthRequired middleware
	auth.HandleFunc("/profile", middleware.AuthRequired(rm.Profile)).Methods("GET")

	// ... Route Achievements akan ditambahkan di Commit #6
}

// --- Implementasi Route Methods ---

// Login menangani POST /api/v1/auth/login (FR-001)
func (rm *RouteManager) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Panggil Service Layer
	tokenResp, httpStatus, err := rm.AuthService.PerformLogin(r.Context(), req.Username, req.Password)
	if err != nil {
		utils.RespondWithError(w, httpStatus, err.Error())
		return
	}

	utils.RespondWithJSON(w, httpStatus, tokenResp)
}

// Profile menangani GET /api/v1/auth/profile
func (rm *RouteManager) Profile(w http.ResponseWriter, r *http.Request) {
	// Ambil claims dari context (telah dijamin ada oleh AuthRequired middleware)
	claims := middleware.GetUserClaims(r.Context()) 
	
	// Panggil service untuk mendapatkan detail profil terbaru
	profile, httpStatus, err := rm.AuthService.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		utils.RespondWithError(w, httpStatus, err.Error())
		return
	}

	resp := map[string]interface{}{
		"status": "success",
		"data": profile,
	}
	utils.RespondWithJSON(w, httpStatus, resp)
}

// ... di dalam func SetupRoutes ...

	rm := NewRouteManager(pgDB, mongoClient) // Update NewRouteManager agar menginjeksikan AchievementRepo

	// ... Auth Routes

	// 5.4 Achievements (Mahasiswa - FR-003, FR-005)
	achievements := v1.PathPrefix("/achievements").Subrouter()
	
	// POST /achievements (FR-003) - requires achievement:create
	achievements.HandleFunc("", middleware.RBACRequired("achievement:create", rm.CreateAchievement)).Methods("POST")
	
	// DELETE /achievements/:id (FR-005) - requires achievement:delete
	achievements.HandleFunc("/{id}", middleware.RBACRequired("achievement:delete", rm.DeleteAchievement)).Methods("DELETE")
	
	// ... Tambahkan route untuk UpdateAchievement di sini
}

// ... Tambahkan NewRouteManager agar menginjeksikan AchievementRepo
func NewRouteManager(pgDB *pgxpool.Pool, mongoClient *mongo.Client) *RouteManager {
	// Repositories
	userRepo := repositories.NewUserRepository(pgDB)
	achieveRepo := repositories.NewAchievementRepository(pgDB, mongoClient) // Injeksi AchievmentRepo

	return &RouteManager{
		AuthService: services.NewAuthService(userRepo),
		AchievementService: services.NewAchievementService(achieveRepo), // Injeksi AchievmentService
	}
}

// --- Tambahkan Implementasi Route Methods untuk Achievement ---

// CreateAchievement menangani POST /api/v1/achievements (FR-003)
func (rm *RouteManager) CreateAchievement(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	
	var req models.CreateAchievementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	studentID := claims.UserID // Asumsi claims.UserID adalah StudentID (Role Mahasiswa)
	
	ref, httpStatus, err := rm.AchievementService.CreateDraft(r.Context(), studentID, &req)
	if err != nil {
		utils.RespondWithError(w, httpStatus, err.Error())
		return
	}

	utils.RespondWithJSON(w, httpStatus, map[string]interface{}{
		"status": "success",
		"message": "Achievement draft created successfully",
		"data": ref,
	})
}

// DeleteAchievement menangani DELETE /api/v1/achievements/:id (FR-005)
func (rm *RouteManager) DeleteAchievement(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing achievement ID")
		return
	}
	
	achievementRefID, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid achievement ID format")
		return
	}

	studentID := claims.UserID
	
	httpStatus, err := rm.AchievementService.DeleteDraft(r.Context(), studentID, achievementRefID)
	if err != nil {
		utils.RespondWithError(w, httpStatus, err.Error())
		return
	}

	utils.RespondWithJSON(w, httpStatus, map[string]string{
		"status": "success",
		"message": "Achievement draft deleted successfully",
	})
}