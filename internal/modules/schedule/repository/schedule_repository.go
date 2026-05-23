package repository

import (
	"context"

	// "dental_clinic/internal"

	"fmt"
	"time"

	"dental_clinic/internal/modules/schedule/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository interface {
	Create(schedule *models.Schedule) (*models.Schedule, error)
	GetSchedules() ([]models.Schedule, error)
	CreateAvailableSlot(doctor_id, clinic_address_id uuid.UUID, slot_start, slot_end time.Time) error
	GetAvailableSlotsByDateAndDoctorAndClinic(doctor_id, clinic_address_id uuid.UUID, date time.Time) ([]models.Slot, error)
	GetScheduleByDoctor(doctor_id uuid.UUID) ([]models.Schedule, error)
	GetSlotById(slotId uuid.UUID) (*models.Slot, error)
	UpdateSlotStatus(slotId uuid.UUID, status string) error
	UpdateSlotStatusTx(slotId uuid.UUID, status string, tx pgx.Tx) error
	GetScheduleById(schedule_id uuid.UUID) (*models.Schedule, error)
	DeleteScheduleById(schedule_id uuid.UUID) error
	UpdateScheduleById(id string, doctor *models.Schedule) error
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

func (r *scheduleRepo) CreateAvailableSlot(doctor_id, clinic_address_id uuid.UUID, slot_start, slot_end time.Time) error {
	slot_id := uuid.New()
	query := `INSERT INTO doctor_time_slots (id, doctor_id, clinic_address_id, slot_start, slot_end, status, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7) 
			ON CONFLICT DO NOTHING
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

func (r *scheduleRepo) UpdateSlotStatus(slotId uuid.UUID, status string) error {
	query := `UPDATE doctor_time_slots SET status = $1 WHERE id = $2 ;`

	result, err := r.db.Exec(context.Background(), query, status, slotId)
	if err != nil {
		return fmt.Errorf("failed to update slot: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("slot not found")
	}

	return nil
}

func (r *scheduleRepo) UpdateSlotStatusTx(slotId uuid.UUID, status string, tx pgx.Tx) error {
	query := `UPDATE doctor_time_slots SET status = $1 WHERE id = $2 ;`

	result, err := tx.Exec(context.Background(), query, status, slotId)
	if err != nil {
		return fmt.Errorf("failed to update slot: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("slot not found")
	}

	return nil
}

func (r *scheduleRepo) GetScheduleById(schedule_id uuid.UUID) (*models.Schedule, error) {
	query := `SELECT id, doctor_id, clinic_address_id, day_of_week, start_time, end_time FROM doctor_working_hours WHERE id = $1 ;`

	var schedule models.Schedule
	err := r.db.QueryRow(context.Background(), query, schedule_id).Scan(&schedule.Id, &schedule.Doctor_id, &schedule.Clinic_address_id, &schedule.Day_of_week, &schedule.Start_time, &schedule.End_time)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepo) DeleteScheduleById(schedule_id uuid.UUID) error {
	query := `DELETE FROM doctor_working_hours WHERE id=$1`
	result, err := r.db.Exec(context.Background(), query, schedule_id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *scheduleRepo) UpdateScheduleById(id string, doctor *models.Schedule) error {
	query := `
		UPDATE doctor_working_hours
		SET doctor_id=$1, clinic_address_id=$2, day_of_week=$3, start_time=$4, end_time=$5
		WHERE id=$6
		RETURNING id, doctor_id, clinic_address_id, day_of_week, start_time, end_time
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		doctor.Doctor_id,
		doctor.Clinic_address_id,
		doctor.Day_of_week,
		doctor.Start_time,
		doctor.End_time,
		id,
	).Scan(&doctor.Id, &doctor.Doctor_id, &doctor.Clinic_address_id, &doctor.Day_of_week, &doctor.Start_time, &doctor.End_time)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
