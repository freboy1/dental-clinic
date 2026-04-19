package services

import (
	"dental_clinic/internal/modules/doctor/dto"
	//"dental_clinic/internal/modules/medical_record/dto"
	"dental_clinic/internal/modules/medical_record/models"
	"dental_clinic/internal/modules/medical_record/repository"
	//"errors"
	//
	//"github.com/google/uuid"
)

type MedicalRecordService struct {
	repo repository.MedicalRecordRepository
}

func NewMedicalRecordService(r repository.MedicalRecordRepository) *MedicalRecordService {
	return &MedicalRecordService{repo: r}
}

func (s *MedicalRecordService) CreateDoctor(req dto.CreateDoctorRequest) (*models.MedicalRecord, error) {
	return nil, nil
}
