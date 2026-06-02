package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	// "dental_clinic/internal"

	"dental_clinic/internal/modules/clinic/models"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClinicRepository interface {
	Create(clinic *models.Clinic) (*models.Clinic, error)
	GetAll() ([]*models.Clinic, error)
	GetByID(id uuid.UUID) (*models.Clinic, error)
	Update(clinic *models.Clinic) (*models.Clinic, error)
	Delete(id uuid.UUID) error
	UpdateLogo(id uuid.UUID, logoURL string) error
	DeleteLogo(id uuid.UUID) error
	AddAddress(id, clinic_id uuid.UUID, address_id string, is_main bool) error
	GetClinicAddress(id uuid.UUID) ([]models.ClinicAddress, error)
	GetClinicAddressByID(id uuid.UUID) (*models.ClinicAddress, error)
	GetClinicByAddressId(id uuid.UUID) (string, error)
	DeleteAddress(id, address_id uuid.UUID) error
	UpdateAddressCover(id uuid.UUID, coverURL string) error
	DeleteAddressCover(id uuid.UUID) error
	AddGalleryImage(image models.ClinicAddressGalleryImage) error
	GetGalleryImages(clinicAddressID uuid.UUID) ([]models.ClinicAddressGalleryImage, error)
	GetGalleryImage(id uuid.UUID) (*models.ClinicAddressGalleryImage, error)
	UpdateGalleryImage(id uuid.UUID, imageURL string) error
	DeleteGalleryImage(id uuid.UUID) error
}

type clinicRepo struct {
	db *pgxpool.Pool
}

func NewClinicRepository(db *pgxpool.Pool) ClinicRepository {
	return &clinicRepo{db: db}
}

func (r *clinicRepo) Create(clinic *models.Clinic) (*models.Clinic, error) {
	query := `INSERT INTO clinics (id, name, description, phone, email, website, is_active, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
            RETURNING id, created_at`

	err := r.db.QueryRow(
		context.Background(),
		query,
		clinic.Id,
		clinic.Name,
		clinic.Description,
		clinic.Phone,
		clinic.Email,
		clinic.Website,
		clinic.IsActive,
		clinic.CreatedAt,
	).Scan(&clinic.Id, &clinic.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create clinic: %w", err)
	}

	return clinic, nil
}

func (r *clinicRepo) GetAll() ([]*models.Clinic, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.description,
			c.phone,
			c.email,
			c.website,
			c.is_active,
			c.created_at,
			COALESCE(c.logo_url, ''),
			COALESCE(ROUND(AVG(cr.rating)::numeric, 2), 0)::float8 AS rating
		FROM clinics c
		LEFT JOIN clinic_reviews cr ON cr.clinic_id = c.id
		WHERE c.is_active = true
		GROUP BY c.id, c.name, c.description, c.phone, c.email, c.website, c.is_active, c.created_at, c.logo_url
		ORDER BY rating DESC, c.created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinics: %w", err)
	}
	defer rows.Close()

	var clinics []*models.Clinic
	for rows.Next() {
		clinic := &models.Clinic{}
		err := rows.Scan(
			&clinic.Id,
			&clinic.Name,
			&clinic.Description,
			&clinic.Phone,
			&clinic.Email,
			&clinic.Website,
			&clinic.IsActive,
			&clinic.CreatedAt,
			&clinic.LogoURL,
			&clinic.Rating,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan clinic: %w", err)
		}
		clinics = append(clinics, clinic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clinics: %w", err)
	}

	return clinics, nil
}

