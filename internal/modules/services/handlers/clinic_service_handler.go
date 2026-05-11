package handlers

import (
	"dental_clinic/internal/modules/services/dto"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// AddServiceToClinic godoc
// @Summary Add a new service
// @Description Add a new dental service (e.g. teeth cleaning, tooth extraction)
// @Tags Services
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Clinic ID (UUID)"
// @Param request body dto.AddServiceRequest true "Service creation data"
// @Success 200 {object} dto.ServiceActionResponse
// @Failure 400 {object} dto.ServiceActionResponse
// @Router /api/add-clinics/{id}/services [post]
func (h *ServiceHandler) AddServiceToClinic(w http.ResponseWriter, r *http.Request) {
	response := dto.ServiceActionResponse{
		Success:   "0",
		Message:   "",
		ServiceID: "",
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.AddServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	service, err := h.service.AddServiceToClinic(id, req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "service added successfully"
	response.ServiceID = service.Id.String()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
