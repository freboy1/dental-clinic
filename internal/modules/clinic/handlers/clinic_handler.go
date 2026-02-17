package handlers

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/clinic/models"
	"dental_clinic/internal/modules/clinic/services"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type ClinicHandler struct {
	service *services.ClinicService
	cfg     config.Config
}

func NewClinicHandler(s *services.ClinicService, cfg config.Config) *ClinicHandler {
	return &ClinicHandler{
		service: s,
		cfg:     cfg,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, ErrorResponse{Error: message})
}
// GetClinics godoc
// @Summary Get all clinics
// @Description Returns a list of all clinics
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SuccessResponse "OK"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/clinics [get]
func (h *ClinicHandler) GetClinics(w http.ResponseWriter, r *http.Request) {
	clinics, err := h.service.GetAllClinics()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if clinics == nil {
		clinics = []*models.Clinic{}
	}
	respondJSON(w, http.StatusOK, SuccessResponse{
		Data: clinics,
	})
}
// GetClinic godoc
// @Summary Get clinic by ID
// @Description Returns a single clinic by its UUID
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Invalid clinic ID format"
// @Failure 404 {object} ErrorResponse "Clinic not found"
// @Router /api/clinics/{id} [get]
func (h *ClinicHandler) GetClinic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic ID format")
		return
	}

	clinic, err := h.service.GetClinicByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Clinic not found")
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Data: clinic,
	})
}
// CreateClinic godoc
// @Summary Create a new clinic
// @Description Creates a new clinic
// @Tags Clinics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.Clinic true "Clinic data"
// @Success 201 {object} SuccessResponse "Created"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/clinics [post]
func (h *ClinicHandler) CreateClinic(w http.ResponseWriter, r *http.Request) {
	var clinic models.Clinic

	if err := json.NewDecoder(r.Body).Decode(&clinic); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	createdClinic, err := h.service.CreateClinic(&clinic)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, SuccessResponse{
		Message: "Clinic created successfully",
		Data:    createdClinic,
	})

}
// UpdateClinic godoc
// @Summary Update clinic
// @Description Updates an existing clinic by UUID
// @Tags Clinics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Param request body models.Clinic true "Updated clinic data"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Clinic not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/clinics/{id} [put]
func (h *ClinicHandler) UpdateClinic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic ID format")
		return
	}

	var clinic models.Clinic

	if err := json.NewDecoder(r.Body).Decode(&clinic); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	updatedClinic, err := h.service.UpdateClinic(id, &clinic)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Clinic updated successfully",
		Data:    updatedClinic,
	})

}
// DeleteClinic godoc
// @Summary Delete clinic
// @Description Deletes a clinic by UUID
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Invalid clinic ID format"
// @Failure 404 {object} ErrorResponse "Clinic not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/clinics/{id} [delete]
func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic ID format")
		return
	}

	if err := h.service.DeleteClinic(id); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Clinic deleted successfully",
	})

}
