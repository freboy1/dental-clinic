package repository

import (
	"context"

	"dental_clinic/internal/modules/reviews/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository interface {
	CreateDoctorRatingTx(rating *models.DoctorRating, tx pgx.Tx) error
	CreateClinicReviewTx(review *models.ClinicReview, tx pgx.Tx) error
}

type reviewRepo struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) ReviewRepository {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) CreateDoctorRatingTx(rating *models.DoctorRating, tx pgx.Tx) error {
	query := `
		INSERT INTO doctor_ratings (id, appointment_id, doctor_id, user_id, rating, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.Exec(
		context.Background(),
		query,
		rating.Id,
		rating.AppointmentId,
		rating.DoctorId,
		rating.UserId,
		rating.Rating,
		rating.CreatedAt,
	)
	return err
}

func (r *reviewRepo) CreateClinicReviewTx(review *models.ClinicReview, tx pgx.Tx) error {
	query := `
		INSERT INTO clinic_reviews (id, appointment_id, clinic_id, user_id, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(
		context.Background(),
		query,
		review.Id,
		review.AppointmentId,
		review.ClinicId,
		review.UserId,
		review.Rating,
		review.Comment,
		review.CreatedAt,
	)
	return err
}
