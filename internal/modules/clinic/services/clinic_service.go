package services

import (
	"errors"
	"fmt"

	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/address/services"
	"dental_clinic/internal/modules/clinic/dto"
	"dental_clinic/internal/modules/clinic/models"
	"dental_clinic/internal/modules/clinic/repository"

	"github.com/google/uuid"
)

type ClinicService struct {
	repo       repository.ClinicRepository
	cfx        config.Config
	addressSrv services.AddressService
}

func NewClinicService(r repository.ClinicRepository, cfx config.Config, addressSrv services.AddressService) *ClinicService {
	return &ClinicService{
		repo:       r,
		cfx:        cfx,
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
	clinic.LogoURL = existing.LogoURL

	return s.repo.Update(clinic)
}

func (s *ClinicService) DeleteClinic(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("clinic not found: %w", err)
	}

	return s.repo.Delete(id)
}

func (s *ClinicService) AddAddress(id uuid.UUID, req dto.AddAddressRequest) error {
	_, err := s.addressSrv.GetAddressByID(req.Address_id)
	if err != nil {
		return fmt.Errorf("address not found: %w", err)
	}

	return s.repo.AddAddress(uuid.New(), id, req.Address_id, req.Is_main)
}

func (s *ClinicService) GetClinicAddress(id uuid.UUID) ([]models.ClinicAddress, error) {
	return s.repo.GetClinicAddress(id)
}

func (s *ClinicService) GetClinicAddressByID(id uuid.UUID) (*models.ClinicAddress, error) {
	return s.repo.GetClinicAddressByID(id)
}

func ToClinicAddressResponse(clinicAddress models.ClinicAddressWithNames) dto.GetClinicAddressResponse {
	gallery := make([]dto.ClinicAddressImageResponse, 0, len(clinicAddress.Gallery))
	for _, image := range clinicAddress.Gallery {
		gallery = append(gallery, dto.ClinicAddressImageResponse{
			Id:       image.Id.String(),
			ImageURL: image.ImageURL,
		})
	}

	return dto.GetClinicAddressResponse{
		Id:               clinicAddress.Id.String(),
		Address_id:       clinicAddress.AddressId.String(),
		Is_main:          clinicAddress.IsMain,
		Address_name:     clinicAddress.AddressName,
		Address_building: clinicAddress.AddressBuilding,
		CoverImageURL:    clinicAddress.CoverImageURL,
		Gallery:          gallery,
	}
}

func ToClinicAddressResponseList(clinicAddress []models.ClinicAddressWithNames) []dto.GetClinicAddressResponse {
	result := make([]dto.GetClinicAddressResponse, 0, len(clinicAddress))
	for _, u := range clinicAddress {
		result = append(result, ToClinicAddressResponse(u))
	}
	return result
}

func (s *ClinicService) DeleteAddress(id, address_id uuid.UUID) error {
	_, err := s.addressSrv.GetAddressByID(address_id.String())
	if err != nil {
		return fmt.Errorf("address not found: %w", err)
	}

	return s.repo.DeleteAddress(id, address_id)
}

func (s *ClinicService) UpdateClinicLogo(id uuid.UUID, logoURL string) error {
	if logoURL == "" {
		return errors.New("logo_url is required")
	}
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("clinic not found: %w", err)
	}
	return s.repo.UpdateLogo(id, logoURL)
}

func (s *ClinicService) DeleteClinicLogo(id uuid.UUID) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("clinic not found: %w", err)
	}
	return s.repo.DeleteLogo(id)
}

func (s *ClinicService) UpdateAddressCover(id uuid.UUID, coverURL string) error {
	if coverURL == "" {
		return errors.New("cover_image_url is required")
	}
	if _, err := s.repo.GetClinicByAddressId(id); err != nil {
		return fmt.Errorf("clinic address not found: %w", err)
	}
	return s.repo.UpdateAddressCover(id, coverURL)
}

func (s *ClinicService) DeleteAddressCover(id uuid.UUID) error {
	if _, err := s.repo.GetClinicByAddressId(id); err != nil {
		return fmt.Errorf("clinic address not found: %w", err)
	}
	return s.repo.DeleteAddressCover(id)
}

func (s *ClinicService) AddGalleryImage(clinicAddressID uuid.UUID, imageURL string) (*models.ClinicAddressGalleryImage, error) {
	if imageURL == "" {
		return nil, errors.New("image_url is required")
	}
	if _, err := s.repo.GetClinicByAddressId(clinicAddressID); err != nil {
		return nil, fmt.Errorf("clinic address not found: %w", err)
	}

	image := models.ClinicAddressGalleryImage{
		Id:              uuid.New(),
		ClinicAddressId: clinicAddressID,
		ImageURL:        imageURL,
	}
	if err := s.repo.AddGalleryImage(image); err != nil {
		return nil, err
	}
	return &image, nil
}

func (s *ClinicService) GetGalleryImages(clinicAddressID uuid.UUID) ([]models.ClinicAddressGalleryImage, error) {
	if _, err := s.repo.GetClinicByAddressId(clinicAddressID); err != nil {
		return nil, fmt.Errorf("clinic address not found: %w", err)
	}
	return s.repo.GetGalleryImages(clinicAddressID)
}

func (s *ClinicService) GetGalleryImage(id uuid.UUID) (*models.ClinicAddressGalleryImage, error) {
	return s.repo.GetGalleryImage(id)
}

func (s *ClinicService) UpdateGalleryImage(id uuid.UUID, imageURL string) (*models.ClinicAddressGalleryImage, error) {
	if imageURL == "" {
		return nil, errors.New("image_url is required")
	}
	currentImage, err := s.repo.GetGalleryImage(id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateGalleryImage(id, imageURL); err != nil {
		return nil, err
	}
	currentImage.ImageURL = imageURL
	return currentImage, nil
}

func (s *ClinicService) DeleteGalleryImage(id uuid.UUID) (*models.ClinicAddressGalleryImage, error) {
	image, err := s.repo.GetGalleryImage(id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.DeleteGalleryImage(id); err != nil {
		return nil, err
	}
	return image, nil
}

func (s *ClinicService) GetClinicAddressWithName(id uuid.UUID) ([]models.ClinicAddressWithNames, error) {
	clinics, err := s.repo.GetClinicAddress(id)

	if err != nil {
		return nil, err
	}

	var clinics_with_name []models.ClinicAddressWithNames

	for _, clinic := range clinics {

		clinic_with_name := models.ClinicAddressWithNames{}

		clinic_with_name.Id = clinic.Id
		clinic_with_name.ClinicId = clinic.ClinicId
		clinic_with_name.AddressId = clinic.AddressId
		clinic_with_name.IsMain = clinic.IsMain
		clinic_with_name.CoverImageURL = clinic.CoverImageURL

		address, err := s.addressSrv.GetAddressByID(clinic_with_name.AddressId.String())
		if err != nil {
			return nil, fmt.Errorf("failed to scan clinic: %w", err)
		}
		clinic_with_name.AddressName = address.Street
		clinic_with_name.AddressBuilding = address.Building
		clinic_with_name.Gallery, err = s.repo.GetGalleryImages(clinic_with_name.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get clinic address gallery: %w", err)
		}

		clinics_with_name = append(clinics_with_name, clinic_with_name)

	}

	return clinics_with_name, nil

}

func (s *ClinicService) GetClinicByAddressId(id uuid.UUID) (string, error) {
	return s.repo.GetClinicByAddressId(id)
}
