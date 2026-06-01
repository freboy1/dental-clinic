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
	GetByUserID(id string) (*models.Doctor, error)
	GetAll() ([]models.Doctor, error)
	Update(id string, doctor *models.Doctor) (*models.Doctor, error)
	Delete(id string) error
	UpdatePhoto(id, photoURL string) error
	DeletePhoto(id string) error
}

type doctorRepo struct {
	db *pgxpool.Pool
}

func NewDoctorRepository(db *pgxpool.Pool) DoctorRepository {
	return &doctorRepo{db: db}
}

func (r *doctorRepo) Create(doctor *models.Doctor) (*models.Doctor, error) {
	query := `
		INSERT INTO doctors (specialization, experience, clinic_id, bio, is_available, name, email, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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
		doctor.UserId,
	).Scan(&doctor.Id)
	return doctor, err
}

func (r *doctorRepo) GetAll() ([]models.Doctor, error) {
	query := `
		SELECT
			d.id,
			d.specialization,
			d.experience,
			d.clinic_id,
			d.bio,
			d.is_available,
			d.name,
			d.email,
			COALESCE(d.photo_url, ''),
			COALESCE(ROUND(AVG(dr.rating)::numeric, 2), 0)::float8 AS rating
		FROM doctors d
		LEFT JOIN doctor_ratings dr ON dr.doctor_id = d.id
		WHERE d.is_deleted=0
		GROUP BY d.id, d.specialization, d.experience, d.clinic_id, d.bio, d.is_available, d.name, d.email, d.photo_url
		ORDER BY rating DESC, d.name
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var d models.Doctor
		if err := rows.Scan(&d.Id, &d.Specialization, &d.Experience, &d.ClinicID, &d.Bio, &d.IsAvailable, &d.Name, &d.Email, &d.PhotoURL, &d.Rating); err != nil {
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
	query := `
		SELECT
			d.id,
			d.specialization,
			d.experience,
			d.clinic_id,
			d.bio,
			d.is_available,
			d.name,
			d.email,
			d.user_id,
			COALESCE(d.photo_url, ''),
			COALESCE(ROUND(AVG(dr.rating)::numeric, 2), 0)::float8 AS rating
		FROM doctors d
		LEFT JOIN doctor_ratings dr ON dr.doctor_id = d.id
		WHERE d.id = $1 AND d.is_deleted=0
		GROUP BY d.id, d.specialization, d.experience, d.clinic_id, d.bio, d.is_available, d.name, d.email, d.user_id, d.photo_url
	`
	var d models.Doctor
	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&d.Id, &d.Specialization, &d.Experience, &d.ClinicID, &d.Bio, &d.IsAvailable, &d.Name, &d.Email, &d.UserId, &d.PhotoURL, &d.Rating)
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

func (r *doctorRepo) GetByUserID(id string) (*models.Doctor, error) {
	query := `
		SELECT
			d.id,
			d.specialization,
			d.experience,
			d.clinic_id,
			d.bio,
			d.is_available,
			d.name,
			d.email,
			d.user_id,
			COALESCE(d.photo_url, ''),
			COALESCE(ROUND(AVG(dr.rating)::numeric, 2), 0)::float8 AS rating
		FROM doctors d
		LEFT JOIN doctor_ratings dr ON dr.doctor_id = d.id
		WHERE d.user_id = $1 AND d.is_deleted=0
		GROUP BY d.id, d.specialization, d.experience, d.clinic_id, d.bio, d.is_available, d.name, d.email, d.user_id, d.photo_url
	`
	var d models.Doctor
	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&d.Id, &d.Specialization, &d.Experience, &d.ClinicID, &d.Bio, &d.IsAvailable, &d.Name, &d.Email, &d.UserId, &d.PhotoURL, &d.Rating)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *doctorRepo) UpdatePhoto(id, photoURL string) error {
	result, err := r.db.Exec(context.Background(), `UPDATE doctors SET photo_url = $2 WHERE id = $1 AND is_deleted = 0`, id, photoURL)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *doctorRepo) DeletePhoto(id string) error {
	result, err := r.db.Exec(context.Background(), `UPDATE doctors SET photo_url = NULL WHERE id = $1 AND is_deleted = 0`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
