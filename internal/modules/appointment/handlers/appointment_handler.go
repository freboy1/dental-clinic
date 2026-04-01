package handlers

import (
	"database/sql"
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/appointment/dto"
	"dental_clinic/internal/modules/appointment/services"
	"dental_clinic/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	// "strings"

	"github.com/gorilla/mux"
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


// GetAppointmentByID godoc
// @Summary Get appointment by ID
// @Description Returns a single appointment by UUID
// @Tags Appointment
// @Security BearerAuth
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Success 200 {object} dto.GetAppointmentsResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/appointment/{id} [get]
func (h *AppointmentHandler) GetAppointmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	appointment, err := h.service.GetAppointmentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToAppointmentResponse(*appointment))
}

// UpdateAppointment godoc
// @Summary Update appointment
// @Description Updates an existing appointment by UUID
// @Tags Appointment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Param request body dto.UpdateAppointmentRequest true "Updated appointment data"
// @Success 200 {object} dto.GetAppointmentsResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/appointment/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.UpdateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	appointment, err := h.service.UpdateAppointment(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToAppointmentResponse(*appointment))
}

// DeleteAppointment godoc
// @Summary Delete appointment 
// @Description Delete appointment
// @Tags Appointment
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.AppointmentResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/appointment/{id} [delete]
func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	response := dto.AppointmentResponse{
		Success:   "0",
		Message:   "",
	}
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteAppointment(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "service not found" {
			response.Message = "Service not found"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		response.Message = "Internal server error"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Success, response.Message = "1", "Successfully deleted"
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}