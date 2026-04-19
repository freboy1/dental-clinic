package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalRecordRepository interface {
	//Create(medical_record *models.MedicalRecord) (*models.MedicalRecord, error)
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
