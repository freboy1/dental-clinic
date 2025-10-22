package services

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/models"
	"dental_clinic/internal/repository"
	"dental_clinic/internal/utils"
	"fmt"
	"regexp"

	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
	cfx  config.Config
}

func NewUserService(r repository.UserRepository, cfx config.Config) *UserService {
	return &UserService{
		repo: r,
		cfx:  cfx,
	}
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (s *UserService) Register(req RegisterRequest) (*models.User, error) {
	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}
	if !isValidPassword(req.Password) {
		return nil, errors.New("weak password")
	}
	if !isValidName(req.Name) {
		return nil, errors.New("invalid name")
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
	created_user, err := s.repo.Create(user)
	if err != nil {
		return created_user, err
	}

	token := uuid.NewString()
	err = s.repo.SaveVerificationToken(created_user.Id.String(), token)
	if err != nil {
		return created_user, err
	}
	utils.SendVerificationEmail(&s.cfx, user.Email, token)

	return created_user, err
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

func (s *UserService) DeleteUser(id, tokenStr string) error {
	claims, _ := utils.GetClaims(tokenStr, s.cfx.JWTSecret)
	if claims["role"] == "admin" || claims["user_id"].(string) == id {
		user, err := s.repo.GetByID(id)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("user not found")
		}

		return s.repo.Delete(id)
	}
	return errors.New("do not have rights")
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidPassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.MatchString(password) || len(password) < 8 {
		return false
	}
	return true
}

func isValidName(name string) bool {
	if name == "" {
		return false
	}
	re := regexp.MustCompile(`^[A-Za-zА-Яа-яЁё]+$`)
	return re.MatchString(name)
}

func (s *UserService) VerifyUserEmail(token string) error {
	userID, err := s.repo.FindUserIdByToken(token)
	if err != nil {
		return err
	}
	return s.repo.MarkUserAsVerified(userID)

}

func (s *UserService) Login(req LoginRequest, ip string) (*models.User, error) {

	// add check for existing user
	if req.Email == "" {
		s.repo.LogLogin("", ip, false)
		return nil, errors.New("email is empty")
	}
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		_ = s.repo.LogLogin("", ip, false)
		return nil, err
	}

	if user == nil {
		_ = s.repo.LogLogin("", ip, false)
		return nil, errors.New("user not found")
	}

	fmt.Println(user.Password)
	fmt.Println(req.Password)

	if !CheckPassword(user.Password, req.Password) {
		s.repo.LogLogin(user.Id.String(), ip, false)
		return nil, errors.New("Invalid credentials")
	}
	s.repo.LogLogin(user.Id.String(), ip, true)
	return user, err
}

func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (s *UserService) UpdatePassword(tokenStr string, req UpdatePasswordRequest) error {
	claims, _ := utils.GetClaims(tokenStr, s.cfx.JWTSecret)
	userIDAny := claims["user_id"]
	userID, _ := userIDAny.(string)
	fmt.Println("User id")
	fmt.Println(userID)
	fmt.Println(req.OldPassword)
	user, err := s.repo.GetUserByID(userID)
	fmt.Println(user.Password)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !CheckPassword(user.Password, req.OldPassword) {
		return errors.New("Invalid credentials")
	}

	if !isValidPassword(req.NewPassword) {
		return errors.New("weak password")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	err = s.repo.UpdatePassword(userID, string(hash))
	if err != nil {
		return err
	}
	err = utils.SendEmail(&s.cfx, user.Email, "You have updated your Password", "You have updated your Password")
	if err != nil {
		return err
	}
	return nil
}
