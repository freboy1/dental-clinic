package services

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/appointment/dto"
	"dental_clinic/internal/modules/appointment/models"
	"dental_clinic/internal/modules/appointment/repository"
	"dental_clinic/internal/utils"
	"time"


	scheduleServices "dental_clinic/internal/modules/schedule/services"
	serviceServices "dental_clinic/internal/modules/services/services"

	"errors"
	// "fmt"

	"github.com/google/uuid"
)

type AppointmentService struct {
	repo repository.AppointmentRepository
	cfx  config.Config
	scheduleSrv scheduleServices.ScheduleService
	serviceSrv serviceServices.ServiceService
}

func NewAppointmentService(r repository.AppointmentRepository, cfx config.Config, scheduleSrv scheduleServices.ScheduleService, serviceSrv serviceServices.ServiceService) *AppointmentService {
	return &AppointmentService{
		repo: r,
		cfx:  cfx,
		scheduleSrv: scheduleSrv,
		serviceSrv: serviceSrv,
	}
}

func (s *AppointmentService) CreateAppointment(tokenStr string, req dto.CreateAppointmentRequest) (*models.Appointment, error) {

	claims, _ := utils.GetClaims(tokenStr, s.cfx.JWTSecret)

	userIDStr, ok := claims["UserID"].(string)
	if !ok {
		return nil, errors.New("invalid UserID type in claims")
	}

	userId, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid UserID")
	}

	doctorId, err := uuid.Parse(req.Doctor_id)
	if err != nil {
		return nil, errors.New("invalid doctorId")
	}

	clinic_addressId, err := uuid.Parse(req.Clinic_address_id)
	if err != nil {
		return nil, errors.New("invalid clinic_addressId")
	}
	
	serviceId, err := uuid.Parse(req.Service_id)
	if err != nil {
		return nil, errors.New("invalid serviceId")
	}

	service, err := s.serviceSrv.GetServiceByID(serviceId.String())
	if err != nil {
		return nil, err
	}

	slotId, err := uuid.Parse(req.Slot_id)
	if err != nil {
		return nil, errors.New("invalid slotId")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}


	slot, err := s.scheduleSrv.GetSlotById(slotId)

	requiredSlots := s.scheduleSrv.HowManySlots(service.Duration)

	rawSlots, err := s.scheduleSrv.GetAvailableSlotsByDateAndDoctorAndClinic(doctorId, clinic_addressId, date)
	if err != nil {
		return nil, err
	}

	slotsToBook, err := s.scheduleSrv.AreSlotsAvailable(rawSlots, slot.Id, requiredSlots)
	if err != nil {
		return nil, err
	}



	appointment := &models.Appointment{
		Id: uuid.New(),
		Doctor_id: doctorId,
		User_id: userId,
		Clinic_address_id: clinic_addressId,
		Service_id: serviceId,
		Start_time: slotsToBook[0].Slot_start,
		End_time: slotsToBook[len(slotsToBook)-1].Slot_end,
		Status: "booked",
		Created_at: time.Time{},
	}

	return s.repo.Create(appointment)
}

