package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	// "dental_clinic/internal"

	"dental_clinic/internal/modules/clinic/models"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClinicRepository interface {
	Create(clinic *models.Clinic) (*models.Clinic, error)
	GetAll() ([]*models.Clinic, error)
	GetByID(id uuid.UUID) (*models.Clinic, error)
	Update(clinic *models.Clinic) (*models.Clinic, error)
	Delete(id uuid.UUID) error
	AddAddress(id, clinic_id uuid.UUID, address_id string, is_main bool) error
	GetClinicAddress(id uuid.UUID) ([]models.ClinicAddress, error)
	DeleteAddress(id, address_id uuid.UUID) error
}

type clinicRepo struct {
	db *pgxpool.Pool
}

func NewClinicRepository(db *pgxpool.Pool) ClinicRepository {
	return &clinicRepo{db: db}
}

func (r *clinicRepo) Create(clinic *models.Clinic) (*models.Clinic, error) {
	query := `INSERT INTO clinics (id, name, description, phone, email, website, is_active, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
            RETURNING id, created_at`

	err := r.db.QueryRow(
		context.Background(),
		query,
		clinic.Id,
		clinic.Name,
		clinic.Description,
		clinic.Phone,
		clinic.Email,
		clinic.Website,
		clinic.IsActive,
		clinic.CreatedAt,
	).Scan(&clinic.Id, &clinic.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create clinic: %w", err)
	}

	return clinic, nil
}

func (r *clinicRepo) GetAll() ([]*models.Clinic, error) {
	query := `SELECT id, name, description, phone, email, website, is_active, created_at
            FROM clinics
            WHERE is_active = true
            ORDER BY created_at DESC`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinics: %w", err)
	}
	defer rows.Close()

	var clinics []*models.Clinic
	for rows.Next() {
		clinic := &models.Clinic{}
		err := rows.Scan(
			&clinic.Id,
			&clinic.Name,
			&clinic.Description,
			&clinic.Phone,
			&clinic.Email,
			&clinic.Website,
			&clinic.IsActive,
			&clinic.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan clinic: %w", err)
		}
		clinics = append(clinics, clinic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clinics: %w", err)
	}

	return clinics, nil
}

func (r *clinicRepo) GetByID(id uuid.UUID) (*models.Clinic, error) {
	query := `SELECT id, name, description, phone, email, website, is_active, created_at
            FROM clinics
            WHERE id = $1`

	clinic := &models.Clinic{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&clinic.Id,
		&clinic.Name,
		&clinic.Description,
		&clinic.Phone,
		&clinic.Email,
		&clinic.Website,
		&clinic.IsActive,
		&clinic.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get clinic: %w", err)
	}

	return clinic, nil
}

func (r *clinicRepo) Update(clinic *models.Clinic) (*models.Clinic, error) {
	query := `UPDATE clinics 
            SET name = $2, 
                description = $3, 
                phone = $4, 
                email = $5, 
                website = $6, 
                is_active = $7
            WHERE id = $1
            RETURNING id, created_at`

	err := r.db.QueryRow(
		context.Background(),
		query,
		clinic.Id,
		clinic.Name,
		clinic.Description,
		clinic.Phone,
		clinic.Email,
		clinic.Website,
		clinic.IsActive,
	).Scan(&clinic.Id, &clinic.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update clinic: %w", err)
	}

	return clinic, nil
}

func (r *clinicRepo) Delete(id uuid.UUID) error {
	query := `UPDATE clinics SET is_active = false WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete clinic: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("clinic not found")
	}

	return nil
}


func (r *clinicRepo) AddAddress(id, clinic_id uuid.UUID, address_id string, is_main bool) error {
	query := `INSERT INTO clinic_addresses (id, clinic_id, address_id, is_main)
            VALUES ($1, $2, $3, $4) 
            `

	if is_main {
		err := r.UnsetMainAddress(clinic_id)
		if err != nil {
			return err
		}
	}
	
	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
		clinic_id,
		address_id,
		is_main,
	)

	if err != nil {
		return fmt.Errorf("failed to create clinic address: %w", err)
	}

	return nil
}

func (r *clinicRepo) UnsetMainAddress(clinic_id uuid.UUID) error {
	query := `UPDATE clinic_addresses SET is_main = false WHERE clinic_id = $1`

	_, err := r.db.Exec(context.Background(), query, clinic_id)
	if err != nil {
		return fmt.Errorf("failed to update clinic address: %w", err)
	}

	// if result.RowsAffected() == 0 {
	// 	return fmt.Errorf("clinic address not found")
	// }

	return nil
}


func (r *clinicRepo) GetClinicAddress(id uuid.UUID) ([]models.ClinicAddress, error) {
	query := `SELECT id, clinic_id, address_id, is_main
            FROM clinic_addresses
            WHERE clinic_id = $1
            `

	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinic address: %w", err)
	}
	defer rows.Close()

	var clinics []models.ClinicAddress
	for rows.Next() {
		clinic := models.ClinicAddress{}
		err := rows.Scan(
			&clinic.Id,
			&clinic.ClinicId,
			&clinic.AddressId,
			&clinic.IsMain,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan clinic: %w", err)
		}
		clinics = append(clinics, clinic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clinics: %w", err)
	}

	return clinics, nil
}


func (r *clinicRepo) DeleteAddress(id, address_id uuid.UUID) error {
	query := `DELETE clinic_addresses WHERE address_id = $1 AND clinic_id = $2`

	result, err := r.db.Exec(context.Background(), query, address_id, id)
	if err != nil {
		return fmt.Errorf("failed to delete clinic address: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("clinic address not found")
	}

	return nil
}