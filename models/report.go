package models

// DashboardStats merepresentasikan ringkasan statistik untuk dashboard
type DashboardStats struct {
	TotalAchievements int            `json:"totalAchievements"`
	ByStatus          map[string]int `json:"byStatus"`
	ByType            map[string]int `json:"byType"`
}

// StudentStats merepresentasikan statistik spesifik mahasiswa
type StudentStats struct {
	StudentID         string         `json:"studentId"`
	TotalPoints       int            `json:"totalPoints"`
	TotalAchievements int            `json:"totalAchievements"`
	ByStatus          map[string]int `json:"byStatus"`
}