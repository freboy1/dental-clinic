package services

import (
	"dental_clinic/internal/modules/medical_record/dto"
	"errors"

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

func (s *MedicalRecordService) UpdateMedicalRecord(id string, req dto.UpdateMedicalRecordRequest) (*models.MedicalRecord, error) {
	medical_record, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if medical_record == nil {
		return nil, errors.New("medical_record not found")
	}

	medical_record.Diagnosis = req.Diagnosis
	medical_record.Notes = req.Notes
	medical_record.Is_checked = req.Is_checked
	medical_record.Updated_at = time.Now()

	return s.repo.Update(id, medical_record)
}

func (s *MedicalRecordService) GetMedicalRecord(id string) (*models.MedicalRecord, error) {
	return s.repo.GetByID(id)
}
func (s *MedicalRecordService) GetMedicalRecordByAppointmentId(id string) (*models.MedicalRecord, error) {
	return s.repo.GetMedicalRecordByAppointmentId(id)
}
