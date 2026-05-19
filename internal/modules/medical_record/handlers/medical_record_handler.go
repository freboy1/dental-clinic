package handlers

import (
	"dental_clinic/internal/modules/medical_record/dto"
	"dental_clinic/internal/modules/medical_record/models"
	"dental_clinic/internal/modules/medical_record/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
// @Accept multipart/form-data
// @Produce  json
// @Param id path string true "MedicalRecord ID"
// @Param diagnosis formData string false "Diagnosis"
// @Param notes formData string false "Notes"
// @Param is_checked formData bool false "Is checked"
// @Param files formData file false "Files"
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

	// парсим multipart form (макс 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	req := dto.UpdateMedicalRecordRequest{
		Diagnosis:  r.FormValue("diagnosis"),
		Notes:      r.FormValue("notes"),
		Is_checked: r.FormValue("is_checked") == "true",
	}

	// сохраняем файлы локально
	files := r.MultipartForm.File["files"]
	var medicalFiles []models.MedicalFile
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// создаём папку если нет
		os.MkdirAll("./uploads/medical_records", os.ModePerm)

		filePath := fmt.Sprintf("./uploads/medical_records/%d_%s", time.Now().UnixNano(), fileHeader.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			continue
		}
		defer dst.Close()
		io.Copy(dst, file)

		medicalFile := models.MedicalFile{
			Filename: fileHeader.Filename,
			FilePath: filePath,
		}

		medicalFiles = append(medicalFiles, medicalFile)
	}

	_, err := h.service.UpdateMedicalRecord(id, req, medicalFiles)
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
