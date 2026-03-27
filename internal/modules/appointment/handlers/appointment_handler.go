package handlers

import (
	// "database/sql"
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/appointment/dto"
	"dental_clinic/internal/modules/appointment/services"
	"dental_clinic/internal/utils"
	"encoding/json"
	// "errors"
	"net/http"
	// "strings"

	// "github.com/gorilla/mux"
)

type AppointmentHandler struct {
	service *services.AppointmentService
	cfg config.Config
}

func NewAppointmentHandler(s *services.AppointmentService, cfg config.Config) *AppointmentHandler {
	return &AppointmentHandler{
		service: s,
		cfg: cfg,
	}
}


// CreateAppointment godoc
// @Summary Create new appointment
// @Description Creates a new appointment
// @Tags Appointment
// @Accept  json
// @Produce  json
// @Param request body dto.CreateAppointmentRequest true "Appointment registration data"
// @Success 200 {object} dto.CreateAppointmentResponse
// @Failure 400 {object} dto.CreateAppointmentResponse
// @Router /api/appointment [post]
func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	response := dto.CreateAppointmentResponse{
		Success: "0",
		Message: "",
		Appointment_id:  "",
	}

	var req dto.CreateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	_ = r
	tokenStr := utils.GetToken(r)

	appointment, err := h.service.CreateAppointment(tokenStr, req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Success = "1"
	response.Message = "successfully created"
	response.Appointment_id = appointment.Id.String()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}




// GetAppointments godoc
// @Summary get appointments
// @Description Get appointments
// @Tags Appointment
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.GetAppointmentsResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/appointment [get]
func (h *AppointmentHandler) GetAllAppointments(w http.ResponseWriter, r *http.Request) {
	
	appointments, err := h.service.GetAllAppointments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(services.ToAppointmentResponseList(appointments))

}
