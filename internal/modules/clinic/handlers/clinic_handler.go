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
