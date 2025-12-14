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

type ProfileRepository interface {
	UpsertStudent(ctx context.Context, student *models.Student) (*models.Student, error)
	UpsertLecturer(ctx context.Context, lecturer *models.Lecturer) (*models.Lecturer, error)
	AssignAdvisor(ctx context.Context, studentUserID uuid.UUID, advisorLecturerID uuid.UUID) error
	GetLecturerByUserID(ctx context.Context, userID uuid.UUID) (*models.Lecturer, error)
}

type profileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) ProfileRepository {
	return &profileRepository{db: db}
}

// UpsertStudent (Insert atau Update jika user_id sudah ada)
func (r *profileRepository) UpsertStudent(ctx context.Context, s *models.Student) (*models.Student, error) {
	query := `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			student_id = EXCLUDED.student_id,
			program_study = EXCLUDED.program_study,
			academic_year = EXCLUDED.academic_year
		RETURNING id, advisor_id`
	
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	err := r.db.QueryRow(ctx, query, s.ID, s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear).Scan(&s.ID, &s.AdvisorID)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert student: %w", err)
	}
	return s, nil
}

// UpsertLecturer
func (r *profileRepository) UpsertLecturer(ctx context.Context, l *models.Lecturer) (*models.Lecturer, error) {
	query := `
		INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET
			lecturer_id = EXCLUDED.lecturer_id,
			department = EXCLUDED.department
		RETURNING id`

	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}

	err := r.db.QueryRow(ctx, query, l.ID, l.UserID, l.LecturerID, l.Department).Scan(&l.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert lecturer: %w", err)
	}
	return l, nil
}

// AssignAdvisor menghubungkan student ke lecturer
func (r *profileRepository) AssignAdvisor(ctx context.Context, studentUserID uuid.UUID, advisorLecturerID uuid.UUID) error {
	query := `UPDATE students SET advisor_id = $1 WHERE user_id = $2`
	cmd, err := r.db.Exec(ctx, query, advisorLecturerID, studentUserID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("student profile not found for this user")
	}
	return nil
}

// GetLecturerByUserID
func (r *profileRepository) GetLecturerByUserID(ctx context.Context, userID uuid.UUID) (*models.Lecturer, error) {
	query := `SELECT id, user_id, lecturer_id, department FROM lecturers WHERE user_id = $1`
	var l models.Lecturer
	err := r.db.QueryRow(ctx, query, userID).Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("lecturer profile not found")
	}
	return &l, err
}