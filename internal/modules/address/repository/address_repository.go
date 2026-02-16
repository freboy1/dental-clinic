package repository

import (
	"context"
	
	// "dental_clinic/internal"


	"dental_clinic/internal/modules/address/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressRepository interface {
	Create(address *models.Address) (*models.Address, error)
	GetByID(id string) (*models.Address, error)
	// Update(id string, address *models.Address) (*models.Address, error)
	Delete(id string) error
	GetAll() ([]models.Address, error)
}

type addressRepo struct {
	db *pgxpool.Pool
}

func NewAddressRepository(db *pgxpool.Pool) AddressRepository {
	return &addressRepo{db: db}
}

func (r *addressRepo) Create(address *models.Address) (*models.Address, error) {
	query := `INSERT INTO addresses (id, country, city, street, building, latitude, longitude)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(context.Background(), query, address.ID, address.Country, address.City, address.Street, address.Building, address.Latitude, address.Longitude).
		Scan(&address.ID)
	return address, err
}

func (r *addressRepo) GetAll() ([]models.Address, error) {
	query := `SELECT id, country, city, street, building, latitude, longitude FROM addresses`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.Address
	for rows.Next() {
		var address models.Address
		if err := rows.Scan(&address.ID, &address.Country, &address.City, &address.Street, &address.Building, &address.Latitude, &address.Longitude); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *addressRepo) GetByID(id string) (*models.Address, error) {
	query := `SELECT country, city, street, building, latitude, longitude FROM addresses WHERE id = $1`
	var address models.Address
	err := r.db.QueryRow(context.Background(), query, id).Scan(&address.Country, &address.City, &address.Street, &address.Building, &address.Latitude, &address.Longitude)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &address, nil
}

func (r *addressRepo) Delete(id string) error {
	query := `DELETE FROM addresses WHERE id=$1`
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
