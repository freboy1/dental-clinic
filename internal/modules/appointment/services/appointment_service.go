package services

import (
	"context"
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/appointment/dto"
	"dental_clinic/internal/modules/appointment/models"
	"dental_clinic/internal/modules/appointment/repository"
	"dental_clinic/internal/utils"

	// "fmt"

	"time"

	clinicServices "dental_clinic/internal/modules/clinic/services"
	medical_recordServices "dental_clinic/internal/modules/medical_record/services"
	reviewServices "dental_clinic/internal/modules/reviews/services"
	scheduleServices "dental_clinic/internal/modules/schedule/services"
	serviceServices "dental_clinic/internal/modules/services/services"

	"errors"
	// "fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentService struct {
	repo              repository.AppointmentRepository
	db                *pgxpool.Pool
	cfx               config.Config
	scheduleSrv       scheduleServices.ScheduleService
	serviceSrv        serviceServices.ServiceService
	medical_recordSrv medical_recordServices.MedicalRecordService
	clinicSrv         clinicServices.ClinicService
	reviewSrv         *reviewServices.ReviewService
}

func NewAppointmentService(r repository.AppointmentRepository, db *pgxpool.Pool, cfx config.Config, scheduleSrv scheduleServices.ScheduleService, serviceSrv serviceServices.ServiceService, medical_recordSrv medical_recordServices.MedicalRecordService, clinicSrv clinicServices.ClinicService, reviewSrv *reviewServices.ReviewService) *AppointmentService {
	return &AppointmentService{
		repo:              r,
		db:                db,
		cfx:               cfx,
		scheduleSrv:       scheduleSrv,
		serviceSrv:        serviceSrv,
		medical_recordSrv: medical_recordSrv,
		clinicSrv:         clinicSrv,
		reviewSrv:         reviewSrv,
	}
}

func (s *AppointmentService) CreateAppointment(tokenStr string, req dto.CreateAppointmentRequest, ctx context.Context) (*models.Appointment, error) {
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

	clinic_id, err := s.clinicSrv.GetClinicByAddressId(clinic_addressId)
	if err != nil {
		return nil, err
	}
	service, err := s.serviceSrv.GetServiceByID(serviceId.String())
	if err != nil {
		return nil, err
	}
	serviceInfo, err := s.serviceSrv.GetByClinicIDAndServiceID(clinic_id, service.Id.String())

	if err != nil {
		return nil, err
	}

	_, err = s.serviceSrv.GetServiceByID(serviceId.String())
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

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	requiredSlots := s.scheduleSrv.HowManySlots(serviceInfo.Duration)
	//requiredSlots := 1
	rawSlots, err := s.scheduleSrv.GetAvailableSlotsByDateAndDoctorAndClinic(doctorId, clinic_addressId, date)
	if err != nil {
		return nil, err
	}

	slotsToBook, err := s.scheduleSrv.AreSlotsAvailable(rawSlots, slot.Id, requiredSlots)
	if err != nil {
		return nil, err
	}

	if len(slotsToBook) == 0 {
		return nil, errors.New("no available slots")
	}

	for _, slot := range slotsToBook {

		err := s.scheduleSrv.ChangeSlotStatusTx(slot, "booked", tx)
		if err != nil {
			return nil, err
		}

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
		Created_at:        time.Now(),
		Name:              req.Name,
		Email:             req.Email,
	}

	appointment, err = s.repo.CreateTx(appointment, tx)
	if err != nil {
		return nil, err
	}

	_, err = s.medical_recordSrv.CreateMedicalRecordTx(
		appointment.Id,
		appointment.Doctor_id,
		appointment.User_id,
		tx,
	)

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
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
		IsReviewed:        appointment.IsReviewed,
		DoctorRating:      appointment.DoctorRating,
		ClinicRating:      appointment.ClinicRating,
		ClinicComment:     appointment.ClinicComment,
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

// GetMedicalRecord
func (s *AppointmentService) GetMedicalRecord(id string) (dto.GetMedicalRecordAppointmentResponse, error) {

	response := dto.GetMedicalRecordAppointmentResponse{
		Status:  "0",
		Message: "",
	}

	medical_record, err := s.medical_recordSrv.GetMedicalRecordByAppointmentId(id)

	if err != nil {
		response.Message = err.Error()
		return response, err
	}

	response.Status = "1"
	response.Notes = medical_record.Notes
	response.Diagnosis = medical_record.Diagnosis
	response.Is_checked = medical_record.Is_checked

	return response, nil
}

func (s *AppointmentService) CreateAppointmentReview(tokenStr, appointmentId string, req dto.CreateAppointmentReviewRequest, ctx context.Context) error {
	appointmentUUID, err := uuid.Parse(appointmentId)
	if err != nil {
		return errors.New("invalid appointmentId")
	}

	claims, err := utils.GetClaims(tokenStr, s.cfx.JWTSecret)
	if err != nil {
		return err
	}

	userIDStr, _ := claims["user_id"].(string)
	if userIDStr == "" {
		return errors.New("No user Id")
	}

	userId, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid UserID")
	}

	appointment, err := s.repo.GetByID(appointmentId)
	if err != nil {
		return err
	}
	if appointment == nil {
		return errors.New("appointment not found")
	}
	if appointment.Id != appointmentUUID {
		return errors.New("appointment not found")
	}
	if appointment.User_id != uuid.Nil && appointment.User_id != userId {
		return errors.New("appointment does not belong to user")
	}
	if appointment.IsReviewed {
		return errors.New("appointment already reviewed")
	}

	clinicIDStr, err := s.clinicSrv.GetClinicByAddressId(appointment.Clinic_address_id)
	if err != nil {
		return err
	}
	clinicId, err := uuid.Parse(clinicIDStr)
	if err != nil {
		return errors.New("invalid clinicId")
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.reviewSrv.CreateAppointmentReviewTx(
		appointment.Id,
		appointment.Doctor_id,
		clinicId,
		userId,
		req.DoctorRating,
		req.ClinicRating,
		req.ClinicComment,
		tx,
	); err != nil {
		return err
	}

	if err := s.repo.MarkReviewedTx(appointmentId, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
