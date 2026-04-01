package repository

import (
	"context"
	
	// "dental_clinic/internal"


	"dental_clinic/internal/modules/appointment/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentRepository interface {
	Create(appointment *models.Appointment) (*models.Appointment, error)
	GetAll() ([]models.Appointment, error)
	GetByID(id string) (*models.Appointment, error)
	Update(appointment *models.Appointment) (*models.Appointment, error)
	Delete(id string) error
}

type appointmentRepo struct {
	db *pgxpool.Pool
}

func NewAppointmentRepository(db *pgxpool.Pool) AppointmentRepository {
	return &appointmentRepo{db: db}
}

func (r *appointmentRepo) Create(appointment *models.Appointment) (*models.Appointment, error) {
	query := `INSERT INTO appointments (id, doctor_id, clinic_address_id, service_id, user_id, start_time, end_time, status, created_at, name, email)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	err := r.db.QueryRow(context.Background(), query, appointment.Id, appointment.Doctor_id, appointment.Clinic_address_id, appointment.Service_id, appointment.User_id, appointment.Start_time, appointment.End_time, appointment.Status, appointment.Created_at, appointment.Name, appointment.Email).
		Scan(&appointment.Id)
	return appointment, err
}



func (r *appointmentRepo) GetAll() ([]models.Appointment, error) {
	query := `SELECT id, doctor_id, clinic_address_id, service_id, user_id, start_time, end_time, status, name, email FROM appointments`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(&appointment.Id, &appointment.Doctor_id, &appointment.Clinic_address_id, &appointment.Service_id, &appointment.User_id, &appointment.Start_time, &appointment.End_time, &appointment.Status, &appointment.Name, &appointment.Email); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepo) GetByID(id string) (*models.Appointment, error) {
	query := `SELECT id, doctor_id, clinic_address_id, service_id, user_id, start_time, end_time, status, name, email FROM appointments WHERE id = $1`

	var appointment models.Appointment
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&appointment.Id, &appointment.Doctor_id, &appointment.Clinic_address_id, &appointment.Service_id,
		&appointment.User_id, &appointment.Start_time, &appointment.End_time, &appointment.Status,
		&appointment.Name, &appointment.Email,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentRepo) Update(appointment *models.Appointment) (*models.Appointment, error) {
	query := `
		UPDATE appointments
		SET doctor_id=$1, clinic_address_id=$2, service_id=$3, start_time=$4, end_time=$5, status=$6, name=$7, email=$8
		WHERE id=$9
		RETURNING id, doctor_id, clinic_address_id, service_id, user_id, start_time, end_time, status, name, email
	`
	err := r.db.QueryRow(context.Background(), query,
		appointment.Doctor_id, appointment.Clinic_address_id, appointment.Service_id,
		appointment.Start_time, appointment.End_time, appointment.Status,
		appointment.Name, appointment.Email, appointment.Id,
	).Scan(
		&appointment.Id, &appointment.Doctor_id, &appointment.Clinic_address_id, &appointment.Service_id,
		&appointment.User_id, &appointment.Start_time, &appointment.End_time, &appointment.Status,
		&appointment.Name, &appointment.Email,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return appointment, nil
}

func (r *appointmentRepo) Delete(id string) error {
	query := `DELETE FROM appointments WHERE id=$1`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}