package services

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/appointment/dto"
	"dental_clinic/internal/modules/appointment/models"
	"dental_clinic/internal/modules/appointment/repository"
	"dental_clinic/internal/utils"

	// "fmt"

	"time"

	scheduleServices "dental_clinic/internal/modules/schedule/services"
	serviceServices "dental_clinic/internal/modules/services/services"

	"errors"
	// "fmt"

	"github.com/google/uuid"
)

type AppointmentService struct {
	repo        repository.AppointmentRepository
	cfx         config.Config
	scheduleSrv scheduleServices.ScheduleService
	serviceSrv  serviceServices.ServiceService
}

func NewAppointmentService(r repository.AppointmentRepository, cfx config.Config, scheduleSrv scheduleServices.ScheduleService, serviceSrv serviceServices.ServiceService) *AppointmentService {
	return &AppointmentService{
		repo:        r,
		cfx:         cfx,
		scheduleSrv: scheduleSrv,
		serviceSrv:  serviceSrv,
	}
}

func (s *AppointmentService) CreateAppointment(tokenStr string, req dto.CreateAppointmentRequest) (*models.Appointment, error) {
	var userId uuid.UUID
	claims, _ := utils.GetClaims(tokenStr, s.cfx.JWTSecret)

	doctorId, err := uuid.Parse(req.Doctor_id)
	if err != nil {
		return nil, errors.New("invalid doctorId")
	}

	userIDStr, _ := claims["user_id"].(string)
	if userIDStr != "" {
		userId, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.New("invalid UserID")
		}
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
	if err != nil {
		return nil, err
	}

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
		Id:                uuid.New(),
		Doctor_id:         doctorId,
		User_id:           userId,
		Clinic_address_id: clinic_addressId,
		Service_id:        serviceId,
		Start_time:        slotsToBook[0].Slot_start,
		End_time:          slotsToBook[len(slotsToBook)-1].Slot_end,
		Status:            "booked",
		Created_at:        time.Time{},
		Name:              req.Name,
		Email:             req.Email,
	}

	appointment, err = s.repo.Create(appointment)
	if err != nil {
		return nil, err
	}

	_ = utils.SendEmail(&s.cfx, appointment.Email, "Appointment was created", "Appointment was created")

	return appointment, nil
}

func (s *AppointmentService) GetAllAppointments() ([]models.Appointment, error) {
	return s.repo.GetAll()
}

func (s *AppointmentService) GetAppointmentByID(id string) (*models.Appointment, error) {
	appointment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, errors.New("appointment not found")
	}
	return appointment, nil
}

func (s *AppointmentService) UpdateAppointment(id string, req dto.UpdateAppointmentRequest) (*models.Appointment, error) {
	appointment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, errors.New("appointment not found")
	}

	if req.Doctor_id != "" {
		doctorId, err := uuid.Parse(req.Doctor_id)
		if err != nil {
			return nil, errors.New("invalid doctor_id")
		}
		appointment.Doctor_id = doctorId
	}

	if req.Clinic_address_id != "" {
		clinicAddressId, err := uuid.Parse(req.Clinic_address_id)
		if err != nil {
			return nil, errors.New("invalid clinic_address_id")
		}
		appointment.Clinic_address_id = clinicAddressId
	}

	if req.Service_id != "" {
		serviceId, err := uuid.Parse(req.Service_id)
		if err != nil {
			return nil, errors.New("invalid service_id")
		}
		appointment.Service_id = serviceId
	}

	if req.Start_time != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", req.Start_time)
		if err != nil {
			return nil, errors.New("invalid start_time format, use: 2006-01-02 15:04:05")
		}
		appointment.Start_time = startTime
	}

	if req.End_time != "" {
		endTime, err := time.Parse("2006-01-02 15:04:05", req.End_time)
		if err != nil {
			return nil, errors.New("invalid end_time format, use: 2006-01-02 15:04:05")
		}
		appointment.End_time = endTime
	}

	if req.Status != "" {
		appointment.Status = req.Status
	}

	if req.Name != "" {
		appointment.Name = req.Name
	}

	if req.Email != "" {
		appointment.Email = req.Email
	}

	return s.repo.Update(appointment)
}

func ToAppointmentResponse(appointment models.Appointment) dto.GetAppointmentsResponse {
	return dto.GetAppointmentsResponse{
		Id:                appointment.Id.String(),
		Doctor_id:         appointment.Doctor_id.String(),
		Clinic_address_id: appointment.Clinic_address_id.String(),
		Service_id:        appointment.Service_id.String(),
		User_id:           appointment.User_id.String(),
		Start_time:        appointment.Start_time.Format("2006-01-02 15:04:05"),
		End_time:          appointment.End_time.Format("2006-01-02 15:04:05"),
		Status:            appointment.Status,
		Name:              appointment.Name,
		Email:             appointment.Email,
	}
}

func ToAppointmentResponseList(appointments []models.Appointment) []dto.GetAppointmentsResponse {
	result := make([]dto.GetAppointmentsResponse, 0, len(appointments))
	for _, u := range appointments {
		result = append(result, ToAppointmentResponse(u))
	}
	return result
}

func (s *AppointmentService) DeleteAppointment(id string) error {
	appointment, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if appointment == nil {
		return errors.New("appointmnet not found")
	}
	return s.repo.Delete(id)
}

func (s *AppointmentService) GetMyAppointments(tokenStr string) ([]models.Appointment, error) {
	claims, _ := utils.GetClaims(tokenStr, s.cfx.JWTSecret)

	userIDStr, _ := claims["user_id"].(string)

	if userIDStr == "" {
		return nil, errors.New("No user Id")
	}
	userId, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid UserID")
	}

	return s.repo.GetMyAppointments(userId.String())
}
