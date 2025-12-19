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
	MongoDatabaseName           = "prestasi_db"
	MongoCollectionAchievements = "achievements"
)

type AchievementRepository interface {
	CreateAchievementAndReference(ctx context.Context, achievement *models.Achievement, studentID uuid.UUID) (*models.AchievementReference, error)
	SoftDeleteAchievementAndReference(ctx context.Context, achievementRefID uuid.UUID, studentID uuid.UUID) error
	UpdateReferenceStatus(ctx context.Context, refID uuid.UUID, currentStatus string, newStatus string, rejectionNote string, verifiedBy uuid.UUID) (*models.AchievementReference, error)
	GetAchievementDetail(ctx context.Context, mongoID string) (*models.Achievement, error)
	GetReferenceByID(ctx context.Context, refID uuid.UUID) (*models.AchievementReference, error)
	UpdateAchievement(ctx context.Context, mongoID string, update interface{}) error
	UpdateReferenceUpdatedAt(ctx context.Context, refID uuid.UUID) (*models.AchievementReference, error)
	ListAchievementReferences(ctx context.Context, studentIDs []uuid.UUID) ([]models.AchievementReference, error)
	GetStatsByStatus(ctx context.Context, studentID *uuid.UUID) (map[string]int, error)
	GetStatsByType(ctx context.Context, studentID *uuid.UUID) (map[string]int, error)
	// NEW: Hard Delete
	HardDeleteAchievement(ctx context.Context, refID uuid.UUID) error
}

type achievementRepository struct {
	pgDB        *pgxpool.Pool
	mongoClient *mongo.Client
}

func NewAchievementRepository(pgDB *pgxpool.Pool, mongoClient *mongo.Client) AchievementRepository {
	return &achievementRepository{pgDB: pgDB, mongoClient: mongoClient}
}

// GetReferenceByID
func (r *achievementRepository) GetReferenceByID(ctx context.Context, refID uuid.UUID) (*models.AchievementReference, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note FROM achievement_references WHERE id = $1 AND is_deleted = FALSE`
	ref := models.AchievementReference{}
	err := r.pgDB.QueryRow(ctx, query, refID).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, 
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("achievement reference not found")
	}
	return &ref, err
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

	_, err = r.pgDB.Exec(ctx, "UPDATE achievement_references SET is_deleted = true, updated_at = NOW() WHERE id = $1", achievementRefID)
	return err
}

// UpdateReferenceStatus
func (r *achievementRepository) UpdateReferenceStatus(ctx context.Context, refID uuid.UUID, currentStatus string, newStatus string, rejectionNote string, verifiedBy uuid.UUID) (*models.AchievementReference, error) {
	
	ref := models.AchievementReference{}
	
	query := `UPDATE achievement_references SET status = $1, updated_at = NOW()`
	args := []interface{}{newStatus}
	argID := 2
	
	// Perbaikan S1039: Hapus fmt.Sprintf jika tidak ada format verbs (%s, %d, dll)
	if newStatus == "submitted" {
		query += ", submitted_at = NOW()"
	}
	if newStatus == "verified" {
		query += fmt.Sprintf(", verified_at = NOW(), verified_by = $%d", argID)
		args = append(args, verifiedBy)
		argID++
	}
	if newStatus == "rejected" {
		query += fmt.Sprintf(", rejection_note = $%d, verified_by = $%d", argID, argID+1)
		args = append(args, rejectionNote, verifiedBy)
		argID += 2
	}
	
	query += fmt.Sprintf(" WHERE id = $%d AND status = $%d RETURNING id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note", argID, argID+1)
	args = append(args, refID, currentStatus)

	err := r.pgDB.QueryRow(ctx, query, args...).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, 
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
	)
	
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("achievement not found or status already changed")
	}
	return &ref, err
}

// UpdateAchievement
func (r *achievementRepository) UpdateAchievement(ctx context.Context, mongoID string, update interface{}) error {
	objID, _ := primitive.ObjectIDFromHex(mongoID)
	coll := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)
	
	_, err := coll.UpdateOne(ctx, bson.M{"_id": objID, "isDeleted": false}, update)
	return err
}

// UpdateReferenceUpdatedAt
func (r *achievementRepository) UpdateReferenceUpdatedAt(ctx context.Context, refID uuid.UUID) (*models.AchievementReference, error) {
	query := `UPDATE achievement_references SET updated_at = NOW() WHERE id = $1 RETURNING id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at, is_deleted`
	ref := models.AchievementReference{}
	
	err := r.pgDB.QueryRow(ctx, query, refID).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, 
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt, &ref.IsDeleted,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("achievement reference not found")
	}
	return &ref, err
}

// GetAchievementDetail
func (r *achievementRepository) GetAchievementDetail(ctx context.Context, mongoID string) (*models.Achievement, error) {
	objID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return nil, errors.New("invalid mongo ID")
	}
	coll := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)
	
	var achievement models.Achievement
	err = coll.FindOne(ctx, bson.M{"_id": objID, "isDeleted": false}).Decode(&achievement)
	return &achievement, err
}

// ListAchievementReferences
func (r *achievementRepository) ListAchievementReferences(ctx context.Context, studentIDs []uuid.UUID) ([]models.AchievementReference, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note FROM achievement_references WHERE is_deleted = FALSE`
	args := []interface{}{}
	
	if len(studentIDs) > 0 {
		query += fmt.Sprintf(" AND student_id = ANY($%d)", 1)
		args = append(args, studentIDs)
	}

	rows, err := r.pgDB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		if err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, 
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		); err != nil {
			return nil, fmt.Errorf("error scanning achievement reference: %w", err)
		}
		achievements = append(achievements, ref)
	}
	return achievements, nil
}

// GetStatsByStatus menghitung jumlah prestasi berdasarkan status dari PostgreSQL
func (r *achievementRepository) GetStatsByStatus(ctx context.Context, studentID *uuid.UUID) (map[string]int, error) {
	query := `SELECT status, COUNT(*) FROM achievement_references WHERE is_deleted = FALSE`
	args := []interface{}{}
	
	if studentID != nil {
		query += ` AND student_id = $1`
		args = append(args, *studentID)
	}
	
	query += ` GROUP BY status`

	rows, err := r.pgDB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}
	return stats, nil
}

// GetStatsByType menghitung jumlah prestasi berdasarkan tipe dari MongoDB
func (r *achievementRepository) GetStatsByType(ctx context.Context, studentID *uuid.UUID) (map[string]int, error) {
	coll := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)
	
	matchStage := bson.M{"isDeleted": false}
	if studentID != nil {
		matchStage["studentId"] = *studentID
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$achievementType",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stats := make(map[string]int)
	var results []struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for _, res := range results {
		stats[res.ID] = res.Count
	}
	return stats, nil
}

func (r *achievementRepository) HardDeleteAchievement(ctx context.Context, refID uuid.UUID) error {
	var mongoID string
	err := r.pgDB.QueryRow(ctx, "SELECT mongo_achievement_id FROM achievement_references WHERE id = $1", refID).Scan(&mongoID)
	if err != nil {
		return fmt.Errorf("reference not found: %w", err)
	}

	objID, _ := primitive.ObjectIDFromHex(mongoID)
	coll := r.mongoClient.Database(MongoDatabaseName).Collection(MongoCollectionAchievements)
	_, err = coll.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete mongo doc: %w", err)
	}

	_, err = r.pgDB.Exec(ctx, "DELETE FROM achievement_references WHERE id = $1", refID)
	return err
}