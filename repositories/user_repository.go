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
	// Metode untuk CRUD User oleh Admin akan ditambahkan di Commit #8
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// scanUserRow adalah helper untuk menscan hasil query user ke struct models.User
func scanUserRow(row pgx.Row) (*models.User, error) {
	user := models.User{}
	var permissionsPgArray []string
	
	
	err := row.Scan(
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


// FindUserByUsernameOrEmail mengambil user, role, dan permissions untuk Login (FR-001)
func (r *userRepository) FindUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	// Termasuk password_hash untuk validasi login
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
		&user.PasswordHash, // Tambahan untuk Login
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

// GetUserByID mengambil user berdasarkan ID untuk endpoint Profile
func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	// Tidak termasuk password_hash
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

// ... (tambah method ke interface UserRepository)

// GetAdvisorIDByUserID mengambil advisor_id (UUID) dari tabel students
func (r *userRepository) GetAdvisorIDByUserID(ctx context.Context, studentUserID uuid.UUID) (uuid.UUID, error) {
    query := `
        SELECT s.advisor_id 
        FROM students s 
        WHERE s.user_id = $1
    `
    var advisorID uuid.UUID // Ini adalah ID dari tabel Lecturers
    err := r.db.QueryRow(ctx, query, studentUserID).Scan(&advisorID)
    
    if errors.Is(err, pgx.ErrNoRows) {
        return uuid.Nil, errors.New("student advisor not found")
    }
    if err != nil {
        return uuid.Nil, fmt.Errorf("database query failed: %w", err)
    }
    return advisorID, nil
}

// GetAdviseeStudentUserIDsByAdvisorUserID mengambil semua Student USER IDs yang dibimbing oleh Dosen Wali ini
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

// GetAdviseeStudentIDsByAdvisorUserID mengambil semua Student IDs yang dibimbing oleh Dosen Wali ini
func (r *userRepository) GetAdviseeStudentIDsByAdvisorUserID(ctx context.Context, advisorUserID uuid.UUID) ([]uuid.UUID, error) {
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