package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"
	// "strings" dihapus karena unused

	"prestasi-mahasiswa-api/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MongoDatabaseName           = "prestasi_db"
	MongoCollectionAchievements = "achievements"
)

type AchievementRepository interface {
	CreateAchievementAndReference(ctx context.Context, achievement *models.Achievement, studentID uuid.UUID) (*models.AchievementReference, error)
	SoftDeleteAchievementAndReference(ctx context.Context, achievementRefID uuid.UUID, studentID uuid.UUID) error
	UpdateReferenceStatus(ctx context.Context, refID uuid.UUID, currentStatus, newStatus, rejectionNote string, verifiedBy uuid.UUID) (*models.AchievementReference, error)
	GetAchievementDetail(ctx context.Context, mongoID string) (*models.Achievement, error)
}

type achievementRepository struct {
	pgDB        *pgxpool.Pool
	mongoClient *mongo.Client
}

func NewAchievementRepository(pgDB *pgxpool.Pool, mongoClient *mongo.Client) AchievementRepository {
	return &achievementRepository{pgDB: pgDB, mongoClient: mongoClient}
}

// CreateAchievementAndReference
func (r *achievementRepository) CreateAchievementAndReference(ctx context.Context, achievement *models.Achievement, studentID uuid.UUID) (*models.AchievementReference, error) {
	mongoCollection := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)

	achievement.ID = primitive.NewObjectID()
	achievement.StudentUUID = studentID
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	achievement.IsDeleted = false

	_, err := mongoCollection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into MongoDB: %w", err)
	}

	ref := models.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          studentID,
		MongoAchievementID: achievement.ID.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		IsDeleted:          false,
	}

	query := `INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, created_at, updated_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = r.pgDB.Exec(ctx, query, ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status, ref.CreatedAt, ref.UpdatedAt, ref.IsDeleted)

	if err != nil {
		_, _ = mongoCollection.DeleteOne(ctx, bson.M{"_id": achievement.ID})
		return nil, fmt.Errorf("failed to insert reference: %w", err)
	}
	return &ref, nil
}

// SoftDeleteAchievementAndReference
func (r *achievementRepository) SoftDeleteAchievementAndReference(ctx context.Context, achievementRefID uuid.UUID, studentID uuid.UUID) error {
	var mongoID string
	queryRef := "SELECT mongo_achievement_id FROM achievement_references WHERE id = $1 AND student_id = $2 AND status = 'draft' AND is_deleted = FALSE"
	err := r.pgDB.QueryRow(ctx, queryRef, achievementRefID, studentID).Scan(&mongoID)

	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("achievement not found or not editable")
	}

	objID, _ := primitive.ObjectIDFromHex(mongoID)
	mongoColl := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)
	_, err = mongoColl.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"isDeleted": true}})
	if err != nil {
		return err
	}

	_, err = r.pgDB.Exec(ctx, "UPDATE achievement_references SET is_deleted = true WHERE id = $1", achievementRefID)
	return err
}

// UpdateReferenceStatus
func (r *achievementRepository) UpdateReferenceStatus(ctx context.Context, refID uuid.UUID, currentStatus, newStatus, rejectionNote string, verifiedBy uuid.UUID) (*models.AchievementReference, error) {
	query := `UPDATE achievement_references SET status = $1, updated_at = $2 WHERE id = $3 AND status = $4 RETURNING id, status`
	
	ref := models.AchievementReference{}
	err := r.pgDB.QueryRow(ctx, query, newStatus, time.Now(), refID, currentStatus).Scan(&ref.ID, &ref.Status)
	
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

// GetAchievementDetail
func (r *achievementRepository) GetAchievementDetail(ctx context.Context, mongoID string) (*models.Achievement, error) {
	objID, _ := primitive.ObjectIDFromHex(mongoID)
	coll := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)
	
	var achievement models.Achievement
	err := coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&achievement)
	return &achievement, err
}