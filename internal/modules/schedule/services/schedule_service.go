package services

import (
	"dental_clinic/internal/config"

	"dental_clinic/internal/modules/schedule/dto"
	"dental_clinic/internal/modules/schedule/models"
	"dental_clinic/internal/modules/schedule/repository"

	"dental_clinic/internal/modules/services/services"

	"errors"
	"fmt"
	"time"

	"math"

	"github.com/google/uuid"
)

type ScheduleService struct {
	repo       repository.ScheduleRepository
	cfx        config.Config
	serviceSrv services.ServiceService
}

func NewScheduleService(r repository.ScheduleRepository, cfx config.Config, serviceSrv services.ServiceService) *ScheduleService {
	return &ScheduleService{
		repo:       r,
		cfx:        cfx,
		serviceSrv: serviceSrv,
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
		Id:                uuid.New(),
		Doctor_id:         doctor_id,
		Clinic_address_id: clinic_address_id,
		Day_of_week:       req.Day_of_week,
		Start_time:        req.Start_time,
		End_time:          req.End_time,
	}

	return s.repo.Create(schedule)
}

func (s *ScheduleService) GetSchedules() ([]models.Schedule, error) {
	return s.repo.GetSchedules()
}

func (s *ScheduleService) GenerateSlots(req dto.GenerateSlotsRequest) error {

	schedules, err := s.GetSchedules()
	if err != nil {
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

func (s *ScheduleService) GetAvailableSlots(doctorID, serviceID, clinic_addressID uuid.UUID, date time.Time) ([]models.Slot, error) {

	service, err := s.serviceSrv.GetServiceByID(serviceID.String())

	if err != nil {
		return nil, err
	}

	raw_slots, err := s.GetAvailableSlotsByDateAndDoctorAndClinic(doctorID, clinic_addressID, date)
	if err != nil {
		return nil, err
	}

	required_slots := s.HowManySlots(service.Duration)

	slots := FindAvailableSlots(raw_slots, required_slots)

	return slots, nil
}

func (s *ScheduleService) GetAvailableSlotsByDateAndDoctorAndClinic(doctorID uuid.UUID, clinic_addressID uuid.UUID, date time.Time) ([]models.Slot, error) {
	return s.repo.GetAvailableSlotsByDateAndDoctorAndClinic(doctorID, clinic_addressID, date)
}

func (s *ScheduleService) GetSlotById(slotId uuid.UUID) (*models.Slot, error) {
	return s.repo.GetSlotById(slotId)
}

func (s *ScheduleService) HowManySlots(duration int) int {
	slotDuration := 30

	slots := int(math.Ceil(float64(duration) / float64(slotDuration)))

	if slots < 1 {
		slots = 1
	}

	return slots

}

func (s *ScheduleService) AreSlotsAvailable(slots []models.Slot, startSlotID uuid.UUID, requiredSlots int) ([]models.Slot, error) {
	var startIndex int = -1

	for i, slot := range slots {
		if slot.Id == startSlotID {
			startIndex = i
			break
		}
	}

	if startIndex == -1 {
		return nil, errors.New("start slot not found")
	}

	if startIndex+requiredSlots > len(slots) {
		return nil, errors.New("not enough slots after start slot")
	}

	for i := 0; i < requiredSlots; i++ {
		current := slots[startIndex+i]

		if current.Status != "available" {
			return nil, errors.New("slot already booked")
		}

		if i > 0 {
			prev := slots[startIndex+i-1]
			if !prev.Slot_end.Equal(current.Slot_start) {
				return nil, errors.New("slots are not consecutive")
			}
		}
	}

	return slots[startIndex : startIndex+requiredSlots], nil
}

func FindAvailableSlots(slots []models.Slot, requiredSlots int) []models.Slot {

	var result []models.Slot

	for i := 0; i <= len(slots)-requiredSlots; i++ {

		valid := true

		for j := 0; j < requiredSlots; j++ {

			if slots[i+j].Status != "available" {
				valid = false
				break
			}

			if j > 0 {
				prev := slots[i+j-1].Slot_end
				curr := slots[i+j].Slot_start

				if !prev.Equal(curr) {
					valid = false
					break
				}
			}
		}

		if valid {
			result = append(result, slots[i])
		}
	}

	return result
}

func ToSlotResponse(slot models.Slot) dto.SlotResponse {
	return dto.SlotResponse{
		Id:         slot.Id.String(),
		Slot_start: slot.Slot_start,
		Slot_end:   slot.Slot_end,
		Status:     slot.Status,
	}
}

func ToSlotResponseList(slots []models.Slot) []dto.SlotResponse {
	result := make([]dto.SlotResponse, 0, len(slots))
	for _, u := range slots {
		result = append(result, ToSlotResponse(u))
	}
	return result
}

func (s *ScheduleService) GetDoctorSchedule(doctor_id uuid.UUID) ([]models.Schedule, error) {
	return s.repo.GetScheduleByDoctor(doctor_id)
}

func ToScheduleResponse(schedule models.Schedule) dto.ScheduleDoctorResponse {
	return dto.ScheduleDoctorResponse{
		Id:                schedule.Id.String(),
		Clinic_address_id: schedule.Clinic_address_id.String(),
		Day_of_week:       schedule.Day_of_week,
		Start_time:        schedule.Start_time,
		End_time:          schedule.End_time,
	}
}

func ToScheduleResponseList(schedules []models.Schedule) []dto.ScheduleDoctorResponse {
	result := make([]dto.ScheduleDoctorResponse, 0, len(schedules))
	for _, u := range schedules {
		result = append(result, ToScheduleResponse(u))
	}
	return result
}
