package repository

import (
	"database/sql"
	"dental_clinic/internal/models"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Update(id string, user *models.User) (*models.User, error)
	Delete(id string) error
	GetAll() ([]models.User, error)
	FindUserIdByToken(token string) (string, error)
	MarkUserAsVerified(user_id string) (error)
	SaveVerificationToken(user_id, token string) (error)
	GetUserByEmail(email string) (*models.User, error)
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

func (r *userRepo) GetAll() ([]models.User, error) {
	query := `SELECT id, role, email, name, gender, age, push_consent FROM users`

	rows, err := r.db.Query(query)
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
	query := `SELECT id, role, email, name, gender, age, push_consent FROM users WHERE id = $1`
	var u models.User
	err := r.db.QueryRow(query, id).Scan(&u.Id, &u.Role, &u.Email, &u.Name, &u.Gender, &u.Age, &u.Push_consent)
	if err != nil {
		if err == sql.ErrNoRows {
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
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Delete(id string) error {
	query := `DELETE FROM users WHERE id=$1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepo) FindUserIdByToken(token string) (string, error) {
	query := "SELECT user_id FROM verification_tokens WHERE token = $1 AND expires_at > NOW()"
	var user_id string
	err := r.db.QueryRow(query, token).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return user_id, nil
}

func (r *userRepo) MarkUserAsVerified(user_id string) (error) {
	query := "UPDATE users SET is_verified=TRUE WHERE id=$1"
	_, err := r.db.Exec(query, user_id)

	return err
}

func (r *userRepo) SaveVerificationToken(user_id, token string) (error) {
	query := `
		INSERT INTO verification_tokens (user_id, token)
		VALUES ($1, $2)
	`
	_, err := r.db.Exec(query, user_id, token)
	return err
}

func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, password, role, email, name, gender, age, push_consent FROM users WHERE email = $1`
	var u models.User
	err := r.db.QueryRow(query, email).Scan(&u.Id, &u.Password, &u.Role, &u.Email, &u.Name, &u.Gender, &u.Age, &u.Push_consent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}