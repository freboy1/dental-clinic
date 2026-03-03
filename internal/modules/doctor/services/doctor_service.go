package services

import (
	"dental_clinic/internal/modules/doctor/dto"
	"dental_clinic/internal/modules/doctor/models"
	"dental_clinic/internal/modules/doctor/repository"
	"errors"

	"github.com/google/uuid"
)

type DoctorService struct {
	repo repository.DoctorRepository
}

func NewDoctorService(r repository.DoctorRepository) *DoctorService {
	return &DoctorService{repo: r}
}

func (s *DoctorService) CreateDoctor(req dto.CreateDoctorRequest) (*models.Doctor, error) {
	if req.Specialization == "" {
		return nil, errors.New("specialization is required")
	}
	if req.Experience < 0 {
		return nil, errors.New("experience cannot be negative")
	}

	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic_id")
	}

	doctor := &models.Doctor{
		Specialization: req.Specialization,
		Experience:     req.Experience,
		ClinicID:       clinicID,
		Bio:            req.Bio,
		IsAvailable:    req.IsAvailable,
	}

	return s.repo.Create(doctor)
}

func (s *DoctorService) GetAllDoctors() ([]models.Doctor, error) {
	return s.repo.GetAll()
}

func (s *DoctorService) GetDoctorByID(id string) (*models.Doctor, error) {
	doctor, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}
	return doctor, nil
}

func (s *DoctorService) UpdateDoctor(id string, req dto.UpdateDoctorRequest) (*models.Doctor, error) {
	doctor, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}

	if req.Specialization == "" {
		return nil, errors.New("specialization is required")
	}
	if req.Experience < 0 {
		return nil, errors.New("experience cannot be negative")
	}


	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic_id")
	}

	doctor.Specialization = req.Specialization
	doctor.Experience = req.Experience
	doctor.ClinicID = clinicID
	doctor.Bio = req.Bio
	doctor.IsAvailable = req.IsAvailable

	return s.repo.Update(id, doctor)
}

func (s *DoctorService) DeleteDoctor(id string) error {
	doctor, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if doctor == nil {
		return errors.New("doctor not found")
	}
	return s.repo.Delete(id)
}

func ToDoctorResponse(d models.Doctor) dto.DoctorResponse {
	return dto.DoctorResponse{
		Id:             d.Id.String(),
		Specialization: d.Specialization,
		Experience:     d.Experience,
		ClinicID:       d.ClinicID.String(),
		Bio:            d.Bio,
		IsAvailable:    d.IsAvailable,
	}
}

func ToDoctorResponseList(doctors []models.Doctor) []dto.DoctorResponse {
	result := make([]dto.DoctorResponse, 0, len(doctors))
	for _, d := range doctors {
		result = append(result, ToDoctorResponse(d))
	}
	return result
}