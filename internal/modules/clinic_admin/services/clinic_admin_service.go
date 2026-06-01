package services

import (
	"errors"
	"time"

	"dental_clinic/internal/modules/clinic_admin/dto"
	"dental_clinic/internal/modules/clinic_admin/models"
	"dental_clinic/internal/modules/clinic_admin/repository"
	userDto "dental_clinic/internal/modules/user/dto"
	userServices "dental_clinic/internal/modules/user/services"

	"github.com/google/uuid"
)

type ClinicAdminService struct {
	repo    repository.ClinicAdminRepository
	userSrv userServices.UserService
}

func NewClinicAdminService(repo repository.ClinicAdminRepository, userSrv userServices.UserService) *ClinicAdminService {
	return &ClinicAdminService{repo: repo, userSrv: userSrv}
}

func (s *ClinicAdminService) CreateClinicAdmin(req dto.CreateClinicAdminRequest) (*models.ClinicAdmin, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic_id")
	}

	user, err := s.userSrv.CreateUser(req.Email, req.Password, req.Name, "clinic_admin", req.IsActive)
	if err != nil {
		return nil, err
	}

	admin, err := s.repo.Create(&models.ClinicAdmin{
		Id:        uuid.New(),
		ClinicID:  clinicID,
		UserID:    user.Id,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	})
	if err != nil {
		_ = s.userSrv.DeleteUserById(user.Id.String())
		return nil, err
	}
	admin.Name = req.Name
	admin.Email = req.Email
	return admin, nil
}

func (s *ClinicAdminService) GetAllClinicAdmins() ([]models.ClinicAdmin, error) {
	return s.repo.GetAll()
}

func (s *ClinicAdminService) GetClinicAdminByID(id string) (*models.ClinicAdmin, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid clinic_admin id")
	}
	admin, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("clinic admin not found")
	}
	return admin, nil
}

func (s *ClinicAdminService) UpdateClinicAdmin(id string, req dto.UpdateClinicAdminRequest) (*models.ClinicAdmin, error) {
	adminID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid clinic_admin id")
	}
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic_id")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	admin, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("clinic admin not found")
	}

	_, err = s.userSrv.UpdateUser(admin.UserID.String(), userDto.RegisterRequest{
		Role:  "clinic_admin",
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		return nil, err
	}
	if req.NewPassword != "" {
		if err := s.userSrv.UpdatePasswordWithUserId(admin.UserID.String(), req.NewPassword); err != nil {
			return nil, err
		}
	}
	if err := s.userSrv.UpdateUserVerification(admin.UserID.String(), req.IsActive); err != nil {
		return nil, err
	}

	admin.Id = adminID
	admin.ClinicID = clinicID
	admin.Name = req.Name
	admin.Email = req.Email
	admin, err = s.repo.Update(admin)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("clinic admin not found")
	}
	admin.Name = req.Name
	admin.Email = req.Email
	return admin, nil
}

func (s *ClinicAdminService) DeleteClinicAdmin(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid clinic_admin id")
	}
	admin, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if admin == nil {
		return errors.New("clinic admin not found")
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return s.userSrv.DeleteUserById(admin.UserID.String())
}

func ToClinicAdminResponse(admin models.ClinicAdmin) dto.ClinicAdminResponse {
	return dto.ClinicAdminResponse{
		Id:        admin.Id.String(),
		ClinicID:  admin.ClinicID.String(),
		UserID:    admin.UserID.String(),
		Name:      admin.Name,
		Email:     admin.Email,
		CreatedAt: admin.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToClinicAdminResponseList(admins []models.ClinicAdmin) []dto.ClinicAdminResponse {
	result := make([]dto.ClinicAdminResponse, 0, len(admins))
	for _, admin := range admins {
		result = append(result, ToClinicAdminResponse(admin))
	}
	return result
}
