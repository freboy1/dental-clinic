package repository

import (
	"context"
	
	// "dental_clinic/internal"


	"dental_clinic/internal/modules/appointment/models"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentRepository interface {
	Create(appointment *models.Appointment) (*models.Appointment, error)
}

type appointmentRepo struct {
	db *pgxpool.Pool
}

func NewAppointmentRepository(db *pgxpool.Pool) AppointmentRepository {
	return &appointmentRepo{db: db}
}

func (r *appointmentRepo) Create(appointment *models.Appointment) (*models.Appointment, error) {
	query := `INSERT INTO appointmentes (id, doctor_id, clinic_address_id, service_id, user_id, start_time, end_time, status, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	err := r.db.QueryRow(context.Background(), query, appointment.Id, appointment.Doctor_id, appointment.Clinic_address_id, appointment.Service_id, appointment.User_id, appointment.Start_time, appointment.End_time, appointment.Status, appointment.Created_at).
		Scan(&appointment.Id)
	return appointment, err
}
