package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"dental_clinic/internal/config"
	aiDto "dental_clinic/internal/modules/ai_assistant/dto"
	"dental_clinic/internal/modules/ai_assistant/models"
	"dental_clinic/internal/modules/ai_assistant/repository"
	appointmentDto "dental_clinic/internal/modules/appointment/dto"
	appointmentServices "dental_clinic/internal/modules/appointment/services"
	scheduleModels "dental_clinic/internal/modules/schedule/models"
	scheduleServices "dental_clinic/internal/modules/schedule/services"
	userServices "dental_clinic/internal/modules/user/services"
	"dental_clinic/internal/utils"

	"github.com/google/uuid"
)

type AIAssistantService struct {
	cfg            config.Config
	repo           repository.AIAssistantRepository
	llm            LLMClient
	appointmentSrv appointmentServices.AppointmentService
	scheduleSrv    scheduleServices.ScheduleService
	userSrv        userServices.UserService
}

func NewAIAssistantService(
	cfg config.Config,
	repo repository.AIAssistantRepository,
	llm LLMClient,
	appointmentSrv appointmentServices.AppointmentService,
	scheduleSrv scheduleServices.ScheduleService,
	userSrv userServices.UserService,
) *AIAssistantService {
	return &AIAssistantService{
		cfg:            cfg,
		repo:           repo,
		llm:            llm,
		appointmentSrv: appointmentSrv,
		scheduleSrv:    scheduleSrv,
		userSrv:        userSrv,
	}
}

func (s *AIAssistantService) ProcessMessage(userID uuid.UUID, tokenStr string, req aiDto.ChatRequest, ctx context.Context) (aiDto.ChatResponse, error) {
	if req.ChoiceID == "" && isResetCommand(req.Message) {
		return s.ResetBooking(userID)
	}

	session, err := s.repo.GetOrCreateSession(userID)
	if err != nil {
		return aiDto.ChatResponse{}, err
	}

	response := aiDto.ChatResponse{
		SessionID: session.Id.String(),
	}

	if strings.TrimSpace(req.Message) == "" && strings.TrimSpace(req.ChoiceID) == "" {
		response.Reply = "Please write what appointment you want to book."
		return response, nil
	}

	messageContent := req.Message
	if req.ChoiceID != "" {
		messageContent = req.ChoiceType + ":" + req.ChoiceID
	}
	if err := s.repo.SaveMessage(session.Id, "user", messageContent); err != nil {
		return response, err
	}

	state, err := s.repo.GetOrCreateState(userID)
	if err != nil {
		return response, err
	}

	extraction := BookingExtraction{}
	if req.ChoiceID != "" {
		if err := s.applyChoice(state, req.ChoiceType, req.ChoiceID); err != nil {
			return response, err
		}
	} else {
		history, err := s.repo.GetRecentMessages(session.Id, 10)
		if err != nil {
			return response, err
		}

		extraction, err = s.llm.ExtractBookingInfo(req.Message, *state, history)
		if err != nil {
			return response, err
		}
		mergeState(state, extraction)
	}

	response, err = s.nextStep(userID, tokenStr, session.Id, state, extraction, ctx)
	if err != nil {
		return response, err
	}

	_ = s.repo.SaveMessage(session.Id, "assistant", response.Reply)
	return response, nil
}

func (s *AIAssistantService) ResetBooking(userID uuid.UUID) (aiDto.ChatResponse, error) {
	if err := s.repo.ClearState(userID); err != nil {
		return aiDto.ChatResponse{}, err
	}

	session, err := s.repo.CreateSession(userID)
	if err != nil {
		return aiDto.ChatResponse{}, err
	}

	state := &models.BookingState{
		UserID: userID.String(),
		Step:   "collect_service",
	}
	if err := s.repo.SaveState(state); err != nil {
		return aiDto.ChatResponse{}, err
	}

	response := aiDto.ChatResponse{
		Reply:     "Booking flow reset. Which dental service do you want?",
		SessionID: session.Id.String(),
		State:     *state,
	}
	_ = s.repo.SaveMessage(session.Id, "assistant", response.Reply)
	return response, nil
}

func isResetCommand(message string) bool {
	switch strings.ToLower(strings.TrimSpace(message)) {
	case "reset", "start over", "restart", "begin again", "new appointment", "начать заново", "сброс", "сбросить", "заново", "новая запись":
		return true
	default:
		return false
	}
}

