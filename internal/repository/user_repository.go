package repository

import (
	"database/sql"
	"dental_clinic/internal/models"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Update(id string, user *models.User) error
	Delete(id string) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (role, email, password, name, gender, age, push_consent)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(query, user.Role, user.Email, user.Password, user.Name, user.Gender, user.Age, user.Push_consent).
		Scan(&user.Id)
	return user, err
}

func (r *userRepo) Delete(id string) error {
	panic("unimplemented")
}

func (r *userRepo) GetByID(id string) (*models.User, error) {
	panic("unimplemented")
}

func (r *userRepo) Update(id string, user *models.User) error {
	panic("unimplemented")
}


