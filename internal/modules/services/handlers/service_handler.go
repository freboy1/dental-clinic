package handlers

import (
	"database/sql"
	"dental_clinic/internal/modules/services/dto"
	"dental_clinic/internal/modules/services/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type ServiceHandler struct {
	service *services.ServiceService
}

func NewServiceHandler(s *services.ServiceService) *ServiceHandler {
	return &ServiceHandler{service: s}
}

// CreateService godoc
// @Summary Create a new service
// @Description Creates a new dental service (e.g. teeth cleaning, tooth extraction)
// @Tags Services
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body dto.CreateServiceRequest true "Service creation data"
// @Success 200 {object} dto.ServiceActionResponse
// @Failure 400 {object} dto.ServiceActionResponse
// @Router /api/services [post]
func (h *ServiceHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	response := dto.ServiceActionResponse{
		Success:   "0",
		Message:   "",
		ServiceID: "",
	}

	var req dto.CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	service, err := h.service.CreateService(req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "service created successfully"
	response.ServiceID = service.Id.String()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAllServices godoc
// @Summary Get all services
// @Description Returns a list of all dental services
// @Tags Services
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.ServiceResponse
// @Failure 500 {object} map[string]string
// @Router /api/services [get]
func (h *ServiceHandler) GetAllServices(w http.ResponseWriter, r *http.Request) {
	servicesList, err := h.service.GetAllServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToServiceResponseList(servicesList))
}

// GetServicesByClinic godoc
// @Summary Get services by clinic
// @Description Returns all services for a specific clinic
// @Tags Services
// @Security BearerAuth
// @Produce json
// @Param clinic_id path string true "Clinic ID"
// @Success 200 {array} dto.ServiceResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinics/{clinic_id}/services [get]
func (h *ServiceHandler) GetServicesByClinic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clinicID := vars["clinic_id"]

	servicesList, err := h.service.GetServicesByClinic(clinicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToServiceResponseList(servicesList))
}

// GetServiceByID godoc
// @Summary Get service by ID
// @Description Returns a single service by ID
// @Tags Services
// @Security BearerAuth
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} dto.ServiceResponse
// @Failure 404 {object} map[string]string
// @Router /api/services/{id} [get]
func (h *ServiceHandler) GetServiceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	service, err := h.service.GetServiceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToServiceResponse(*service))
}

// UpdateService godoc
// @Summary Update service
// @Description Updates an existing service
// @Tags Services
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Service ID"
// @Param request body dto.UpdateServiceRequest true "Service update data"
// @Success 200 {object} dto.ServiceResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/services/{id} [put]
func (h *ServiceHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	service, err := h.service.UpdateService(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToServiceResponse(*service))
}

// DeleteService godoc
// @Summary Delete service
// @Description Deletes a service by ID
// @Tags Services
// @Security BearerAuth
// @Param id path string true "Service ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/services/{id} [delete]
func (h *ServiceHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteService(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "service not found" {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Service deleted successfully"})
}