func mergeState(state *models.BookingState, extraction BookingExtraction) {
	if extraction.DoctorID != "" {
		state.DoctorID = extraction.DoctorID
	}
	if extraction.ServiceID != "" {
		state.ServiceID = extraction.ServiceID
	}
	if extraction.ClinicAddressID != "" {
		state.ClinicAddressID = extraction.ClinicAddressID
	}
	if extraction.Date != "" {
		state.Date = extraction.Date
	}
	if extraction.Time != "" {
		state.Time = normalizeTime(extraction.Time)
	}
}

func (s *AIAssistantService) applyChoice(state *models.BookingState, choiceType, choiceID string) error {
	choiceType = strings.TrimSpace(choiceType)
	choiceID = strings.TrimSpace(choiceID)
	if choiceID == "" {
		return errors.New("choice_id is required")
	}

	switch choiceType {
	case "service":
		state.ServiceID = choiceID
		state.ClinicAddressID = ""
		state.DoctorID = ""
		state.Time = ""
	case "clinic":
		state.ClinicAddressID = choiceID
		state.DoctorID = ""
		state.Time = ""
	case "doctor":
		state.DoctorID = choiceID
		state.Time = ""
	case "slot":
		slotID, err := uuid.Parse(choiceID)
		if err != nil {
			return errors.New("invalid slot choice")
		}
		slot, err := s.scheduleSrv.GetSlotById(slotID)
		if err != nil {
			return err
		}
		state.Time = slot.Slot_start.Format("15:04:05")
	case "date":
		if _, err := time.Parse("2006-01-02", choiceID); err != nil {
			return errors.New("invalid date choice")
		}
		state.Date = choiceID
		state.Time = ""
	default:
		return errors.New("invalid choice_type")
	}

	return nil
}

func (s *AIAssistantService) nextStep(userID uuid.UUID, tokenStr string, sessionID uuid.UUID, state *models.BookingState, extraction BookingExtraction, ctx context.Context) (aiDto.ChatResponse, error) {
	response := aiDto.ChatResponse{
		SessionID: sessionID.String(),
		State:     *state,
	}

	if state.ServiceID == "" {
		serviceQuery := strings.TrimSpace(extraction.ServiceQuery)
		if serviceQuery == "" {
			state.Step = "collect_service"
			response.Reply = "Which dental service do you want?"
			response.ChoiceRequired = false
			response.State = *state
			return response, s.repo.SaveState(state)
		}

		services, err := s.repo.SearchServices(serviceQuery)
		if err != nil {
			return response, err
		}
		if len(services) == 0 {
			state.Step = "collect_service"
			response.Reply = "I could not find that service. Please choose another service."
			response.State = *state
			return response, s.repo.SaveState(state)
		}
		if len(services) == 1 {
			state.ServiceID = services[0].Id
		} else {
			state.Step = "collect_service"
			response.Reply = "Which dental service do you want?"
			response.ChoiceRequired = true
			response.ChoiceType = "service"
			response.Services = services
			response.State = *state
			return response, s.repo.SaveState(state)
		}
	}

	if state.ClinicAddressID == "" {
		clinics, err := s.repo.GetClinicOptions(state.ServiceID)
		if err != nil {
			return response, err
		}
		if len(clinics) == 0 {
			state.Step = "collect_clinic"
			response.Reply = "I could not find clinics for this service."
			response.State = *state
			return response, s.repo.SaveState(state)
		}
		if len(clinics) == 1 {
			state.ClinicAddressID = clinics[0].ClinicAddressID
		} else {
			state.Step = "collect_clinic"
			response.Reply = "Which clinic do you prefer?"
			response.ChoiceRequired = true
			response.ChoiceType = "clinic"
			response.Clinics = clinics
			response.State = *state
			return response, s.repo.SaveState(state)
		}
	}

	if state.DoctorID == "" {
		doctors, err := s.repo.GetDoctorOptions(state.ServiceID, state.ClinicAddressID)
		if err != nil {
			return response, err
		}
		if len(doctors) == 0 {
			state.Step = "collect_doctor"
			response.Reply = "I could not find available doctors for this clinic and service."
			response.State = *state
			return response, s.repo.SaveState(state)
		}
		if len(doctors) == 1 {
			state.DoctorID = doctors[0].Id
		} else {
			state.Step = "collect_doctor"
			response.Reply = "Which doctor would you prefer?"
			response.ChoiceRequired = true
			response.ChoiceType = "doctor"
			response.Doctors = doctors
			response.State = *state
			return response, s.repo.SaveState(state)
		}
	}

	if state.Date == "" {
		state.Step = "collect_date"
		response.Reply = "What date should I book?"
		response.State = *state
		return response, s.repo.SaveState(state)
	}

	slots, err := s.getAvailableSlots(*state)
	if err != nil {
		return response, err
	}

	if state.Time == "" {
		state.Step = "collect_time"
		if len(slots) == 0 {
			response.Reply = "No available slots were found for this date. Please choose another date."
			response.State = *state
			return response, s.repo.SaveState(state)
		}
		response.Reply = "What time works for you?"
		response.ChoiceRequired = true
		response.ChoiceType = "slot"
		response.AvailableSlots = toSlotResponseList(slots)
		response.State = *state
		return response, s.repo.SaveState(state)
	}

	slotID, err := findSlotIDByTime(slots, state.Time)
	if err != nil {
		state.Time = ""
		state.Step = "collect_time"
		response.Reply = "That time is not available. Please choose one of the available slots."
		response.ChoiceRequired = true
		response.ChoiceType = "slot"
		response.AvailableSlots = toSlotResponseList(slots)
		response.State = *state
		return response, s.repo.SaveState(state)
	}

	appointmentID, err := s.createAppointment(userID, tokenStr, *state, slotID, ctx)
	if err != nil {
		return response, err
	}

	response.Reply = "Appointment booked successfully."
	response.AppointmentID = appointmentID
	response.State = *state
	_ = s.repo.ClearState(userID)
	return response, nil
}

