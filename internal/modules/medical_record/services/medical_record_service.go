package services

import (
	//"dental_clinic/internal/modules/medical_record/dto"
	"dental_clinic/internal/modules/medical_record/models"
	"dental_clinic/internal/modules/medical_record/repository"
	"time"

	//"errors"
	//
	//"github.com/google/uuid"

	"github.com/google/uuid"
)

type MedicalRecordService struct {
	repo repository.MedicalRecordRepository
}

func NewMedicalRecordService(r repository.MedicalRecordRepository) *MedicalRecordService {
	return &MedicalRecordService{repo: r}
}

func (s *MedicalRecordService) CreateMedicalRecord(appointment_id, doctor_id, patient_id uuid.UUID) (*models.MedicalRecord, error) {
	medical_record := &models.MedicalRecord{
		Id:             uuid.New(),
		Appointment_id: appointment_id,
		Doctor_id:      doctor_id,
		Patient_id:     patient_id,
		Diagnosis:      "",
		Notes:          "",
		Is_checked:     false,
		Created_at:     time.Now(),
		Updated_at:     time.Now(),
	}

	return s.repo.Create(medical_record)
}
