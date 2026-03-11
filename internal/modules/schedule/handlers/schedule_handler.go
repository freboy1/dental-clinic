package handlers

import (
	"dental_clinic/internal/config"
	"fmt"

	// "dental_clinic/internal/modules/schedule/models"
	"dental_clinic/internal/modules/schedule/dto"
	"dental_clinic/internal/modules/schedule/services"
	"encoding/json"

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