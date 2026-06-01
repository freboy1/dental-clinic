package handlers

import (
	"encoding/json"
	"net/http"

	"dental_clinic/internal/modules/clinic_admin/dto"
	"dental_clinic/internal/modules/clinic_admin/services"

	"github.com/gorilla/mux"
)

type ClinicAdminHandler struct {
	service *services.ClinicAdminService
}

func NewClinicAdminHandler(service *services.ClinicAdminService) *ClinicAdminHandler {
	return &ClinicAdminHandler{service: service}
}

// CreateClinicAdmin godoc
// @Summary Create clinic admin
// @Description Creates a clinic admin profile and linked user with role clinic_admin
// @Tags Clinic Admins
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateClinicAdminRequest true "Clinic admin data"
// @Success 200 {object} dto.ClinicAdminActionResponse
// @Failure 400 {object} dto.ClinicAdminActionResponse
// @Router /api/clinic-admins [post]
func (h *ClinicAdminHandler) CreateClinicAdmin(w http.ResponseWriter, r *http.Request) {
	response := dto.ClinicAdminActionResponse{Success: "0", Message: ""}

	var req dto.CreateClinicAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	defer r.Body.Close()

	admin, err := h.service.CreateClinicAdmin(req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "clinic admin created successfully"
	response.ClinicAdminID = admin.Id.String()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// GetClinicAdmins godoc
// @Summary Get clinic admins
// @Description Returns all clinic admins
// @Tags Clinic Admins
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.ClinicAdminResponse
// @Failure 500 {object} map[string]string
// @Router /api/clinic-admins [get]
func (h *ClinicAdminHandler) GetClinicAdmins(w http.ResponseWriter, r *http.Request) {
	admins, err := h.service.GetAllClinicAdmins()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(services.ToClinicAdminResponseList(admins))
}

// GetClinicAdminByID godoc
// @Summary Get clinic admin by ID
// @Description Returns one clinic admin by UUID
// @Tags Clinic Admins
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic Admin ID"
// @Success 200 {object} dto.ClinicAdminResponse
// @Failure 404 {object} map[string]string
// @Router /api/clinic-admins/{id} [get]
func (h *ClinicAdminHandler) GetClinicAdminByID(w http.ResponseWriter, r *http.Request) {
	admin, err := h.service.GetClinicAdminByID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(services.ToClinicAdminResponse(*admin))
}

// UpdateClinicAdmin godoc
// @Summary Update clinic admin
// @Description Updates clinic admin profile and linked user account
// @Tags Clinic Admins
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Clinic Admin ID"
// @Param request body dto.UpdateClinicAdminRequest true "Clinic admin data"
// @Success 200 {object} dto.ClinicAdminResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinic-admins/{id} [put]
func (h *ClinicAdminHandler) UpdateClinicAdmin(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateClinicAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	admin, err := h.service.UpdateClinicAdmin(mux.Vars(r)["id"], req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(services.ToClinicAdminResponse(*admin))
}

// DeleteClinicAdmin godoc
// @Summary Delete clinic admin
// @Description Deletes clinic admin profile and linked user account
// @Tags Clinic Admins
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic Admin ID"
// @Success 200 {object} dto.ClinicAdminActionResponse
// @Failure 400 {object} dto.ClinicAdminActionResponse
// @Router /api/clinic-admins/{id} [delete]
func (h *ClinicAdminHandler) DeleteClinicAdmin(w http.ResponseWriter, r *http.Request) {
	response := dto.ClinicAdminActionResponse{Success: "0", Message: ""}

	if err := h.service.DeleteClinicAdmin(mux.Vars(r)["id"]); err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "clinic admin deleted successfully"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
