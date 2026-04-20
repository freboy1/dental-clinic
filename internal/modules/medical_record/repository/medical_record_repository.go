package repository

import (
	"context"
	"dental_clinic/internal/modules/medical_record/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalRecordRepository interface {
	Create(medical_record *models.MedicalRecord) (*models.MedicalRecord, error)
	//GetByID(id string) (*models.Doctor, error)
	//GetAll() ([]models.Doctor, error)
	//Update(id string, doctor *models.Doctor) (*models.Doctor, error)
	//Delete(id string) error
}

type medical_report_Repo struct {
	db *pgxpool.Pool
}

func NewDoctorRepository(db *pgxpool.Pool) MedicalRecordRepository {
	return &medical_report_Repo{db: db}
}

func (r *medical_report_Repo) Create(medical_record *models.MedicalRecord) (*models.MedicalRecord, error) {
	query := `
		INSERT INTO medical_records (id, appointment_id, doctor_id, patient_id, diagnosis, notes, is_checked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		medical_record.Id,
		medical_record.Appointment_id,
		medical_record.Doctor_id,
		medical_record.Patient_id,
		medical_record.Diagnosis,
		medical_record.Notes,
		medical_record.Is_checked,
		medical_record.Created_at,
		medical_record.Updated_at,
	).Scan(&medical_record.Id)
	return medical_record, err
}
