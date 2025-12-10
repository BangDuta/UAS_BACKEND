package models

import (
	"time"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AchievementReference merepresentasikan data dari tabel achievement_references (PostgreSQL)
type AchievementReference struct {
    ID                   uuid.UUID `json:"id"`
    StudentID            uuid.UUID `json:"studentId"`
    MongoAchievementID   string    `json:"mongoAchievementId"` // FOREIGN KEY ke MongoDB
    Status               string    `json:"status"` // ENUM: 'draft', 'submitted', 'verified', 'rejected' [cite: 96]
    SubmittedAt          *time.Time `json:"submittedAt"`
    VerifiedAt           *time.Time `json:"verifiedAt"`
    VerifiedBy           *uuid.UUID `json:"verifiedBy"`
    RejectionNote        *string   `json:"rejectionNote"`
    CreatedAt            time.Time `json:"createdAt"`
    UpdatedAt            time.Time `json:"updatedAt"`
}

// Achievement merepresentasikan data dari collection achievements (MongoDB) [cite: 107]
type Achievement struct {
    ID              primitive.ObjectID `bson:"_id,omitempty"`
    StudentUUID     uuid.UUID          `bson:"studentId"` // Reference to PostgreSQL [cite: 110]
    AchievementType string             `bson:"achievementType"` // e.g., 'competition', 'publication' [cite: 111]
    Title           string             `bson:"title"` [cite: 112]
    Description     string             `bson:"description"` [cite: 113]
    Details         map[string]interface{} `bson:"details"` // Field dinamis [cite: 114]
    Attachments     []AttachmentFile   `bson:"attachments"`
    Tags            []string           `bson:"tags"` [cite: 153]
    Points          int                `bson:"points"` [cite: 154]
    CreatedAt       time.Time          `bson:"createdAt"` [cite: 155]
    UpdatedAt       time.Time          `bson:"updatedAt"` [cite: 156]
}

type AttachmentFile struct {
    FileName string    `bson:"fileName"` [cite: 148]
    FileUrl  string    `bson:"fileUrl"` [cite: 149]
    FileType string    `bson:"fileType"` [cite: 150]
    UploadedAt time.Time `bson:"uploadedAt"` [cite: 151]
}