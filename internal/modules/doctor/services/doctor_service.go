package services

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/utils"
	"errors"
	"fmt"
	"math/rand"

	"dental_clinic/internal/modules/doctor/dto"
	"dental_clinic/internal/modules/doctor/models"
	"dental_clinic/internal/modules/doctor/repository"

	medical_recordServices "dental_clinic/internal/modules/medical_record/services"
	userServices "dental_clinic/internal/modules/user/services"

	"github.com/google/uuid"
)

type DoctorService struct {
	repo              repository.DoctorRepository
	userSrv           userServices.UserService
	medical_recordSrv medical_recordServices.MedicalRecordService
	cfg               config.Config
}

func NewDoctorService(r repository.DoctorRepository, userSrv userServices.UserService, medical_recordSrv medical_recordServices.MedicalRecordService, cfg config.Config) *DoctorService {
	return &DoctorService{
		repo:              r,
		userSrv:           userSrv,
		medical_recordSrv: medical_recordSrv,
		cfg:               cfg,
	}
}
func generateConfirmationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

type CreateDoctorResult struct {
	Doctor           *models.Doctor
	ConfirmationCode string
}

func (s *DoctorService) CreateDoctor(req dto.CreateDoctorRequest) (*CreateDoctorResult, error) {
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

	confirmationCode := generateConfirmationCode()
	_ = utils.SendDoctorWelcomeEmail(&s.cfg, doctor.Email, doctor.Name, confirmationCode)

	return &CreateDoctorResult{
		Doctor:           doctor,
		ConfirmationCode: confirmationCode,
	}, nil
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
	err = s.repo.Delete(id)
	if err != nil {
		return err
	}
	return s.userSrv.DeleteUserById(doctor.UserId.String())
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

func (s *DoctorService) GetDoctorByIdMedicalRecords(id string) ([]dto.GetMedicalRecordDoctorResponse, error) {
	doctor, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}
	medical_records, err := s.medical_recordSrv.GetMedicalRecordsByDoctorId(id)
	if err != nil {
		return nil, err
	}

	var responses []dto.GetMedicalRecordDoctorResponse
	for _, medical_record := range medical_records {

		response := dto.GetMedicalRecordDoctorResponse{
			Diagnosis:  medical_record.Diagnosis,
			Notes:      medical_record.Notes,
			Is_checked: medical_record.Is_checked,
		}
		responses = append(responses, response)
	}
	return responses, nil
}
func (s *DoctorService) GetDoctorByUserIdMedicalRecords(id string) ([]dto.GetMedicalRecordDoctorResponse, error) {
	doctor, err := s.repo.GetByUserID(id)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}
	medical_records, err := s.medical_recordSrv.GetMedicalRecordsByDoctorId(id)
	if err != nil {
		return nil, err
	}

	var responses []dto.GetMedicalRecordDoctorResponse
	for _, medical_record := range medical_records {

		response := dto.GetMedicalRecordDoctorResponse{
			Diagnosis:  medical_record.Diagnosis,
			Notes:      medical_record.Notes,
			Is_checked: medical_record.Is_checked,
		}
		responses = append(responses, response)
	}
	return responses, nil
}
