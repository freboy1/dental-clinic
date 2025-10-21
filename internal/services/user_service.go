package services

import (
	"dental_clinic/internal/models"
	"dental_clinic/internal/repository"

	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

type RegisterRequest struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	PushConsent bool   `json:"push_consent"`
}

func (s *UserService) Register(req RegisterRequest) (*models.User, error) {
	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}
	if !isValidPassword(req.Password) {
		return nil, errors.New("weak password")
	}

	// add check for existing user

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &models.User{
		Role:         req.Role,
		Email:        req.Email,
		Password:     string(hash),
		Name:         req.Name,
		Gender:       req.Gender,
		Age:          req.Age,
		Push_consent: req.PushConsent,
	}

	return s.repo.Create(user)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateUser(id string, req RegisterRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	user.Gender = req.Gender
	user.Age = req.Age
	user.Push_consent = req.PushConsent

	return s.repo.Update(id, user)
}

func (s *UserService) DeleteUser(id string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.repo.Delete(id)
}

func isValidEmail(email string) bool {
	// do it a bit later
	return true
}

func isValidPassword(password string) bool {
	// do it a bit later
	return true
}
