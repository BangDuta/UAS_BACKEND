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
	GetAdviseeStudentUserIDsByAdvisorUserID(ctx context.Context, advisorUserID uuid.UUID) ([]uuid.UUID, error)
	
	// Admin CRUD
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req *models.UpdateUserRequest, roleID *uuid.UUID) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ListAllUsers(ctx context.Context) ([]models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// scanUserRow helper
func scanUserRow(row pgx.Row) (*models.User, error) {
	user := models.User{}
	var permissionsPgArray []string
	var roleID uuid.UUID 

	err := row.Scan(
		&user.ID, 
		&user.Username, 
		&user.Email, 
		&user.PasswordHash,
		&user.FullName, 
		&roleID, 
		&user.Role, 
		&user.IsActive, 
		&permissionsPgArray,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	user.RoleID = roleID 
	user.Permissions = permissionsPgArray
	return &user, nil
}

func (r *userRepository) FindUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, r.name AS role, u.is_active, 
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
	return scanUserRow(r.db.QueryRow(ctx, query, strings.ToLower(identifier)))
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, r.name AS role, u.is_active, 
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
	return scanUserRow(r.db.QueryRow(ctx, query, id))
}

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

// --- Admin CRUD Implementation ---

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id`
	
	var createdID uuid.UUID
	err := r.db.QueryRow(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive,
	).Scan(&createdID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return r.GetUserByID(ctx, createdID)
}

func (r *userRepository) UpdateUser(ctx context.Context, userID uuid.UUID, req *models.UpdateUserRequest, newRoleID *uuid.UUID) (*models.User, error) {
	sets := []string{}
	values := []interface{}{}
	i := 1

	if req.FullName != nil {
		sets = append(sets, fmt.Sprintf("full_name = $%d", i)); values = append(values, *req.FullName); i++
	}
	if req.Email != nil {
		sets = append(sets, fmt.Sprintf("email = $%d", i)); values = append(values, *req.Email); i++
	}
	if req.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", i)); values = append(values, *req.IsActive); i++
	}
	if newRoleID != nil {
		sets = append(sets, fmt.Sprintf("role_id = $%d", i)); values = append(values, *newRoleID); i++
	}

	if len(sets) == 0 {
		return r.GetUserByID(ctx, userID)
	}

	sets = append(sets, "updated_at = NOW()")
	values = append(values, userID)
	
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(sets, ", "), i)

	_, err := r.db.Exec(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return r.GetUserByID(ctx, userID)
}

func (r *userRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *userRepository) ListAllUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT 
			u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, r.name AS role_name, u.is_active, 
            ARRAY(
                SELECT p.name 
                FROM role_permissions rp 
                JOIN permissions p ON rp.permission_id = p.id 
                WHERE rp.role_id = r.id
            ) AS permissions
		FROM users u
		JOIN roles r ON u.role_id = r.id
		ORDER BY u.username`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
        return nil, fmt.Errorf("database query failed: %w", err)
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        // Kita tidak bisa pakai scanUserRow karena itu untuk QueryRow (pgx.Row) bukan pgx.Rows
		user := models.User{}
        var permissionsPgArray []string
        var roleID uuid.UUID 

        err := rows.Scan(
            &user.ID, &user.Username, &user.Email, &user.PasswordHash,
            &user.FullName, &roleID, &user.Role, &user.IsActive, &permissionsPgArray,
        )
        if err != nil {
            return nil, err
        }
        user.RoleID = roleID
        user.Permissions = permissionsPgArray
        users = append(users, user)
    }
	return users, nil
}