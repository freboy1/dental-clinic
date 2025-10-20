package repository

import (
	"context"
	"errors"

	"github.com/freboy1/dental-clinic/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, phone, first_name, last_name, role, push_consent, activated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, registered_at;
	`
	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Phone,
		user.FirstName,
		user.LastName,
		user.Role,
		user.PushConsent,
		user.Activated,
	).Scan(&user.ID, &user.RegisteredAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, phone, first_name, last_name, role, push_consent, activated, registered_at
		FROM users WHERE email = $1;
	`

	var u models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Phone, &u.FirstName,
		&u.LastName, &u.Role, &u.PushConsent, &u.Activated, &u.RegisteredAt,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, phone, first_name, last_name, role, push_consent, activated, registered_at
		FROM users WHERE id = $1;
	`

	var u models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Phone, &u.FirstName,
		&u.LastName, &u.Role, &u.PushConsent, &u.Activated, &u.RegisteredAt,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, phone = $2, first_name = $3, last_name = $4, push_consent = $5
		WHERE id = $6;
	`
	_, err := r.db.Exec(ctx, query,
		user.Email, user.Phone, user.FirstName, user.LastName, user.PushConsent, user.ID,
	)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