func (r *clinicRepo) GetByID(id uuid.UUID) (*models.Clinic, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.description,
			c.phone,
			c.email,
			c.website,
			c.is_active,
			c.created_at,
			COALESCE(c.logo_url, ''),
			COALESCE(ROUND(AVG(cr.rating)::numeric, 2), 0)::float8 AS rating
		FROM clinics c
		LEFT JOIN clinic_reviews cr ON cr.clinic_id = c.id
		WHERE c.id = $1
		GROUP BY c.id, c.name, c.description, c.phone, c.email, c.website, c.is_active, c.created_at, c.logo_url
	`

	clinic := &models.Clinic{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&clinic.Id,
		&clinic.Name,
		&clinic.Description,
		&clinic.Phone,
		&clinic.Email,
		&clinic.Website,
		&clinic.IsActive,
		&clinic.CreatedAt,
		&clinic.LogoURL,
		&clinic.Rating,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get clinic: %w", err)
	}

	return clinic, nil
}

func (r *clinicRepo) UpdateLogo(id uuid.UUID, logoURL string) error {
	result, err := r.db.Exec(context.Background(), `UPDATE clinics SET logo_url = $2 WHERE id = $1 AND is_active = true`, id, logoURL)
	if err != nil {
		return fmt.Errorf("failed to update clinic logo: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("clinic not found")
	}
	return nil
}

func (r *clinicRepo) DeleteLogo(id uuid.UUID) error {
	result, err := r.db.Exec(context.Background(), `UPDATE clinics SET logo_url = NULL WHERE id = $1 AND is_active = true`, id)
	if err != nil {
		return fmt.Errorf("failed to delete clinic logo: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("clinic not found")
	}
	return nil
}

func (r *clinicRepo) Update(clinic *models.Clinic) (*models.Clinic, error) {
	query := `UPDATE clinics 
            SET name = $2, 
                description = $3, 
                phone = $4, 
                email = $5, 
                website = $6, 
                is_active = $7
            WHERE id = $1
            RETURNING id, created_at`

	err := r.db.QueryRow(
		context.Background(),
		query,
		clinic.Id,
		clinic.Name,
		clinic.Description,
		clinic.Phone,
		clinic.Email,
		clinic.Website,
		clinic.IsActive,
	).Scan(&clinic.Id, &clinic.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update clinic: %w", err)
	}

	return clinic, nil
}

func (r *clinicRepo) Delete(id uuid.UUID) error {
	query := `UPDATE clinics SET is_active = false WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete clinic: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("clinic not found")
	}

	return nil
}

func (r *clinicRepo) AddAddress(id, clinic_id uuid.UUID, address_id string, is_main bool) error {
	query := `INSERT INTO clinic_addresses (id, clinic_id, address_id, is_main)
            VALUES ($1, $2, $3, $4) 
            `

	if is_main {
		err := r.UnsetMainAddress(clinic_id)
		if err != nil {
			return err
		}
	}

	_, err := r.db.Exec(
		context.Background(),
		query,
		id,
		clinic_id,
		address_id,
		is_main,
	)
	if err != nil {
		return fmt.Errorf("failed to create clinic address: %w", err)
	}

	return nil
}

func (r *clinicRepo) UnsetMainAddress(clinic_id uuid.UUID) error {
	query := `UPDATE clinic_addresses SET is_main = false WHERE clinic_id = $1`

	_, err := r.db.Exec(context.Background(), query, clinic_id)
	if err != nil {
		return fmt.Errorf("failed to update clinic address: %w", err)
	}

	// if result.RowsAffected() == 0 {
	// 	return fmt.Errorf("clinic address not found")
	// }

	return nil
}

func (r *clinicRepo) GetClinicAddress(id uuid.UUID) ([]models.ClinicAddress, error) {
	query := `SELECT id, clinic_id, address_id, is_main, COALESCE(cover_image_url, '')
            FROM clinic_addresses
            WHERE clinic_id = $1
            `

	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinic address: %w", err)
	}
	defer rows.Close()

	var clinics []models.ClinicAddress
	for rows.Next() {
		clinic := models.ClinicAddress{}
		err := rows.Scan(
			&clinic.Id,
			&clinic.ClinicId,
			&clinic.AddressId,
			&clinic.IsMain,
			&clinic.CoverImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan clinic: %w", err)
		}
		clinics = append(clinics, clinic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clinics: %w", err)
	}

	return clinics, nil
}

