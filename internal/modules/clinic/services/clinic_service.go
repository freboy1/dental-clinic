package services

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/address/services"
	"dental_clinic/internal/modules/clinic/dto"
	"dental_clinic/internal/modules/clinic/models"
	"dental_clinic/internal/modules/clinic/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ClinicService struct {
	repo repository.ClinicRepository
	cfx  config.Config
	addressSrv services.AddressService
}

func NewClinicService(r repository.ClinicRepository, cfx config.Config, addressSrv services.AddressService) *ClinicService {
	return &ClinicService{
		repo: r,
		cfx:  cfx,
		addressSrv: addressSrv,
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


func (s *ClinicService) AddAddress(id uuid.UUID, req dto.AddAddressRequest) (error) {
	_, err := s.addressSrv.GetAddressByID(req.Address_id)
	if err != nil {
		return fmt.Errorf("address not found: %w", err)
	}

	return s.repo.AddAddress(uuid.New(), id, req.Address_id, req.Is_main)
}

func (s *ClinicService) GetClinicAddress(id uuid.UUID) ([]models.ClinicAddress, error) {
	return s.repo.GetClinicAddress(id)
}

func ToClinicAddressResponse(clinicAddress models.ClinicAddress) dto.GetClinicAddressResponse {
	return dto.GetClinicAddressResponse{
		Address_id:     clinicAddress.AddressId.String(),
		Is_main:  clinicAddress.IsMain,
	}
}

func ToClinicAddressResponseList(clinicAddress []models.ClinicAddress) []dto.GetClinicAddressResponse {
	result := make([]dto.GetClinicAddressResponse, 0, len(clinicAddress))
	for _, u := range clinicAddress {
		result = append(result, ToClinicAddressResponse(u))
	}
	return result
}