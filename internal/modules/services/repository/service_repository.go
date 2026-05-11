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
	GetByClinicID(clinicID string) ([]models.Clinic_Service, error)
	Update(id string, service *models.Service) (*models.Service, error)
	Delete(id string) error
	AddServiceToClinic(clinic_service *models.Clinic_Service) (*models.Clinic_Service, error)
	DeleteServiceToClinic(clinicID, serviceID string) error
}

type serviceRepo struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) ServiceRepository {
	return &serviceRepo{db: db}
}

func (r *serviceRepo) Create(service *models.Service) (*models.Service, error) {
	query := `
		INSERT INTO services (name, description)
		VALUES ($1, $2)
		RETURNING id
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		service.Name,
		service.Description,
	).Scan(&service.Id)
	return service, err
}

func (r *serviceRepo) GetAll() ([]models.Service, error) {
	query := `SELECT id, name, description FROM services`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.Id, &s.Name, &s.Description); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *serviceRepo) GetByClinicID(clinicID string) ([]models.Clinic_Service, error) {
	query := `SELECT id, clinic_id, service_id, price, duration_minutes, is_active FROM clinic_services WHERE clinic_id = $1`

	rows, err := r.db.Query(context.Background(), query, clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Clinic_Service
	for rows.Next() {
		var s models.Clinic_Service
		if err := rows.Scan(&s.Id, &s.ClinicID, &s.ServiceID, &s.Price, &s.Duration, &s.IsActive); err != nil {
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
	query := `SELECT id, name, description FROM services WHERE id = $1`
	var s models.Service
	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&s.Id, &s.Name, &s.Description)
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
		SET name=$1, description=$2
		WHERE id=$3
		RETURNING id, name, description
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		service.Name,
		service.Description,
		id,
	).Scan(&service.Id, &service.Name, &service.Description)

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

func (r *serviceRepo) AddServiceToClinic(clinic_service *models.Clinic_Service) (*models.Clinic_Service, error) {
	query := `
		INSERT INTO clinic_services (id, clinic_id, service_id, price, duration_minutes, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		clinic_service.Id,
		clinic_service.ClinicID,
		clinic_service.ServiceID,
		clinic_service.Price,
		clinic_service.Duration,
		clinic_service.IsActive,
	).Scan(&clinic_service.Id)
	return clinic_service, err
}

func (r *serviceRepo) DeleteServiceToClinic(clinicID, serviceID string) error {
	query := `DELETE FROM clinic_services WHERE clinic_id=$1 AND service_id=$2`
	result, err := r.db.Exec(context.Background(), query, clinicID, serviceID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