func (r *clinicRepo) GetClinicAddressByID(id uuid.UUID) (*models.ClinicAddress, error) {
	query := `SELECT id, clinic_id, address_id, is_main, COALESCE(cover_image_url, '')
            FROM clinic_addresses
            WHERE id = $1
            `

	clinic := &models.ClinicAddress{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&clinic.Id,
		&clinic.ClinicId,
		&clinic.AddressId,
		&clinic.IsMain,
		&clinic.CoverImageURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinic address: %w", err)
	}
	return clinic, nil
}

func (r *clinicRepo) UpdateAddressCover(id uuid.UUID, coverURL string) error {
	result, err := r.db.Exec(context.Background(), `UPDATE clinic_addresses SET cover_image_url = $2 WHERE id = $1`, id, coverURL)
	if err != nil {
		return fmt.Errorf("failed to update clinic address cover: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("clinic address not found")
	}
	return nil
}

func (r *clinicRepo) DeleteAddressCover(id uuid.UUID) error {
	result, err := r.db.Exec(context.Background(), `UPDATE clinic_addresses SET cover_image_url = NULL WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete clinic address cover: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("clinic address not found")
	}
	return nil
}

func (r *clinicRepo) AddGalleryImage(image models.ClinicAddressGalleryImage) error {
	query := `
		INSERT INTO clinic_address_gallery (id, clinic_address_id, image_url)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(context.Background(), query, image.Id, image.ClinicAddressId, image.ImageURL)
	if err != nil {
		return fmt.Errorf("failed to add gallery image: %w", err)
	}
	return nil
}

func (r *clinicRepo) GetGalleryImages(clinicAddressID uuid.UUID) ([]models.ClinicAddressGalleryImage, error) {
	query := `
		SELECT id, clinic_address_id, image_url
		FROM clinic_address_gallery
		WHERE clinic_address_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(context.Background(), query, clinicAddressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gallery images: %w", err)
	}
	defer rows.Close()

	var images []models.ClinicAddressGalleryImage
	for rows.Next() {
		var image models.ClinicAddressGalleryImage
		err := rows.Scan(&image.Id, &image.ClinicAddressId, &image.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gallery image: %w", err)
		}
		images = append(images, image)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating gallery images: %w", err)
	}
	return images, nil
}

func (r *clinicRepo) GetGalleryImage(id uuid.UUID) (*models.ClinicAddressGalleryImage, error) {
	query := `
		SELECT id, clinic_address_id, image_url
		FROM clinic_address_gallery
		WHERE id = $1
	`
	var image models.ClinicAddressGalleryImage
	err := r.db.QueryRow(context.Background(), query, id).Scan(&image.Id, &image.ClinicAddressId, &image.ImageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get gallery image: %w", err)
	}
	return &image, nil
}

func (r *clinicRepo) UpdateGalleryImage(id uuid.UUID, imageURL string) error {
	result, err := r.db.Exec(context.Background(), `UPDATE clinic_address_gallery SET image_url = $2 WHERE id = $1`, id, imageURL)
	if err != nil {
		return fmt.Errorf("failed to update gallery image: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("gallery image not found")
	}
	return nil
}

func (r *clinicRepo) DeleteGalleryImage(id uuid.UUID) error {
	result, err := r.db.Exec(context.Background(), `DELETE FROM clinic_address_gallery WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete gallery image: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("gallery image not found")
	}
	return nil
}

func (r *clinicRepo) DeleteAddress(id, address_id uuid.UUID) error {
	query := `DELETE FROM clinic_addresses WHERE address_id = $1 AND clinic_id = $2`

	result, err := r.db.Exec(context.Background(), query, address_id, id)
	if err != nil {
		return fmt.Errorf("failed to delete clinic address: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("clinic address not found")
	}

	return nil
}

func (r *clinicRepo) GetClinicByAddressId(id uuid.UUID) (string, error) {
	var clinic_id string
	query := `SELECT clinic_id
            FROM clinic_addresses
            WHERE id = $1
            `
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&clinic_id,
	)

	if err != nil {
		return "", fmt.Errorf("failed to get clinic: %w", err)
	}

	return clinic_id, nil
}
