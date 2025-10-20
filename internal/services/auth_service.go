package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/freboy1/dental-clinic/internal/models"
	"github.com/freboy1/dental-clinic/internal/repository"
	"github.com/freboy1/dental-clinic/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, user *models.User, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	loginRepo  repository.LoginHistoryRepository
	jwtManager *utils.JWTManager
}

// NewAuthService — конструктор
func NewAuthService(
	userRepo repository.UserRepository,
	loginRepo repository.LoginHistoryRepository,
	jwtManager *utils.JWTManager,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		loginRepo:  loginRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, user *models.User, password string) (string, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, user.Email)
	if existing != nil {
		return "", errors.New("user with this email already exists")
	}

	if err := utils.ValidateEmail(user.Email); err != nil {
		return "", err
	}
	if err := utils.ValidatePassword(password); err != nil {
		return "", err
	}
	if err := utils.ValidateName(user.FirstName); err != nil {
		return "", err
	}
	if err := utils.ValidateName(user.LastName); err != nil {
		return "", err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user.PasswordHash = hash
	user.Role = "user"
	user.Activated = true
	user.RegisteredAt = time.Now()

	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	token, err := s.jwtManager.Generate(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	log.Println("✅ User registered:", user.Email)
	return token, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.CheckPassword(user.PasswordHash, password) {
		_ = s.loginRepo.LogLogin(ctx, user.ID, "unknown", false, "unknown")
		return "", errors.New("invalid email or password")
	}

	_ = s.loginRepo.LogLogin(ctx, user.ID, "unknown", true, "unknown")

	token, err := s.jwtManager.Generate(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
