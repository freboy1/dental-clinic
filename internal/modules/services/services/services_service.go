package services

import (
	"dental_clinic/internal/modules/clinic/services"
	"errors"

	"dental_clinic/internal/modules/services/dto"
	"dental_clinic/internal/modules/services/models"
	"dental_clinic/internal/modules/services/repository"

	"github.com/google/uuid"
)

type ServiceService struct {
	repo      repository.ServiceRepository
	clinicSrv services.ClinicService
}

func NewServiceService(r repository.ServiceRepository, clinicSrv services.ClinicService) *ServiceService {
	return &ServiceService{
		repo:      r,
		clinicSrv: clinicSrv,
	}
}

func (s *ServiceService) CreateService(req dto.CreateServiceRequest) (*models.Service, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	service := &models.Service{
		Name:        req.Name,
		Description: req.Description,
	}

	return s.repo.Create(service)
}

func (s *ServiceService) GetAllServices() ([]models.Service, error) {
	return s.repo.GetAll()
}

func (s *ServiceService) GetClinicNames(services []models.Clinic_Service) ([]models.ServiceWithClinicName, error) {
	var servicesListWithClinicNames []models.ServiceWithClinicName

	for _, service := range services {

		clinic, err := s.clinicSrv.GetClinicByID(service.ClinicID)
		if err != nil {
			return servicesListWithClinicNames, err
		}

		service_info, err := s.GetServiceByID(service.ServiceID.String())

		if err != nil {
			return servicesListWithClinicNames, err
		}

		serviceWithClinicName := models.ServiceWithClinicName{
			Id:          service.Id,
			Name:        service_info.Name,
			Description: service_info.Description,
			Price:       service.Price,
			Duration:    service.Duration,
			ClinicID:    service.ClinicID,
			IsActive:    service.IsActive,
			ClinicName:  clinic.Name,
		}

		servicesListWithClinicNames = append(servicesListWithClinicNames, serviceWithClinicName)
	}

	return servicesListWithClinicNames, nil
}

func (s *ServiceService) GetServicesByClinic(clinicID string) ([]models.ServiceWithClinicName, error) {
	if _, err := uuid.Parse(clinicID); err != nil {
		return nil, errors.New("invalid clinic_id")
	}
	clinic_services, err := s.repo.GetByClinicID(clinicID)
	if err != nil {
		return nil, err
	}
	return s.GetClinicNames(clinic_services)
}

func (s *ServiceService) GetServiceByID(id string) (*models.Service, error) {
	service, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("service not found")
	}
	return service, nil
}

func (s *ServiceService) UpdateService(id string, req dto.UpdateServiceRequest) (*models.Service, error) {
	service, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("service not found")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	service.Name = req.Name
	service.Description = req.Description

	return s.repo.Update(id, service)
}

func (s *ServiceService) DeleteService(id string) error {
	service, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if service == nil {
		return errors.New("service not found")
	}
	return s.repo.Delete(id)
}

func ToServiceResponse(s models.Service) dto.ServiceResponse {
	return dto.ServiceResponse{
		Id:          s.Id.String(),
		Name:        s.Name,
		Description: s.Description,
		// ClinicName: ,
	}
}

func ToServiceResponseList(services []models.Service) []dto.ServiceResponse {
	result := make([]dto.ServiceResponse, 0, len(services))
	for _, s := range services {
		result = append(result, ToServiceResponse(s))
	}
	return result
}

func ToServiceNameResponse(s models.ServiceWithClinicName) dto.ServiceResponseWithName {
	return dto.ServiceResponseWithName{
		Id:          s.Id.String(),
		Name:        s.Name,
		Description: s.Description,
		Price:       s.Price,
		Duration:    s.Duration,
		ClinicID:    s.ClinicID.String(),
		IsActive:    s.IsActive,
		ClinicName:  s.ClinicName,
	}
}

func ToServiceNameResponseList(services []models.ServiceWithClinicName) []dto.ServiceResponseWithName {
	result := make([]dto.ServiceResponseWithName, 0, len(services))
	for _, s := range services {
		result = append(result, ToServiceNameResponse(s))
	}
	return result
}

func (s *ServiceService) GetServices() ([]models.Service, error) {
	return s.repo.GetAll()
}

func (s *ServiceService) AddServiceToClinic(id string, req dto.AddServiceRequest) (*models.Clinic_Service, error) {
	if req.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if req.Duration <= 0 {
		return nil, errors.New("duration must be greater than 0")
	}

	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, errors.New("invalid service_id")
	}

	clinicID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid clinic_id")
	}

	service := &models.Clinic_Service{
		Id:        uuid.New(),
		ClinicID:  clinicID,
		Price:     req.Price,
		Duration:  req.Duration,
		ServiceID: serviceID,
		IsActive:  req.IsActive,
	}

	return s.repo.AddServiceToClinic(service)
}

func (s *ServiceService) DeleteServiceByClinic(clinicID, serviceID string) error {
	if _, err := uuid.Parse(clinicID); err != nil {
		return errors.New("invalid clinic_id")
	}

	if _, err := uuid.Parse(serviceID); err != nil {
		return errors.New("invalid service_id")
	}

	return s.repo.DeleteServiceToClinic(clinicID, serviceID)
}

func (s *ServiceService) GetByClinicIDAndServiceID(clinicID, serviceID string) (*models.Clinic_Service, error) {
	if _, err := uuid.Parse(clinicID); err != nil {
		return nil, errors.New("invalid clinic_id")
	}

	if _, err := uuid.Parse(serviceID); err != nil {
		return nil, errors.New("invalid service_id")
	}

	return s.repo.GetByClinicIDAndServiceID(clinicID, serviceID)
}
