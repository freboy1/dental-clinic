package services

import (
	"errors"

	"dental_clinic/internal/modules/doctor/dto"
	"dental_clinic/internal/modules/doctor/models"
	"dental_clinic/internal/modules/doctor/repository"

	userServices "dental_clinic/internal/modules/user/services"

	"github.com/google/uuid"
)

type DoctorService struct {
	repo    repository.DoctorRepository
	userSrv userServices.UserService
}

func NewDoctorService(r repository.DoctorRepository, userSrv userServices.UserService) *DoctorService {
	return &DoctorService{
		repo:    r,
		userSrv: userSrv,
	}
}

func (s *DoctorService) CreateDoctor(req dto.CreateDoctorRequest) (*models.Doctor, error) {
	if req.Specialization == "" {
		return nil, errors.New("specialization is required")
	}
	if req.Experience < 0 {
		return nil, errors.New("experience cannot be negative")
	}
	if req.Name == "" {
		return nil, errors.New("name is empty")
	}
	if req.Email == "" {
		return nil, errors.New("email is empty")
	}

	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic_id")
	}

	doctor := &models.Doctor{
		Name:           req.Name,
		Email:          req.Email,
		Specialization: req.Specialization,
		Experience:     req.Experience,
		ClinicID:       clinicID,
		Bio:            req.Bio,
		IsAvailable:    req.IsAvailable,
	}

	user, err := s.userSrv.CreateUser(doctor.Email, req.Password, doctor.Name, "doctor", req.Is_active)
	if err != nil {
		return nil, err
	}

	doctor.UserId = user.Id

	doctor, err = s.repo.Create(doctor)
	if err != nil {
		return nil, err
	}

	return doctor, nil
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

	err = s.userSrv.UpdatePasswordWithUserId(doctor.UserId.String(), req.NewPassword)
	if err != nil {
		return nil, err
	}

	err = s.userSrv.UpdateUserVerification(doctor.UserId.String(), req.Is_active)
	if err != nil {
		return nil, err
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
		Name:           d.Name,
		Email:          d.Email,
	}
}

func ToDoctorResponseList(doctors []models.Doctor) []dto.DoctorResponse {
	result := make([]dto.DoctorResponse, 0, len(doctors))
	for _, d := range doctors {
		result = append(result, ToDoctorResponse(d))
	}
	return result
}
