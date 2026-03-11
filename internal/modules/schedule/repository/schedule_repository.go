package repository

import (
	"context"
	
	// "dental_clinic/internal"


	"dental_clinic/internal/modules/schedule/models"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository interface {
	Create(schedule *models.Schedule) (*models.Schedule, error)
}

type scheduleRepo struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) ScheduleRepository {
	return &scheduleRepo{db: db}
}

func (r *scheduleRepo) Create(schedule *models.Schedule) (*models.Schedule, error) {
	query := `INSERT INTO doctor_working_hours (id, doctor_id, clinic_address_id, day_of_week, start_time, end_time)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(context.Background(), query, schedule.Id, schedule.Doctor_id, schedule.Clinic_address_id, schedule.Day_of_week, schedule.Start_time, schedule.End_time).
		Scan(&schedule.Id)
	return schedule, err
}