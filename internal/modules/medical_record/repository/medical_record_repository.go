package repository

import (
	"context"
	"dental_clinic/internal/modules/medical_record/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalRecordRepository interface {
	Create(medical_record *models.MedicalRecord) (*models.MedicalRecord, error)
	GetByID(id string) (*models.MedicalRecord, error)
	GetMedicalRecordByAppointmentId(id string) (*models.MedicalRecord, error)
	GetMedicalRecordsByDoctorId(id string) ([]models.MedicalRecord, error)
	//GetAll() ([]models.Doctor, error)
	Update(id string, doctor *models.MedicalRecord) (*models.MedicalRecord, error)
	//Delete(id string) error
	SaveMedicalFile(medicalRecordID, fileURL string) error
}

type medical_report_Repo struct {
	db *pgxpool.Pool
}

func NewMedicalRecordRepository(db *pgxpool.Pool) MedicalRecordRepository {
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

func (r *medical_report_Repo) GetByID(id string) (*models.MedicalRecord, error) {
	query := `SELECT appointment_id, doctor_id, patient_id, diagnosis, notes, is_checked, created_at, updated_at FROM medical_records WHERE id = $1`
	var medical_record models.MedicalRecord
	err := r.db.QueryRow(context.Background(), query, id).Scan(&medical_record.Appointment_id, &medical_record.Doctor_id, &medical_record.Patient_id, &medical_record.Diagnosis, &medical_record.Notes, &medical_record.Is_checked, &medical_record.Created_at, &medical_record.Updated_at)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &medical_record, nil
}

func (r *medical_report_Repo) Update(id string, medical_record *models.MedicalRecord) (*models.MedicalRecord, error) {
	query := `
		UPDATE medical_records
		SET diagnosis=$1, notes=$2, is_checked=$3, updated_at=$4
		WHERE id=$5
		RETURNING id, diagnosis, notes, is_checked, updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		medical_record.Diagnosis,
		medical_record.Notes,
		medical_record.Is_checked,
		medical_record.Updated_at,
		id,
	).Scan(
		&medical_record.Id,
		&medical_record.Diagnosis,
		&medical_record.Notes,
		&medical_record.Is_checked,
		&medical_record.Updated_at,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return medical_record, nil
}

func (r *medical_report_Repo) GetMedicalRecordByAppointmentId(id string) (*models.MedicalRecord, error) {
	query := `SELECT appointment_id, doctor_id, patient_id, diagnosis, notes, is_checked, created_at, updated_at FROM medical_records WHERE appointment_id = $1`
	var medical_record models.MedicalRecord
	err := r.db.QueryRow(context.Background(), query, id).Scan(&medical_record.Appointment_id, &medical_record.Doctor_id, &medical_record.Patient_id, &medical_record.Diagnosis, &medical_record.Notes, &medical_record.Is_checked, &medical_record.Created_at, &medical_record.Updated_at)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &medical_record, nil
}

func (r *medical_report_Repo) GetMedicalRecordsByDoctorId(id string) ([]models.MedicalRecord, error) {
	query := `SELECT id, appointment_id, doctor_id, patient_id, diagnosis, notes, is_checked, created_at, updated_at FROM medical_records WHERE doctor_id = $1`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medical_records []models.MedicalRecord
	for rows.Next() {
		var medical_record models.MedicalRecord
		if err := rows.Scan(&medical_record.Id, &medical_record.Appointment_id, &medical_record.Doctor_id, &medical_record.Patient_id, &medical_record.Diagnosis, &medical_record.Notes, &medical_record.Is_checked, &medical_record.Created_at, &medical_record.Updated_at); err != nil {
			return nil, err
		}
		medical_records = append(medical_records, medical_record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return medical_records, nil
}

func (r *medical_report_Repo) SaveMedicalFile(medicalRecordID, fileURL string) error {
	query := `INSERT INTO medical_files (id, medical_record_id, file_url, created_at) 
              VALUES (gen_random_uuid(), $1, $2, NOW())`
	_, err := r.db.Exec(context.Background(), query, medicalRecordID, fileURL)
	return err
}
