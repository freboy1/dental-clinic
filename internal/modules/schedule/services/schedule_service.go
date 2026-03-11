package services

import (
	"dental_clinic/internal/config"

	"dental_clinic/internal/modules/schedule/models"
	"dental_clinic/internal/modules/schedule/dto"
	"dental_clinic/internal/modules/schedule/repository"

	"errors"
	"fmt"
	"time"
	
	"github.com/google/uuid"
)

type ScheduleService struct {
	repo repository.ScheduleRepository
	cfx  config.Config
}

func NewScheduleService(r repository.ScheduleRepository, cfx config.Config) *ScheduleService {
	return &ScheduleService{
		repo: r,
		cfx:  cfx,
	}
}




func (s *ScheduleService) CreateSchedule(doctor_id uuid.UUID, req dto.CreateScheduleRequest) (*models.Schedule, error) {

	if req.Clinic_address_id == "" {
		return nil, fmt.Errorf("schedule clinic_address_id is required")
	}

	clinic_address_id, err := uuid.Parse(req.Clinic_address_id)
	if err != nil {
		return nil, errors.New("invalid clinic_address_id")
	}

	schedule := &models.Schedule{
		Id: uuid.New(),
		Doctor_id: doctor_id,
		Clinic_address_id: clinic_address_id,
		Day_of_week: req.Day_of_week,
		Start_time: req.Start_time,
		End_time: req.End_time,
	}

	return s.repo.Create(schedule)
}



func (s *ScheduleService) GetSchedules() ([]models.Schedule, error) {
	return s.repo.GetSchedules()
}


func (s *ScheduleService) GenerateSlots(req dto.GenerateSlotsRequest) error {

	schedules, err := s.GetSchedules()
	if (err != nil) {
		return err
	}

	fromDate, err := time.Parse("2006-01-02", req.From_date)
	if err != nil {
		return err
	}

	toDate, err := time.Parse("2006-01-02", req.To_date)
	if err != nil {
		return err
	}

	for _, schedule := range schedules {

		startTime, err := time.Parse("15:04:05", schedule.Start_time)
		if err != nil {
			return err
		}

		endTime, err := time.Parse("15:04:05", schedule.End_time)
		if err != nil {
			return err
		}
		

		for date := fromDate; !date.After(toDate); date = date.AddDate(0, 0, 1) {

			if int(date.Weekday()) != schedule.Day_of_week {
				continue
			}

			start := time.Date(
				date.Year(), date.Month(), date.Day(),
				startTime.Hour(), startTime.Minute(),
				0, 0, time.UTC,
			)

			end := time.Date(
				date.Year(), date.Month(), date.Day(),
				endTime.Hour(), endTime.Minute(),
				0, 0, time.UTC,
			)

			for t := start; t.Before(end); t = t.Add(30 * time.Minute) {

				slotEnd := t.Add(30 * time.Minute)

				err := s.repo.CreateAvailableSlot(
					schedule.Doctor_id,
					schedule.Clinic_address_id,
					t,
					slotEnd,
				)
				if err != nil {
					return err
				}
			}
		}

	}
	
	// return s.repo.Create(schedule)
	return nil
}