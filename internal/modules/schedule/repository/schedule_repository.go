package repository

import (
	"context"
	
	// "dental_clinic/internal"


	"dental_clinic/internal/modules/schedule/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"fmt"
	"time"
	"github.com/google/uuid"
)

type ScheduleRepository interface {
	Create(schedule *models.Schedule) (*models.Schedule, error)
	GetSchedules() ([]models.Schedule, error)
	CreateAvailableSlot(doctor_id, clinic_address_id uuid.UUID, slot_start, slot_end time.Time) (error)
	GetAvailableSlotsByDateAndDoctorAndClinic(doctor_id, clinic_address_id uuid.UUID, date time.Time) ([]models.Slot, error)
	GetScheduleByDoctor(doctor_id uuid.UUID) ([]models.Schedule, error)
	GetSlotById(slotId uuid.UUID) (*models.Slot, error)
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


func (r *scheduleRepo) GetSchedules() ([]models.Schedule, error) {
	query := `SELECT id, doctor_id, clinic_address_id, day_of_week, start_time, end_time FROM doctor_working_hours`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.Id, &schedule.Doctor_id, &schedule.Clinic_address_id, &schedule.Day_of_week, &schedule.Start_time, &schedule.End_time); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}


func (r *scheduleRepo) CreateAvailableSlot(doctor_id, clinic_address_id uuid.UUID, slot_start, slot_end time.Time) (error) {
	slot_id := uuid.New()
	query := `INSERT INTO doctor_time_slots (id, doctor_id, clinic_address_id, slot_start, slot_end, status, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7) 
            `
	_, err := r.db.Exec(
		context.Background(),
		query,
		slot_id,
		doctor_id,
		clinic_address_id,
		slot_start,
		slot_end,
		"available",
		time.Now(),
	)
	
	if err != nil {
		return fmt.Errorf("failed to create slot: %w", err)
	}
	
	return nil
}




func (r *scheduleRepo) GetAvailableSlotsByDateAndDoctorAndClinic(doctor_id, clinic_address_id uuid.UUID, date time.Time) ([]models.Slot, error) {
	query := `SELECT id, slot_start, slot_end, status FROM doctor_time_slots WHERE doctor_id = $1 AND DATE(slot_start) = $2 AND clinic_address_id = $3 ORDER BY slot_start;`

	rows, err := r.db.Query(context.Background(), query, doctor_id, date, clinic_address_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var slot models.Slot
		if err := rows.Scan(&slot.Id, &slot.Slot_start, &slot.Slot_end, &slot.Status); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}


func (r *scheduleRepo) GetScheduleByDoctor(doctor_id uuid.UUID) ([]models.Schedule, error) {
	query := `SELECT id, doctor_id, clinic_address_id, day_of_week, start_time, end_time FROM doctor_working_hours WHERE doctor_id = $1`

	rows, err := r.db.Query(context.Background(), query, doctor_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.Id, &schedule.Doctor_id, &schedule.Clinic_address_id, &schedule.Day_of_week, &schedule.Start_time, &schedule.End_time); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}


func (r *scheduleRepo) GetSlotById(slotId uuid.UUID) (*models.Slot, error) {
	query := `SELECT id, slot_start, slot_end, status FROM doctor_time_slots WHERE id = $1 ;`

	var slot models.Slot
	err := r.db.QueryRow(context.Background(), query, slotId).Scan(&slot.Id, &slot.Slot_start, &slot.Slot_end, &slot.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &slot, nil
}