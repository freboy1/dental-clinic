package services

import (
	"errors"
	"time"

	"dental_clinic/internal/modules/reviews/models"
	"dental_clinic/internal/modules/reviews/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ReviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(r repository.ReviewRepository) *ReviewService {
	return &ReviewService{repo: r}
}

func (s *ReviewService) CreateAppointmentReviewTx(appointmentId, doctorId, clinicId, userId uuid.UUID, doctorRating, clinicRating int, clinicComment string, tx pgx.Tx) error {
	if doctorRating < 1 || doctorRating > 5 {
		return errors.New("doctor_rating must be between 1 and 5")
	}
	if clinicRating < 1 || clinicRating > 5 {
		return errors.New("clinic_rating must be between 1 and 5")
	}

	now := time.Now()
	rating := &models.DoctorRating{
		Id:            uuid.New(),
		AppointmentId: appointmentId,
		DoctorId:      doctorId,
		UserId:        userId,
		Rating:        doctorRating,
		CreatedAt:     now,
	}
	review := &models.ClinicReview{
		Id:            uuid.New(),
		AppointmentId: appointmentId,
		ClinicId:      clinicId,
		UserId:        userId,
		Rating:        clinicRating,
		Comment:       clinicComment,
		CreatedAt:     now,
	}

	if err := s.repo.CreateDoctorRatingTx(rating, tx); err != nil {
		return err
	}
	return s.repo.CreateClinicReviewTx(review, tx)
}
