package services

import (
	"dental_clinic/internal/models"
	"dental_clinic/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"errors"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

type RegisterRequest struct {
	Role    string `json:"role"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name   string `json:"name"`
	Gender       string `json:"gender"`
	Age       int `json:"age"`
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
		Role:    req.Role,
		Email:       req.Email,
		Password:    string(hash),
		Name:   req.Name,
		Gender:       req.Gender,
		Age: req.Age,
		Push_consent: req.PushConsent,
	}

	return s.repo.Create(user)
}

func isValidEmail(email string) bool {
	// do it a bit later
	return true
} 

func isValidPassword(password string) bool {
	// do it a bit later
	return true
} 