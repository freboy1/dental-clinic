package repository

import (
	"context"

	"dental_clinic/internal/modules/clinic_admin/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClinicAdminRepository interface {
	Create(admin *models.ClinicAdmin) (*models.ClinicAdmin, error)
	GetAll() ([]models.ClinicAdmin, error)
	GetByID(id string) (*models.ClinicAdmin, error)
	Update(admin *models.ClinicAdmin) (*models.ClinicAdmin, error)
	Delete(id string) error
}

type clinicAdminRepo struct {
	db *pgxpool.Pool
}

func NewClinicAdminRepository(db *pgxpool.Pool) ClinicAdminRepository {
	return &clinicAdminRepo{db: db}
}

func (r *clinicAdminRepo) Create(admin *models.ClinicAdmin) (*models.ClinicAdmin, error) {
	query := `
		INSERT INTO clinic_admins (id, clinic_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, clinic_id, user_id, created_at
	`
	err := r.db.QueryRow(context.Background(), query, admin.Id, admin.ClinicID, admin.UserID, admin.CreatedAt).
		Scan(&admin.Id, &admin.ClinicID, &admin.UserID, &admin.CreatedAt)
	return admin, err
}

func (r *clinicAdminRepo) GetAll() ([]models.ClinicAdmin, error) {
	query := `
		SELECT ca.id, ca.clinic_id, ca.user_id, u.name, u.email, ca.created_at
		FROM clinic_admins ca
		JOIN users u ON u.id = ca.user_id
		ORDER BY ca.created_at DESC
	`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	admins := make([]models.ClinicAdmin, 0)
	for rows.Next() {
		var admin models.ClinicAdmin
		if err := rows.Scan(&admin.Id, &admin.ClinicID, &admin.UserID, &admin.Name, &admin.Email, &admin.CreatedAt); err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}
	return admins, rows.Err()
}

func (r *clinicAdminRepo) GetByID(id string) (*models.ClinicAdmin, error) {
	query := `
		SELECT ca.id, ca.clinic_id, ca.user_id, u.name, u.email, ca.created_at
		FROM clinic_admins ca
		JOIN users u ON u.id = ca.user_id
		WHERE ca.id = $1
	`
	admin := &models.ClinicAdmin{}
	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&admin.Id, &admin.ClinicID, &admin.UserID, &admin.Name, &admin.Email, &admin.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return admin, nil
}

func (r *clinicAdminRepo) Update(admin *models.ClinicAdmin) (*models.ClinicAdmin, error) {
	query := `
		UPDATE clinic_admins
		SET clinic_id = $2
		WHERE id = $1
		RETURNING id, clinic_id, user_id, created_at
	`
	err := r.db.QueryRow(context.Background(), query, admin.Id, admin.ClinicID).
		Scan(&admin.Id, &admin.ClinicID, &admin.UserID, &admin.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return admin, nil
}

func (r *clinicAdminRepo) Delete(id string) error {
	result, err := r.db.Exec(context.Background(), `DELETE FROM clinic_admins WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
