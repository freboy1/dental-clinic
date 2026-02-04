package services

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/clinic/models"
	"dental_clinic/internal/modules/clinic/repository"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ClinicService struct {
	repo repository.ClinicRepository
	cfx  config.Config
}

func NewClinicService(r repository.ClinicRepository, cfx config.Config) *ClinicService {
	return &ClinicService{
		repo: r,
		cfx:  cfx,
	}
}

func (s *ClinicService) CreateClinic(clinic *models.Clinic) (*models.Clinic, error) {
	if clinic.Name == "" {
		return nil, fmt.Errorf("clinic name is required")
	}

	if clinic.Phone == "" {
		return nil, fmt.Errorf("clinic phone is required")
	}

	clinic.Id = uuid.New()

	clinic.CreatedAt = time.Now()

	clinic.IsActive = true

	return s.repo.Create(clinic)
}

func (s *ClinicService) GetAllClinics() ([]*models.Clinic, error) {
	return s.repo.GetAll()
}

func (s *ClinicService) GetClinicByID(id uuid.UUID) (*models.Clinic, error) {
	return s.repo.GetByID(id)
}

func (s *ClinicService) UpdateClinic(id uuid.UUID, clinic *models.Clinic) (*models.Clinic, error) {
	if clinic.Name == "" {
		return nil, fmt.Errorf("clinic name is required")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("clinic not found: %w", err)
	}

	clinic.Id = id
	clinic.CreatedAt = existing.CreatedAt

	return s.repo.Update(clinic)
}

func (s *ClinicService) DeleteClinic(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("clinic not found: %w", err)
	}

	return s.repo.Delete(id)
}
