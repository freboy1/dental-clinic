package repository

import (
	"context"
	"dental_clinic/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Update(id string, user *models.User) (*models.User, error)
	Delete(id string) error
	GetAll() ([]models.User, error)
	FindUserIdByToken(token string) (string, error)
	MarkUserAsVerified(user_id string) error
	SaveVerificationToken(user_id, token string) error
	GetUserByEmail(email string) (*models.User, error)
	UpdatePassword(user_id, new_password string) error
	GetUserByID(id string) (*models.User, error)
	LogLogin(userID, ip string, success bool) error
	SaveEmailVerificationToken(userID, NewEmail, verifyToken string) error
	VerifyEmailToken(token string) (string, string, error)
	UpdateEmailInDatabase(userId, newEmail string) error
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (role, email, password, name, gender, age, push_consent)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(context.Background(), query, user.Role, user.Email, user.Password, user.Name, user.Gender, user.Age, user.Push_consent).
		Scan(&user.Id)
	return user, err
}

func (r *userRepo) GetAll() ([]models.User, error) {
	query := `SELECT id, role, email, name, gender, age, push_consent FROM users`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Role, &u.Email, &u.Name, &u.Gender, &u.Age, &u.Push_consent); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) GetByID(id string) (*models.User, error) {
	query := `SELECT role, email, name, gender, age, push_consent FROM users WHERE id = $1`
	var u models.User
	err := r.db.QueryRow(context.Background(), query, id).Scan(&u.Role, &u.Email, &u.Name, &u.Gender, &u.Age, &u.Push_consent)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(id string, user *models.User) (*models.User, error) {
	query := `
		UPDATE users
		SET name=$1, email=$2, role=$3, gender=$4, age=$5, push_consent=$6
		WHERE id=$7
		RETURNING id, role, email, name, gender, age, push_consent
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		user.Name,
		user.Email,
		user.Role,
		user.Gender,
		user.Age,
		user.Push_consent,
		id,
	).Scan(
		&user.Id,
		&user.Role,
		&user.Email,
		&user.Name,
		&user.Gender,
		&user.Age,
		&user.Push_consent,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Delete(id string) error {
	query := `DELETE FROM users WHERE id=$1`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *userRepo) FindUserIdByToken(token string) (string, error) {
	query := "SELECT user_id FROM verification_tokens WHERE token = $1 AND expires_at > NOW()"
	var user_id string
	err := r.db.QueryRow(context.Background(), query, token).Scan(&user_id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return user_id, nil
}

func (r *userRepo) MarkUserAsVerified(user_id string) error {
	query := "UPDATE users SET is_verified=TRUE WHERE id=$1"
	_, err := r.db.Exec(context.Background(), query, user_id)

	return err
}

func (r *userRepo) SaveVerificationToken(user_id, token string) error {
	query := `
		INSERT INTO verification_tokens (user_id, token)
		VALUES ($1, $2)
	`
	_, err := r.db.Exec(context.Background(), query, user_id, token)
	return err
}

func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, password, role, email, name, gender, age, push_consent FROM users WHERE email = $1 AND is_verified = true`
	var u models.User
	err := r.db.QueryRow(context.Background(), query, email).Scan(&u.Id, &u.Password, &u.Role, &u.Email, &u.Name, &u.Gender, &u.Age, &u.Push_consent)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) UpdatePassword(user_id, new_password string) error {
	query := "UPDATE users SET password=$1 WHERE id=$2"
	_, err := r.db.Exec(context.Background(), query, new_password, user_id)
	return err
}

func (r *userRepo) GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, password, role, email, name, gender, age, push_consent FROM users WHERE id = $1`
	var u models.User
	err := r.db.QueryRow(context.Background(), query, id).Scan(&u.Id, &u.Password, &u.Role, &u.Email, &u.Name, &u.Gender, &u.Age, &u.Push_consent)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) LogLogin(userID, ip string, success bool) error {
	if userID == "" {
		query := `INSERT INTO login_logs (user_id, ip_address, success)
                  VALUES (NULL, $1, $2)`
		_, err := r.db.Exec(context.Background(), query, ip, success)
		return err
	}
	query := `INSERT INTO login_logs (user_id, ip_address, success) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(context.Background(), query, userID, ip, success)
	return err
}
func (r *userRepo) SaveEmailVerificationToken(userID, NewEmail, verifyToken string) error {
	query := `
		INSERT INTO email_change_tokens (user_id, new_email, token)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(context.Background(), query, userID, NewEmail, verifyToken)
	return err
}

func (r *userRepo) VerifyEmailToken(token string) (string, string, error) {
	var userID, newEmail string
	query := `SELECT user_id, new_email FROM email_change_tokens WHERE token = $1`
	err := r.db.QueryRow(context.Background(), query, token).Scan(&userID, &newEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", nil
		}
		return "", "", err
	}
	return userID, newEmail, nil
}

func (r *userRepo) UpdateEmailInDatabase(userId, newEmail string) error {
	query := "UPDATE users SET email=$1 WHERE id=$2"
	_, err := r.db.Exec(context.Background(), query, newEmail, userId)
	return err
}
