package handlers

import (
	"database/sql"
	"dental_clinic/internal/modules/doctor/dto"
	"dental_clinic/internal/modules/doctor/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type DoctorHandler struct {
	service *services.DoctorService
}

func NewDoctorHandler(s *services.DoctorService) *DoctorHandler {
	return &DoctorHandler{service: s}
}

// CreateDoctor godoc
// @Summary Create a new doctor
// @Description Creates a new doctor profile linked to a user
// @Tags Doctors
// @Accept  json
// @Produce  json
// @Param request body dto.CreateDoctorRequest true "Doctor creation data"
// @Success 200 {object} dto.DoctorActionResponse
// @Failure 400 {object} dto.DoctorActionResponse
// @Router /api/doctors [post]
func (h *DoctorHandler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	response := dto.DoctorActionResponse{
		Success:  "0",
		Message:  "",
		DoctorID: "",
	}

	var req dto.CreateDoctorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	doctor, err := h.service.CreateDoctor(req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "doctor created successfully"
	response.DoctorID = doctor.Id.String()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAllDoctors godoc
// @Summary Get all doctors
// @Description Returns a list of all doctors
// @Tags Doctors
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.DoctorResponse
// @Failure 500 {object} map[string]string
// @Router /api/doctors [get]
func (h *DoctorHandler) GetAllDoctors(w http.ResponseWriter, r *http.Request) {
	doctors, err := h.service.GetAllDoctors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToDoctorResponseList(doctors))
}

// GetDoctorByID godoc
// @Summary Get doctor by ID
// @Description Returns a single doctor by ID
// @Tags Doctors
// @Security BearerAuth
// @Produce json
// @Param id path string true "Doctor ID"
// @Success 200 {object} dto.DoctorResponse
// @Failure 404 {object} map[string]string
// @Router /api/doctors/{id} [get]
func (h *DoctorHandler) GetDoctorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	doctor, err := h.service.GetDoctorByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToDoctorResponse(*doctor))
}

// UpdateDoctor godoc
// @Summary Update doctor
// @Description Updates an existing doctor's information
// @Tags Doctors
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Doctor ID"
// @Param request body dto.UpdateDoctorRequest true "Doctor update data"
// @Success 200 {object} dto.DoctorResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/doctors/{id} [put]
func (h *DoctorHandler) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.UpdateDoctorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	doctor, err := h.service.UpdateDoctor(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToDoctorResponse(*doctor))
}

// DeleteDoctor godoc
// @Summary Delete doctor
// @Description Deletes a doctor by ID
// @Tags Doctors
// @Security BearerAuth
// @Param id path string true "Doctor ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/doctors/{id} [delete]
func (h *DoctorHandler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteDoctor(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "doctor not found" {
			http.Error(w, "Doctor not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Doctor deleted successfully"})
}