package repository

import (
	"context"

	"dental_clinic/internal/modules/doctor/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DoctorRepository interface {
	Create(doctor *models.Doctor) (*models.Doctor, error)
	GetByID(id string) (*models.Doctor, error)
	GetAll() ([]models.Doctor, error)
	Update(id string, doctor *models.Doctor) (*models.Doctor, error)
	Delete(id string) error
}

type doctorRepo struct {
	db *pgxpool.Pool
}

func NewDoctorRepository(db *pgxpool.Pool) DoctorRepository {
	return &doctorRepo{db: db}
}

func (r *doctorRepo) Create(doctor *models.Doctor) (*models.Doctor, error) {
	query := `
		INSERT INTO doctors (specialization, experience, clinic_id, bio, is_available, name, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		doctor.Specialization,
		doctor.Experience,
		doctor.ClinicID,
		doctor.Bio,
		doctor.IsAvailable,
		doctor.Name,
		doctor.Email,
	).Scan(&doctor.Id)
	return doctor, err
}

func (r *doctorRepo) GetAll() ([]models.Doctor, error) {
	query := `SELECT id, specialization, experience, clinic_id, bio, is_available, name, email FROM doctors WHERE is_deleted=0`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var d models.Doctor
		if err := rows.Scan(&d.Id, &d.Specialization, &d.Experience, &d.ClinicID, &d.Bio, &d.IsAvailable, &d.Name, &d.Email); err != nil {
			return nil, err
		}
		doctors = append(doctors, d)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (r *doctorRepo) GetByID(id string) (*models.Doctor, error) {
	query := `SELECT id, specialization, experience, clinic_id, bio, is_available FROM doctors WHERE id = $1 AND is_deleted=0`
	var d models.Doctor
	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&d.Id, &d.Specialization, &d.Experience, &d.ClinicID, &d.Bio, &d.IsAvailable)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *doctorRepo) Update(id string, doctor *models.Doctor) (*models.Doctor, error) {
	query := `
		UPDATE doctors
		SET specialization=$1, experience=$2, clinic_id=$3, bio=$4, is_available=$5
		WHERE id=$6
		RETURNING id, specialization, experience, clinic_id, bio, is_available
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		doctor.Specialization,
		doctor.Experience,
		doctor.ClinicID,
		doctor.Bio,
		doctor.IsAvailable,
		id,
	).Scan(&doctor.Id, &doctor.Specialization, &doctor.Experience, &doctor.ClinicID, &doctor.Bio, &doctor.IsAvailable)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return doctor, nil
}

func (r *doctorRepo) Delete(id string) error {
	query := `
				UPDATE doctors
				SET is_deleted=1
				WHERE id=$1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
