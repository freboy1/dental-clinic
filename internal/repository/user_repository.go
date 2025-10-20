package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kamilakamilkami/dental-clinic/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, phone, first_name, last_name, role, push_notifications, activated, activation_token)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, registered_at;
	`
	return r.db.QueryRow(ctx, query,
		user.Email, user.PasswordHash, user.Phone, user.FirstName, user.LastName,
		user.Role, user.PushConsent, user.Activated, user.ActivationToken,
	).Scan(&user.ID, &user.RegisteredAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, phone, first_name, last_name, role, push_notifications, activated, registered_at, activation_token
		FROM users WHERE email=$1`, email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Phone, &user.FirstName, &user.LastName,
		&user.Role, &user.PushConsent, &user.Activated, &user.RegisteredAt, &user.ActivationToken)

	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