func (s *AIAssistantService) getAvailableSlots(state models.BookingState) ([]scheduleModels.Slot, error) {
	doctorID, err := uuid.Parse(state.DoctorID)
	if err != nil {
		return nil, errors.New("invalid doctor_id")
	}
	clinicAddressID, err := uuid.Parse(state.ClinicAddressID)
	if err != nil {
		return nil, errors.New("invalid clinic_address_id")
	}
	serviceID, err := uuid.Parse(state.ServiceID)
	if err != nil {
		return nil, errors.New("invalid service_id")
	}
	date, err := time.Parse("2006-01-02", state.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	return s.scheduleSrv.GetAvailableSlots(doctorID, serviceID, clinicAddressID, date)
}

func (s *AIAssistantService) createAppointment(userID uuid.UUID, tokenStr string, state models.BookingState, slotID string, ctx context.Context) (string, error) {
	user, err := s.userSrv.GetUserByID(userID.String())
	if err != nil {
		return "", err
	}

	name := ""
	email := ""
	if user != nil {
		name = user.Name
		email = user.Email
	}
	if email == "" {
		claims, _ := utils.GetClaims(tokenStr, s.cfg.JWTSecret)
		email, _ = claims["email"].(string)
	}

	appointment, err := s.appointmentSrv.CreateAppointment(tokenStr, appointmentDto.CreateAppointmentRequest{
		Doctor_id:         state.DoctorID,
		Clinic_address_id: state.ClinicAddressID,
		Service_id:        state.ServiceID,
		Slot_id:           slotID,
		Date:              state.Date,
		Name:              name,
		Email:             email,
	}, ctx)
	if err != nil {
		return "", err
	}

	return appointment.Id.String(), nil
}

func findSlotIDByTime(slots []scheduleModels.Slot, requested string) (string, error) {
	requested = normalizeTime(requested)
	for _, slot := range slots {
		if slot.Slot_start.Format("15:04:05") == requested {
			return slot.Id.String(), nil
		}
	}
	return "", fmt.Errorf("selected time is not available")
}

func normalizeTime(value string) string {
	value = strings.TrimSpace(value)
	if len(value) == 5 {
		return value + ":00"
	}
	return value
}

func toSlotResponse(slot scheduleModels.Slot) aiDto.SlotResponse {
	return aiDto.SlotResponse{
		Id:        slot.Id.String(),
		SlotStart: slot.Slot_start,
		SlotEnd:   slot.Slot_end,
		Status:    slot.Status,
	}
}

func toSlotResponseList(slots []scheduleModels.Slot) []aiDto.SlotResponse {
	result := make([]aiDto.SlotResponse, 0, len(slots))
	for _, slot := range slots {
		result = append(result, toSlotResponse(slot))
	}
	return result
}
