package repositories

import (
	"context"
	"errors"
	"fmt"

	"prestasi-mahasiswa-api/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository interface {
	GetAllRoles(ctx context.Context) ([]models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error)
}

type roleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &roleRepository{db: db}
}

// GetRoleByName mengambil detail role berdasarkan nama, termasuk permissions
func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	query := `
		SELECT 
			r.id, r.name, ARRAY_AGG(p.name) AS permissions
		FROM 
			roles r
		LEFT JOIN 
			role_permissions rp ON r.id = rp.role_id
		LEFT JOIN 
			permissions p ON rp.permission_id = p.id
		WHERE 
			r.name = $1
		GROUP BY 
			r.id, r.name`

	role := models.Role{}
	var permissionsPgArray []string
	err := r.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &permissionsPgArray)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Role not found
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	
	role.Permissions = permissionsPgArray
	return &role, nil
}

// GetAllRoles, GetRoleByID (metode lainnya dapat diimplementasikan serupa jika diperlukan)
// Saat ini hanya GetRoleByName yang diimplementasikan karena langsung dibutuhkan oleh UserService.
func (r *roleRepository) GetAllRoles(ctx context.Context) ([]models.Role, error) {
    // Implementasi untuk mengambil semua roles
    return nil, nil // Placeholder
}
func (r *roleRepository) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error) {
    // Implementasi untuk mengambil role berdasarkan ID
    return nil, nil // Placeholder
}