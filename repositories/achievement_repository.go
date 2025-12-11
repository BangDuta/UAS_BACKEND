package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"prestasi-mahasiswa-api/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MongoDatabaseName  = "prestasi_db"
	MongoCollectionAchievements = "achievements"
)

type AchievementRepository interface {
	CreateAchievementAndReference(ctx context.Context, achievement *models.Achievement, studentID uuid.UUID) (*models.AchievementReference, error)
	SoftDeleteAchievementAndReference(ctx context.Context, achievementRefID uuid.UUID, studentID uuid.UUID) error
	// Metode untuk workflow/read akan ditambahkan di commit berikutnya
}

type achievementRepository struct {
	pgDB *pgxpool.Pool
	mongoClient *mongo.Client
}

func NewAchievementRepository(pgDB *pgxpool.Pool, mongoClient *mongo.Client) AchievementRepository {
	return &achievementRepository{pgDB: pgDB, mongoClient: mongoClient}
}

// CreateAchievementAndReference menangani penyimpanan ke MongoDB dan PostgreSQL (FR-003)
func (r *achievementRepository) CreateAchievementAndReference(ctx context.Context, achievement *models.Achievement, studentID uuid.UUID) (*models.AchievementReference, error) {
	// 1. Simpan ke MongoDB
	mongoCollection := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)
	
	achievement.ID = primitive.NewObjectID()
	achievement.StudentUUID = studentID
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	achievement.IsDeleted = false // Tambahkan field IsDeleted untuk soft delete

	_, err := mongoCollection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, fmt.Errorf("failed to insert achievement into MongoDB: %w", err)
	}

	// 2. Simpan reference ke PostgreSQL
	ref := models.AchievementReference{
		ID: uuid.New(),
		StudentID: studentID,
		MongoAchievementID: achievement.ID.Hex(),
		Status: "draft", // Status awal: 'draft' (FR-003)
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	query := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, created_at, updated_at, is_deleted) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = r.pgDB.Exec(ctx, query, 
		ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status, ref.CreatedAt, ref.UpdatedAt, ref.IsDeleted)
	
	if err != nil {
		// Rollback MongoDB insertion
		_, _ = mongoCollection.DeleteOne(ctx, bson.M{"_id": achievement.ID})
		return nil, fmt.Errorf("failed to insert reference into PostgreSQL: %w", err)
	}

	return &ref, nil
}

// SoftDeleteAchievementAndReference (FR-005)
func (r *achievementRepository) SoftDeleteAchievementAndReference(ctx context.Context, achievementRefID uuid.UUID, studentID uuid.UUID) error {
	// 1. Cari reference, pastikan milik studentID dan berstatus 'draft'
	var mongoID string
	queryRef := "SELECT mongo_achievement_id FROM achievement_references WHERE id = $1 AND student_id = $2 AND status = 'draft' AND is_deleted = FALSE"
	err := r.pgDB.QueryRow(ctx, queryRef, achievementRefID, studentID).Scan(&mongoID)
	
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("achievement reference not found, status is not 'draft', or access denied")
	}
	if err != nil {
		return fmt.Errorf("failed to get mongo ID: %w", err)
	}

	// 2. Soft delete data di MongoDB
	objID, _ := primitive.ObjectIDFromHex(mongoID)
	mongoCollection := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)
	updateResult, err := mongoCollection.UpdateOne(ctx, 
		bson.M{"_id": objID, "isDeleted": false}, // Pastikan belum dihapus
		bson.M{"$set": bson.M{"isDeleted": true, "updatedAt": time.Now()}})
	
	if err != nil || updateResult.ModifiedCount == 0 {
		return errors.New("failed to soft delete achievement in MongoDB (or already deleted)")
	}

	// 3. Update reference di PostgreSQL
	queryUpdateRef := "UPDATE achievement_references SET is_deleted = true, updated_at = $1 WHERE id = $2"
	_, err = r.pgDB.Exec(ctx, queryUpdateRef, time.Now(), achievementRefID)
	if err != nil {
		return fmt.Errorf("failed to update reference status in PostgreSQL: %w", err)
	}

	return nil
}