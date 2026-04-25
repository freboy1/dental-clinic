package handlers

import (
	"dental_clinic/internal/modules/medical_record/dto"
	"dental_clinic/internal/modules/medical_record/services"
	"encoding/json"
	"net/http"

	//"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MedicalRecordHandler struct {
	service *services.MedicalRecordService
}

func NewMedicalRecordHandler(s *services.MedicalRecordService) *MedicalRecordHandler {
	return &MedicalRecordHandler{service: s}
}

// UpdateMedicalRecord godoc
// @Summary Update MedicalRecord
// @Description Updates an existing MedicalRecord's information
// @Tags MedicalRecord
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "MedicalRecord ID"
// @Param request body dto.UpdateMedicalRecordRequest true "Medical record update data"
// @Success 200 {object} dto.MedicalRecordResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/medical-records/{id} [put]
func (h *MedicalRecordHandler) UpdateMedicalRecord(w http.ResponseWriter, r *http.Request) {
	response := dto.MedicalRecordResponse{
		Success: "0",
		Message: "",
	}
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.UpdateMedicalRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	_, err := h.service.UpdateMedicalRecord(id, req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "successfully updated"

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// GetMedicalRecord godoc
// @Summary Get MedicalRecord
// @Description gets an existing MedicalRecord's information
// @Tags MedicalRecord
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "MedicalRecord ID"
// @Success 200 {object} dto.GetMedicalRecordResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/medical-records/{id} [get]
func (h *MedicalRecordHandler) GetMedicalRecord(w http.ResponseWriter, r *http.Request) {
	response := dto.GetMedicalRecordResponse{
		Status:  "0",
		Message: "",
	}
	vars := mux.Vars(r)
	id := vars["id"]

	medical_record, err := h.service.GetMedicalRecord(id)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = "1"
	response.Message = "successfully retrieved"
	response.Diagnosis = medical_record.Diagnosis
	response.Notes = medical_record.Notes
	response.Is_checked = medical_record.Is_checked

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
