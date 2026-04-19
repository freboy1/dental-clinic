package handlers

import (
	//"dental_clinic/internal/modules/medical_record/dto"
	"dental_clinic/internal/modules/medical_record/services"
	"net/http"
)

type MedicalRecordHandler struct {
	service *services.MedicalRecordService
}

func NewDoctorHandler(s *services.MedicalRecordService) *MedicalRecordHandler {
	return &MedicalRecordHandler{service: s}
}

// UpdateDoctor godoc
// @Summary Update MedicalRecord
// @Description Updates an existing MedicalRecord's information
// @Tags MedicalRecord
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "MedicalRecord ID"
// @Param request body dto.UpdateDoctorRequest true "Doctor update data"
// @Success 200 {object} dto.DoctorResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/medical-records/{id} [put]
func (h *MedicalRecordHandler) UpdateMedicalRecord(w http.ResponseWriter, r *http.Request) {

}
