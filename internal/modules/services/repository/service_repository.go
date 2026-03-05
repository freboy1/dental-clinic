package repository

import (
	"context"
	"dental_clinic/internal/modules/services/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository interface {
	Create(service *models.Service) (*models.Service, error)
	GetByID(id string) (*models.Service, error)
	GetAll() ([]models.Service, error)
	GetByClinicID(clinicID string) ([]models.Service, error)
	Update(id string, service *models.Service) (*models.Service, error)
	Delete(id string) error
}

type serviceRepo struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) ServiceRepository {
	return &serviceRepo{db: db}
}

func (r *serviceRepo) Create(service *models.Service) (*models.Service, error) {
	query := `
		INSERT INTO services (name, description, price, duration, clinic_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		service.Name,
		service.Description,
		service.Price,
		service.Duration,
		service.ClinicID,
		service.IsActive,
	).Scan(&service.Id)
	return service, err
}

func (r *serviceRepo) GetAll() ([]models.Service, error) {
	query := `SELECT id, name, description, price, duration, clinic_id, is_active FROM services`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.Id, &s.Name, &s.Description, &s.Price, &s.Duration, &s.ClinicID, &s.IsActive); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *serviceRepo) GetByClinicID(clinicID string) ([]models.Service, error) {
	query := `SELECT id, name, description, price, duration, clinic_id, is_active FROM services WHERE clinic_id = $1`

	rows, err := r.db.Query(context.Background(), query, clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.Id, &s.Name, &s.Description, &s.Price, &s.Duration, &s.ClinicID, &s.IsActive); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *serviceRepo) GetByID(id string) (*models.Service, error) {
	query := `SELECT id, name, description, price, duration, clinic_id, is_active FROM services WHERE id = $1`
	var s models.Service
	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&s.Id, &s.Name, &s.Description, &s.Price, &s.Duration, &s.ClinicID, &s.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *serviceRepo) Update(id string, service *models.Service) (*models.Service, error) {
	query := `
		UPDATE services
		SET name=$1, description=$2, price=$3, duration=$4, clinic_id=$5, is_active=$6
		WHERE id=$7
		RETURNING id, name, description, price, duration, clinic_id, is_active
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		service.Name,
		service.Description,
		service.Price,
		service.Duration,
		service.ClinicID,
		service.IsActive,
		id,
	).Scan(&service.Id, &service.Name, &service.Description, &service.Price, &service.Duration, &service.ClinicID, &service.IsActive)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return service, nil
}

func (r *serviceRepo) Delete(id string) error {
	query := `DELETE FROM services WHERE id=$1`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}