package services

import (
	"dental_clinic/internal/config"

	"dental_clinic/internal/modules/schedule/models"
	"dental_clinic/internal/modules/schedule/dto"
	"dental_clinic/internal/modules/schedule/repository"

	"errors"
	"fmt"
	
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