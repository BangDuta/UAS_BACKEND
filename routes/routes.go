package routes

import (
	"encoding/json"
	"net/http"

	"prestasi-mahasiswa-api/middleware"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/repositories"
	"prestasi-mahasiswa-api/services"
	"prestasi-mahasiswa-api/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

// RouteManager Struct
type RouteManager struct {
	AuthService        services.AuthService
	AchievementService services.AchievementService
}

// Konstruktor RouteManager
func NewRouteManager(pgDB *pgxpool.Pool, mongoClient *mongo.Client) *RouteManager {
	userRepo := repositories.NewUserRepository(pgDB)
	achieveRepo := repositories.NewAchievementRepository(pgDB, mongoClient)

	return &RouteManager{
		AuthService:        services.NewAuthService(userRepo),
		AchievementService: services.NewAchievementService(achieveRepo),
	}
}

// SetupRoutes
func SetupRoutes(r *mux.Router, pgDB *pgxpool.Pool, mongoClient *mongo.Client) {
	rm := NewRouteManager(pgDB, mongoClient)
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Auth
	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", rm.Login).Methods("POST")
	auth.HandleFunc("/profile", middleware.AuthRequired(rm.Profile)).Methods("GET")

	// Achievements
	ach := v1.PathPrefix("/achievements").Subrouter()
	ach.HandleFunc("", middleware.RBACRequired("achievement:create", rm.CreateAchievement)).Methods("POST")
	ach.HandleFunc("/{id}", middleware.RBACRequired("achievement:delete", rm.DeleteAchievement)).Methods("DELETE")
	ach.HandleFunc("/{id}/submit", middleware.RBACRequired("achievement:update", rm.SubmitAchievement)).Methods("POST")
	ach.HandleFunc("/{id}/verify", middleware.RBACRequired("achievement:verify", rm.VerifyAchievement)).Methods("POST")
	ach.HandleFunc("/{id}/reject", middleware.RBACRequired("achievement:verify", rm.RejectAchievement)).Methods("POST")
}

// --- Handler Methods ---

func (rm *RouteManager) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, 400, "Invalid Request")
		return
	}
	resp, status, err := rm.AuthService.PerformLogin(r.Context(), req.Username, req.Password)
	if err != nil {
		utils.RespondWithError(w, status, err.Error())
		return
	}
	utils.RespondWithJSON(w, status, resp)
}

func (rm *RouteManager) Profile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	resp, status, err := rm.AuthService.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		utils.RespondWithError(w, status, err.Error())
		return
	}
	utils.RespondWithJSON(w, status, resp)
}

func (rm *RouteManager) CreateAchievement(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req models.CreateAchievementRequest
	json.NewDecoder(r.Body).Decode(&req)
	
	resp, status, err := rm.AchievementService.CreateDraft(r.Context(), claims.UserID, &req)
	if err != nil {
		utils.RespondWithError(w, status, err.Error())
		return
	}
	utils.RespondWithJSON(w, status, resp)
}

func (rm *RouteManager) DeleteAchievement(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])

	status, err := rm.AchievementService.DeleteDraft(r.Context(), claims.UserID, id)
	if err != nil {
		utils.RespondWithError(w, status, err.Error())
		return
	}
	utils.RespondWithJSON(w, status, map[string]string{"message": "deleted"})
}

func (rm *RouteManager) SubmitAchievement(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])

	resp, status, err := rm.AchievementService.SubmitForVerification(r.Context(), claims.UserID, id)
	if err != nil {
		utils.RespondWithError(w, status, err.Error())
		return
	}
	utils.RespondWithJSON(w, status, resp)
}

func (rm *RouteManager) VerifyAchievement(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])

	resp, status, err := rm.AchievementService.VerifyAchievement(r.Context(), claims.UserID, id)
	if err != nil {
		utils.RespondWithError(w, status, err.Error())
		return
	}
	utils.RespondWithJSON(w, status, resp)
}

func (rm *RouteManager) RejectAchievement(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])

	resp, status, err := rm.AchievementService.RejectAchievement(r.Context(), claims.UserID, id, "Rejected")
	if err != nil {
		utils.RespondWithError(w, status, err.Error())
		return
	}
	utils.RespondWithJSON(w, status, resp)
}