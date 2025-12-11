package models

import (
	"time"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AchievementReference (PostgreSQL)
type AchievementReference struct {
	ID                 uuid.UUID  `json:"id"`
	StudentID          uuid.UUID  `json:"studentId"`
	MongoAchievementID string     `json:"mongoAchievementId"`
	Status             string     `json:"status"`
	SubmittedAt        *time.Time `json:"submittedAt"`
	VerifiedAt         *time.Time `json:"verifiedAt"`
	VerifiedBy         *uuid.UUID `json:"verifiedBy"`
	RejectionNote      *string    `json:"rejectionNote"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	IsDeleted          bool       `json:"isDeleted"`
}

// Achievement (MongoDB)
type Achievement struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty"`
	StudentUUID     uuid.UUID              `bson:"studentId"`
	AchievementType string                 `bson:"achievementType"`
	Title           string                 `bson:"title"`
	Description     string                 `bson:"description"`
	Details         map[string]interface{} `bson:"details"`
	Attachments     []AttachmentFile       `bson:"attachments"`
	Tags            []string               `bson:"tags"`
	Points          int                    `bson:"points"`
	CreatedAt       time.Time              `bson:"createdAt"`
	UpdatedAt       time.Time              `bson:"updatedAt"`
	IsDeleted       bool                   `bson:"isDeleted"`
}

type AttachmentFile struct {
	FileName   string    `bson:"fileName"`
	FileUrl    string    `bson:"fileUrl"`
	FileType   string    `bson:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt"`
}

// Struct untuk Request Creation
type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievementType"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
}

// Struct untuk Response Detail (Gabungan SQL + Mongo) - INI YANG HILANG SEBELUMNYA
type AchievementDetailResponse struct {
	RefID           uuid.UUID              `json:"refId"`
	Status          string                 `json:"status"`
	AchievementType string                 `json:"achievementType"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Points          int                    `json:"points"`
	RejectionNote   *string                `json:"rejectionNote,omitempty"`
	SubmittedAt     *time.Time             `json:"submittedAt,omitempty"`
	VerifiedAt      *time.Time             `json:"verifiedAt,omitempty"`
}