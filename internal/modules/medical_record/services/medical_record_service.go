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

func (s *MedicalRecordService) UpdateMedicalRecord(id string, req dto.UpdateMedicalRecordRequest, medicalFiles []models.MedicalFile) (*models.MedicalRecord, error) {
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

	updated, err := s.repo.Update(id, medical_record)
	if err != nil {
		return nil, err
	}

	// сохраняем пути файлов в БД
	for _, medicalFile := range medicalFiles {
		_ = s.repo.SaveMedicalFile(id, medicalFile.FilePath, medicalFile.Filename, medicalFile.MimeType)
	}

	return updated, nil
}

func (s *MedicalRecordService) GetMedicalRecord(id string) (*models.MedicalRecord, error) {
	return s.repo.GetByID(id)
}
func (s *MedicalRecordService) GetMedicalRecordByAppointmentId(id string) (*models.MedicalRecord, error) {
	return s.repo.GetMedicalRecordByAppointmentId(id)
}

func (s *MedicalRecordService) GetMedicalRecordsByDoctorId(id string) ([]models.MedicalRecord, error) {
	return s.repo.GetMedicalRecordsByDoctorId(id)
}

func (s *MedicalRecordService) GetMedicalRecordFiles(id string) ([]dto.MedicalFileResponse, error) {
	medical_files, err := s.repo.GetMedicalFiles(id)
	if err != nil {
		return nil, err
	}

	return ToMedicalFileResponseList(medical_files), nil
}

func ToMedicalFileResponse(s models.MedicalFile) dto.MedicalFileResponse {
	return dto.MedicalFileResponse{
		ID:       s.Id.String(),
		Name:     s.Filename,
		MimeType: s.MimeType,
		// ClinicName: ,
	}
}

func ToMedicalFileResponseList(services []models.MedicalFile) []dto.MedicalFileResponse {
	result := make([]dto.MedicalFileResponse, 0, len(services))
	for _, s := range services {
		result = append(result, ToMedicalFileResponse(s))
	}
	return result
}

func (s *MedicalRecordService) GetFileByID(id string) (*models.MedicalFile, error) {
	return s.repo.GetFileByID(id)
}
