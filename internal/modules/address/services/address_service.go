package services

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/address/models"
	"dental_clinic/internal/modules/address/dto"
	"dental_clinic/internal/modules/address/repository"

	"errors"
	"fmt"

	"github.com/google/uuid"
)

type AddressService struct {
	repo repository.AddressRepository
	cfx  config.Config
}

func NewAddressService(r repository.AddressRepository, cfx config.Config) *AddressService {
	return &AddressService{
		repo: r,
		cfx:  cfx,
	}
}

func (s *AddressService) CreateAddress(req dto.CreateRequest) (*models.Address, error) {
	if req.Street == "" {
		return nil, fmt.Errorf("address stree is required")
	}

	if req.City == "" {
		return nil, fmt.Errorf("address city is required")
	}

	address := &models.Address{
		ID: uuid.New(),
		Street:         req.Street,
		City:        req.City,
		Country:     req.Country,
		Building:         req.Building,
		Latitude:       req.Latitude,
		Longitude:          req.Longitude,
	}

	return s.repo.Create(address)
}

func (s *AddressService) GetAllAddresss(tokenStr string) ([]models.Address, error) {
	return s.repo.GetAll()
}

func (s *AddressService) GetAddressByID(id string) (*models.Address, error) {
	address, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("address not found")
	}
	return address, nil
}

func (s *AddressService) UpdateAddress(id string, req dto.CreateRequest) (*models.Address, error) {
	address, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("address not found")
	}

	address.Country = req.Country
	address.City = req.City
	address.Street = req.Street
	address.Building = req.Building
	address.Latitude = req.Latitude
	address.Longitude = req.Longitude

	return s.repo.Update(id, address)
}

func (s *AddressService) DeleteAddress(id, tokenStr string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("clinic not found: %w", err)
	}

	return s.repo.Delete(id)
}


func ToAddressResponse(address models.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:     address.ID.String(),
		Country:  address.Country,
		City:   address.City,
		Street:    address.Street,
		Building: address.Building,
		Latitude:   address.Latitude,
		Longitude:   address.Longitude,
	}
}

func ToAddressResponseList(addresss []models.Address) []dto.AddressResponse {
	result := make([]dto.AddressResponse, 0, len(addresss))
	for _, u := range addresss {
		result = append(result, ToAddressResponse(u))
	}
	return result
}