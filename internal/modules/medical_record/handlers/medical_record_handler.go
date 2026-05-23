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

		buffer := make([]byte, 512)

		_, err = file.Read(buffer)
		if err != nil {
			continue
		}

		mimeType := http.DetectContentType(buffer)

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			continue
		}

		// создаём папку если нет
		os.MkdirAll("./uploads/medical_records", os.ModePerm)

		filePath := fmt.Sprintf("./uploads/medical_records/%d_%s", time.Now().UnixNano(), fileHeader.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			continue
		}

		_, err = io.Copy(dst, file)

		dst.Close()
		file.Close()

		if err != nil {
			continue
		}

		medicalFile := models.MedicalFile{
			Filename: fileHeader.Filename,
			FilePath: filePath,
			MimeType: mimeType,
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

	medical_files, err := h.service.GetMedicalRecordFiles(id)
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
	response.Files = medical_files

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// PreviewMedicalFile godoc
// @Summary Preview medical file
// @Description Preview medical file in browser
// @Tags MedicalRecord
// @Security BearerAuth
// @Produce octet-stream
// @Param id path string true "File ID"
// @Success 200 {file} file
// @Failure 404 {object} map[string]string
// @Router /api/files/medical-records/{id} [get]
func (h *MedicalRecordHandler) GetPreviewMedicalRecordFile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	file, err := h.service.GetFileByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "file not found",
		})
		return
	}

	w.Header().Set("Content-Type", file.MimeType)

	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`inline; filename="%s"`, file.Filename),
	)

	http.ServeFile(w, r, file.FilePath)
}

// DownloadMedicalFile godoc
// @Summary Download medical file
// @Description Download medical file
// @Tags MedicalRecord
// @Security BearerAuth
// @Produce octet-stream
// @Param id path string true "File ID"
// @Success 200 {file} file
// @Failure 404 {object} map[string]string
// @Router /api/files/medical-records/{id}/download [get]
func (h *MedicalRecordHandler) DownloadMedicalRecordFile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	file, err := h.service.GetFileByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "file not found",
		})
		return
	}

	w.Header().Set("Content-Type", file.MimeType)

	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, file.Filename),
	)

	http.ServeFile(w, r, file.FilePath)
}

// DeleteMedicalFile godoc
// @Summary Delete medical file
// @Description Delete medical file by ID
// @Tags MedicalRecord
// @Security BearerAuth
// @Produce octet-stream
// @Param id path string true "File ID"
// @Success 200 {file} file
// @Failure 404 {object} map[string]string
// @Router /api/files/medical-records/{id} [delete]
func (h *MedicalRecordHandler) DeleteRecordFile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	file, err := h.service.GetFileByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "file not found",
		})
		return
	}

	err = os.Remove(file.FilePath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to delete file from storage",
		})
		return
	}

	err = h.service.DeleteFile(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to delete file record",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
