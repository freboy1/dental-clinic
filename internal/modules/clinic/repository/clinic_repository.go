package repository

import (
	"context"
	
	// "dental_clinic/internal"


	"dental_clinic/internal/modules/clinic/models"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClinicRepository interface {
	Create(clinic *models.Clinic) (*models.Clinic, error)
}

type clinicRepo struct {
	db *pgxpool.Pool
}

func NewClinicRepository(db *pgxpool.Pool) ClinicRepository {
	return &clinicRepo{db: db}
}

func (r *clinicRepo) Create(clinic *models.Clinic) (*models.Clinic, error) {
	query := `INSERT INTO clinics (id, name, description, phone, email, website, is_active, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(context.Background(), query, clinic.Id, clinic.Name, clinic.Description, clinic.Phone, clinic.Email, clinic.Website, clinic.IsActive).
		Scan(&clinic.Id)
	return clinic, err
}
