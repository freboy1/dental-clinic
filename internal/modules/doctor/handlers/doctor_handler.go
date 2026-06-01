package handlers

import (
	"database/sql"
	"dental_clinic/internal/config"
	"dental_clinic/internal/middleware"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dental_clinic/internal/modules/doctor/dto"
	"dental_clinic/internal/modules/doctor/services"

	"github.com/gorilla/mux"
	// "fmt"
)

type DoctorHandler struct {
	service *services.DoctorService
	cfg     config.Config
}

func NewDoctorHandler(s *services.DoctorService, cfg config.Config) *DoctorHandler {
	return &DoctorHandler{
		service: s,
		cfg:     cfg,
	}
}

// CreateDoctor godoc
// @Summary Create a new doctor
// @Description Creates a new doctor profile linked to a user
// @Tags Doctors
// @Security BearerAuth
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
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	result, err := h.service.CreateDoctor(req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "doctor created successfully"
	response.DoctorID = result.Doctor.Id.String()
	response.ConfirmationCode = result.ConfirmationCode

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
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
	_ = json.NewEncoder(w).Encode(services.ToDoctorResponseList(doctors))
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
	_ = json.NewEncoder(w).Encode(services.ToDoctorResponse(*doctor))
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
	_ = json.NewEncoder(w).Encode(services.ToDoctorResponse(*doctor))
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
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Doctor deleted successfully"})
}

// UpdateDoctorPhoto godoc
// @Summary Update doctor photo
// @Description Uploads and stores doctor photo locally
// @Tags Doctors
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Doctor ID"
// @Param photo formData file true "Doctor photo"
// @Success 200 {object} dto.DoctorActionResponse
// @Failure 400 {object} dto.DoctorActionResponse
// @Router /api/doctors/{id}/photo [post]
func (h *DoctorHandler) UpdateDoctorPhoto(w http.ResponseWriter, r *http.Request) {
	response := dto.DoctorActionResponse{
		Success: "0",
		Message: "",
	}

	doctorID := mux.Vars(r)["id"]
	currentDoctor, err := h.service.GetDoctorByID(doctorID)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		response.Message = "photo is required"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		response.Message = "failed to read photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	contentType := http.DetectContentType(buffer[:n])
	if !strings.HasPrefix(contentType, "image/") {
		response.Message = "photo must be an image"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		response.Message = "failed to read photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if err := os.MkdirAll("./uploads/doctors", os.ModePerm); err != nil {
		response.Message = "failed to create upload directory"
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s_%d%s", doctorID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(".", "uploads", "doctors", filename)
	dst, err := os.Create(filePath)
	if err != nil {
		response.Message = "failed to save photo"
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		_ = dst.Close()
		response.Message = "failed to save photo"
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	_ = dst.Close()

	photoURL := "/uploads/doctors/" + filename
	if err := h.service.UpdateDoctorPhoto(doctorID, dto.DoctorPhotoRequest{PhotoURL: photoURL}); err != nil {
		_ = os.Remove(filePath)
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if currentDoctor.PhotoURL != "" {
		_ = os.Remove(filepath.FromSlash("." + currentDoctor.PhotoURL))
	}

	response.Success = "1"
	response.Message = "doctor photo updated successfully"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// DeleteDoctorPhoto godoc
// @Summary Delete doctor photo
// @Description Clears doctor photo URL
// @Tags Doctors
// @Security BearerAuth
// @Produce json
// @Param id path string true "Doctor ID"
// @Success 200 {object} dto.DoctorActionResponse
// @Failure 400 {object} dto.DoctorActionResponse
// @Router /api/doctors/{id}/photo [delete]
func (h *DoctorHandler) DeleteDoctorPhoto(w http.ResponseWriter, r *http.Request) {
	response := dto.DoctorActionResponse{
		Success: "0",
		Message: "",
	}

	doctorID := mux.Vars(r)["id"]
	doctor, err := h.service.GetDoctorByID(doctorID)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if err := h.service.DeleteDoctorPhoto(doctorID); err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if doctor.PhotoURL != "" {
		_ = os.Remove(filepath.FromSlash("." + doctor.PhotoURL))
	}

	response.Success = "1"
	response.Message = "doctor photo deleted successfully"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// GetDoctorByIdMedicalRecords godoc
// @Summary Get doctor medical records
// @Description Get doctor medical records by ID
// @Tags Doctors
// @Security BearerAuth
// @Param id path string true "Doctor ID"
// @Success 200 {array} dto.GetMedicalRecordDoctorResponse
// @Failure 404 {object} map[string]string
// @Router /api/doctors/medical-records/{id} [get]
func (h *DoctorHandler) GetDoctorByIdMedicalRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	responses, err := h.service.GetDoctorByIdMedicalRecords(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(responses)
}

// GetDoctorMedicalRecords godoc
// @Summary Get doctor medical records
// @Description Get doctor medical records
// @Tags Doctors
// @Security BearerAuth
// @Success 200 {array} dto.GetMedicalRecordDoctorResponse
// @Failure 404 {object} map[string]string
// @Router /api/doctors-test/my-medical-records [get]
func (h *DoctorHandler) GetDoctorMedicalRecords(w http.ResponseWriter, r *http.Request) {

	user_id, err := middleware.GetUserID(r, h.cfg.JWTSecret)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	responses, err := h.service.GetDoctorByUserIdMedicalRecords(user_id.String())
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responses)
}
