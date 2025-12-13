package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"prestasi-mahasiswa-api/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	FindUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAdviseeStudentUserIDsByAdvisorUserID(ctx context.Context, advisorUserID uuid.UUID) ([]uuid.UUID, error) // Tambahan untuk FR-006
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// FindUserByUsernameOrEmail (FR-001)
func (r *userRepository) FindUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.full_name, r.name AS role, u.is_active, 
            ARRAY(
                SELECT p.name 
                FROM role_permissions rp 
                JOIN permissions p ON rp.permission_id = p.id 
                WHERE rp.role_id = r.id
            ) AS permissions
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.username = $1 OR u.email = $1
    `
	user := models.User{}
	var permissionsPgArray []string
	
	err := r.db.QueryRow(ctx, query, strings.ToLower(identifier)).Scan(
		&user.ID, 
		&user.Username, 
		&user.Email, 
		&user.PasswordHash,
		&user.FullName, 
		&user.Role, 
		&user.IsActive, 
		&permissionsPgArray,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	user.Permissions = permissionsPgArray
	return &user, nil
}

// GetUserByID
func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.email, u.full_name, r.name AS role, u.is_active, 
            ARRAY(
                SELECT p.name 
                FROM role_permissions rp 
                JOIN permissions p ON rp.permission_id = p.id 
                WHERE rp.role_id = r.id
            ) AS permissions
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.id = $1
    `
	user := models.User{}
	var permissionsPgArray []string
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, 
		&user.Username, 
		&user.Email, 
		&user.FullName, 
		&user.Role, 
		&user.IsActive, 
		&permissionsPgArray,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	user.Permissions = permissionsPgArray
	return &user, nil
}

// GetAdviseeStudentUserIDsByAdvisorUserID (FR-006)
// Mengambil semua User IDs Mahasiswa yang dibimbing oleh Dosen Wali ini
func (r *userRepository) GetAdviseeStudentUserIDsByAdvisorUserID(ctx context.Context, advisorUserID uuid.UUID) ([]uuid.UUID, error) {
    query := `
        SELECT s.user_id
        FROM students s
        JOIN lecturers l ON s.advisor_id = l.id
        WHERE l.user_id = $1
    `
    rows, err := r.db.Query(ctx, query, advisorUserID)
    if err != nil {
        return nil, fmt.Errorf("database query failed: %w", err)
    }
    defer rows.Close()

    var studentUserIDs []uuid.UUID
    for rows.Next() {
        var id uuid.UUID
        if err := rows.Scan(&id); err != nil {
            return nil, fmt.Errorf("error scanning student ID: %w", err)
        }
        studentUserIDs = append(studentUserIDs, id)
    }
    return studentUserIDs, nil
}