package handlers

import (
	"dental_clinic/internal/config"
	// "fmt"

	// "dental_clinic/internal/modules/schedule/models"
	"dental_clinic/internal/modules/schedule/dto"
	"dental_clinic/internal/modules/schedule/services"
	"encoding/json"
	"time"

	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ScheduleHandler struct {
	service *services.ScheduleService
	cfg     config.Config
}

func NewScheduleHandler(s *services.ScheduleService, cfg config.Config) *ScheduleHandler {
	return &ScheduleHandler{
		service: s,
		cfg:     cfg,
	}
}


// CreateSchedule godoc
// @Summary Create new schedule
// @Description Creates a new schedule
// @Tags Schedule
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param doctorId path string true "Doctor ID (UUID)"
// @Param request body dto.CreateScheduleRequest true "Schedule registration data"
// @Success 200 {object} dto.CreateScheduleResponse
// @Failure 400 {object} dto.CreateScheduleResponse
// @Router /api/schedule/doctors/{doctorId}/working-hours [post]
func (h *ScheduleHandler) CreateDoctorSchedule(w http.ResponseWriter, r *http.Request) {
	response := dto.CreateScheduleResponse{
		Success: "0",
		Message: "",
		Schedule_id:  "",
	}
	vars := mux.Vars(r)
	doctorIDStr := vars["doctorId"]
 
	doctor_id, err := uuid.Parse(doctorIDStr)
	if err != nil {
		response.Message = "Invalid doctorId"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req dto.CreateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	schedule, err := h.service.CreateSchedule(doctor_id, req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Success = "1"
	response.Message = "successfully created"
	response.Schedule_id = schedule.Id.String()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}




// GenerateSlots godoc
// @Summary Generate new slots
// @Description Generate a new slots
// @Tags Schedule
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body dto.GenerateSlotsRequest true "Generate slots data"
// @Success 200 {object} dto.ScheduleResponse
// @Failure 400 {object} dto.ScheduleResponse
// @Router /api/schedule/generate [post]
func (h *ScheduleHandler) GenerateSlots(w http.ResponseWriter, r *http.Request) {
	response := dto.ScheduleResponse{
		Success: "0",
		Message: "",
	}
	

	var req dto.GenerateSlotsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := h.service.GenerateSlots(req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Success = "1"
	response.Message = "successfully generated"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}




// Get Slots godoc
// @Summary Get available slots
// @Description Get available slots
// @Tags Schedule
// @Accept  json
// @Produce  json
// Query parameters
// @Param doctor_id query string true "Doctor ID (UUID)"
// @Param service_id query string true "Service ID (UUID)"
// @Param clinic_address_id query string true "Clinic Address ID (UUID)"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {array} dto.SlotResponse
// @Failure 400 {array} dto.SlotResponse
// @Router /api/schedule/available-slots [get]
func (h *ScheduleHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	
	query := r.URL.Query()

	doctorIDStr := query.Get("doctor_id")
	serviceIDStr := query.Get("service_id")
	clinic_addressIDStr := query.Get("clinic_address_id")
	dateStr := query.Get("date")

	if doctorIDStr == "" || serviceIDStr == "" || dateStr == "" || clinic_addressIDStr == "" {
		http.Error(w, "missing query parameters", http.StatusBadRequest)
		return
	}

	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		http.Error(w, "invalid doctor_id", http.StatusBadRequest)
		return
	}

	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		http.Error(w, "invalid service_id", http.StatusBadRequest)
		return
	}

	clinic_addressID, err := uuid.Parse(clinic_addressIDStr)
	if err != nil {
		http.Error(w, "invalid service_id", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	slots, err := h.service.GetAvailableSlots(doctorID, serviceID, clinic_addressID, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToSlotResponseList(slots))
}




// GetDoctorSchedule godoc
// @Summary Get doctor schedule
// @Description Returns a doctor schedule
// @Tags Schedule
// @Security BearerAuth
// @Produce json
// @Param doctorId path string true "Doctor ID (UUID)"
// @Success 200 {array} dto.ScheduleDoctorResponse
// @Failure 401 {object} map[string]string
// @Router /api/schedule/doctors/{doctorId}/working-hours [get]
func (h *ScheduleHandler) GetDoctorSchedule(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	doctorIDStr := vars["doctorId"]

	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		http.Error(w, "invalid doctorId", http.StatusBadRequest)
		return
	}

	schedules, err := h.service.GetDoctorSchedule(doctorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(services.ToScheduleResponseList(schedules))
}




