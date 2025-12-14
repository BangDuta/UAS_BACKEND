package models

import (
	"time"
	"github.com/google/uuid"
)

// Student merepresentasikan tabel students
type Student struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"userId"`
	StudentID    string     `json:"studentId"` // NIM
	ProgramStudy string     `json:"programStudy"`
	AcademicYear string     `json:"academicYear"`
	AdvisorID    *uuid.UUID `json:"advisorId"` // Link ke ID tabel Lecturers (bukan User ID)
	CreatedAt    time.Time  `json:"createdAt"`
}

// Lecturer merepresentasikan tabel lecturers
type Lecturer struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	LecturerID string    `json:"lecturerId"` // NIP
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Request Payload untuk Set Profil Mahasiswa
type StudentProfileRequest struct {
	StudentID    string `json:"studentId"`    // NIM
	ProgramStudy string `json:"programStudy"`
	AcademicYear string `json:"academicYear"`
}

// Request Payload untuk Set Profil Dosen
type LecturerProfileRequest struct {
	LecturerID string `json:"lecturerId"` // NIP
	Department string `json:"department"`
}

// Request Payload untuk Assign Dosen Wali
type AssignAdvisorRequest struct {
	AdvisorUserID string `json:"advisorUserId"` // Kita input User ID Dosen, nanti sistem cari Lecturer ID-nya
}