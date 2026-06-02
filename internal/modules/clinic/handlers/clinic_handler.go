package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/clinic/dto"
	"dental_clinic/internal/modules/clinic/models"
	"dental_clinic/internal/modules/clinic/services"
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

func saveUploadedImage(r *http.Request, fieldName, uploadDir, ownerID string) (string, string, error) {
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		return "", "", fmt.Errorf("invalid request body")
	}

	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		return "", "", fmt.Errorf("%s is required", fieldName)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return "", "", fmt.Errorf("failed to read image")
	}
	contentType := http.DetectContentType(buffer[:n])
	if !strings.HasPrefix(contentType, "image/") {
		return "", "", fmt.Errorf("%s must be an image", fieldName)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", fmt.Errorf("failed to read image")
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", "", fmt.Errorf("failed to create upload directory")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s_%d%s", ownerID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to save image")
	}
	if _, err := io.Copy(dst, file); err != nil {
		_ = dst.Close()
		return "", "", fmt.Errorf("failed to save image")
	}
	_ = dst.Close()

	return mediaURLFromPath(filePath), filePath, nil
}

func removeUploadedFile(fileURL string) {
	if fileURL == "" {
		return
	}
	_ = os.Remove(filepath.FromSlash("." + fileURL))
}

func mediaURLFromPath(filePath string) string {
	urlPath := filepath.ToSlash(filePath)
	urlPath = strings.TrimPrefix(urlPath, ".")
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	return urlPath
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
	var req dto.CreateClinicRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	clinic := &models.Clinic{
		Name:        req.Name,
		Description: req.Description,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
		IsActive:    req.IsActive,
	}

	createdClinic, err := h.service.CreateClinic(clinic)
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

	var req dto.CreateClinicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	clinic := &models.Clinic{
		Name:        req.Name,
		Description: req.Description,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
		IsActive:    req.IsActive,
	}

	updatedClinic, err := h.service.UpdateClinic(id, clinic)
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

// UpdateClinicLogo godoc
// @Summary Update clinic logo
// @Description Uploads a clinic logo image, stores it locally, and updates the clinic logo URL
// @Tags Clinics
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Param logo formData file true "Clinic logo image"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Clinic not found"
// @Router /api/clinics/{id}/logo [post]
func (h *ClinicHandler) UpdateClinicLogo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic ID format")
		return
	}

	clinic, err := h.service.GetClinicByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Clinic not found")
		return
	}

	logoURL, filePath, err := saveUploadedImage(r, "logo", "./uploads/clinics", id.String())
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateClinicLogo(id, logoURL); err != nil {
		_ = os.Remove(filePath)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	removeUploadedFile(clinic.LogoURL)
	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Clinic logo updated successfully",
		Data:    map[string]string{"logo_url": logoURL},
	})
}

// DeleteClinicLogo godoc
// @Summary Delete clinic logo
// @Description Clears the clinic logo URL and removes the local logo file when present
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Clinic not found"
// @Router /api/clinics/{id}/logo [delete]
func (h *ClinicHandler) DeleteClinicLogo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic ID format")
		return
	}

	clinic, err := h.service.GetClinicByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Clinic not found")
		return
	}

	if err := h.service.DeleteClinicLogo(id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	removeUploadedFile(clinic.LogoURL)
	respondJSON(w, http.StatusOK, SuccessResponse{Message: "Clinic logo deleted successfully"})
}

// AddClinicAddress godoc
// @Summary Add Address
// @Description Add an address by UUID
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Param request body dto.AddAddressRequest true "Address add data"
// @Success 200 {array} dto.ClinicResponse
// @Failure 404 {array} dto.ClinicResponse
// @Router /api/clinics/{id}/address [post]
func (h *ClinicHandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	response := dto.ClinicResponse{
		Success: "0",
		Message: "",
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Message = "Invalid clinic ID format"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	var req dto.AddAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	clinic, err := h.service.GetClinicByID(id)
	if err != nil {
		response.Message = "Clinic not found"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if err = h.service.AddAddress(clinic.Id, req); err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "successfully added"

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// GetClinicAddress godoc
// @Summary Get all clinic address
// @Description Returns a list of all clinic address
// @Tags Clinics
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Success 200 {array} dto.GetClinicAddressResponse
// @Failure 500 {object} map[string]string
// @Router /api/clinics/{id}/address [get]
func (h *ClinicHandler) GetClinicAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid clinic ID format", http.StatusNotFound)
		return
	}

	clinic, err := h.service.GetClinicByID(id)
	if err != nil {
		http.Error(w, "Clinic not found", http.StatusNotFound)
		return
	}

	clinicAddress, err := h.service.GetClinicAddressWithName(clinic.Id)
	if err != nil {
		http.Error(w, "Failed to get clinic address", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(services.ToClinicAddressResponseList(clinicAddress))
}

// DeleteClinicAddress godoc
// @Summary Delete Address
// @Description Delete an address by UUID
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Param addressId path string true "Address ID (UUID)"
// @Success 200 {array} dto.ClinicResponse
// @Failure 404 {array} dto.ClinicResponse
// @Router /api/clinics/{id}/address/{addressId} [delete]
func (h *ClinicHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	response := dto.ClinicResponse{
		Success: "0",
		Message: "",
	}

	vars := mux.Vars(r)
	idStr, addressIdStr := vars["id"], vars["addressId"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Message = "Invalid clinic ID format"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	addressId, err := uuid.Parse(addressIdStr)
	if err != nil {
		response.Message = "Invalid address ID format"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	clinic, err := h.service.GetClinicByID(id)
	if err != nil {
		response.Message = "Clinic not found"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if err = h.service.DeleteAddress(clinic.Id, addressId); err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "successfully deleted"

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// UpdateClinicAddressCover godoc
// @Summary Update clinic address cover image
// @Description Uploads a cover image for a clinic address, stores it locally, and updates the cover image URL
// @Tags Clinics
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Clinic address ID (UUID)"
// @Param cover formData file true "Clinic address cover image"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Clinic address not found"
// @Router /api/clinic-addresses/{id}/cover [post]
func (h *ClinicHandler) UpdateClinicAddressCover(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic address ID format")
		return
	}

	clinicAddress, err := h.service.GetClinicAddressByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Clinic address not found")
		return
	}

	coverURL, filePath, err := saveUploadedImage(r, "cover", "./uploads/clinic-addresses/covers", id.String())
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateAddressCover(id, coverURL); err != nil {
		_ = os.Remove(filePath)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	removeUploadedFile(clinicAddress.CoverImageURL)
	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Clinic address cover updated successfully",
		Data:    map[string]string{"cover_image_url": coverURL},
	})
}

// DeleteClinicAddressCover godoc
// @Summary Delete clinic address cover image
// @Description Clears the cover image URL for a clinic address and removes the local cover file when present
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic address ID (UUID)"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Clinic address not found"
// @Router /api/clinic-addresses/{id}/cover [delete]
func (h *ClinicHandler) DeleteClinicAddressCover(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic address ID format")
		return
	}

	clinicAddress, err := h.service.GetClinicAddressByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Clinic address not found")
		return
	}

	if err := h.service.DeleteAddressCover(id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	removeUploadedFile(clinicAddress.CoverImageURL)
	respondJSON(w, http.StatusOK, SuccessResponse{Message: "Clinic address cover deleted successfully"})
}

// AddClinicAddressGalleryImage godoc
// @Summary Add clinic address gallery image
// @Description Uploads an image, stores it locally, and creates a gallery image record for the clinic address
// @Tags Clinics
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Clinic address ID (UUID)"
// @Param image formData file true "Gallery image"
// @Success 201 {object} SuccessResponse "Created"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Router /api/clinic-addresses/{id}/gallery [post]
func (h *ClinicHandler) AddClinicAddressGalleryImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic address ID format")
		return
	}

	imageURL, filePath, err := saveUploadedImage(r, "image", "./uploads/clinic-addresses/gallery", id.String())
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	image, err := h.service.AddGalleryImage(id, imageURL)
	if err != nil {
		_ = os.Remove(filePath)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, SuccessResponse{
		Message: "Clinic address gallery image added successfully",
		Data: dto.ClinicAddressImageResponse{
			Id:       image.Id.String(),
			ImageURL: image.ImageURL,
		},
	})
}

// GetClinicAddressGallery godoc
// @Summary Get clinic address gallery images
// @Description Returns all gallery images for a clinic address
// @Tags Clinics
// @Produce json
// @Param id path string true "Clinic address ID (UUID)"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Router /api/clinic-addresses/{id}/gallery [get]
func (h *ClinicHandler) GetClinicAddressGallery(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic address ID format")
		return
	}

	images, err := h.service.GetGalleryImages(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := make([]dto.ClinicAddressImageResponse, 0, len(images))
	for _, image := range images {
		response = append(response, dto.ClinicAddressImageResponse{
			Id:       image.Id.String(),
			ImageURL: image.ImageURL,
		})
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Data: response})
}

// UpdateClinicAddressGalleryImage godoc
// @Summary Update clinic address gallery image
// @Description Replaces a gallery image file and updates the gallery image URL
// @Tags Clinics
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Clinic address ID (UUID)"
// @Param imageId path string true "Gallery image ID (UUID)"
// @Param image formData file true "Gallery image"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Gallery image not found"
// @Router /api/clinic-addresses/{id}/gallery/{imageId} [put]
func (h *ClinicHandler) UpdateClinicAddressGalleryImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clinicAddressID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic address ID format")
		return
	}

	imageID, err := uuid.Parse(vars["imageId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid gallery image ID format")
		return
	}

	currentImage, err := h.service.GetGalleryImage(imageID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Gallery image not found")
		return
	}
	if currentImage.ClinicAddressId != clinicAddressID {
		respondError(w, http.StatusNotFound, "Gallery image not found for clinic address")
		return
	}

	imageURL, filePath, err := saveUploadedImage(r, "image", "./uploads/clinic-addresses/gallery", currentImage.ClinicAddressId.String())
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedImage, err := h.service.UpdateGalleryImage(imageID, imageURL)
	if err != nil {
		_ = os.Remove(filePath)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	removeUploadedFile(currentImage.ImageURL)
	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Clinic address gallery image updated successfully",
		Data: dto.ClinicAddressImageResponse{
			Id:       updatedImage.Id.String(),
			ImageURL: updatedImage.ImageURL,
		},
	})
}

// DeleteClinicAddressGalleryImage godoc
// @Summary Delete clinic address gallery image
// @Description Deletes a gallery image record and removes the local image file
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic address ID (UUID)"
// @Param imageId path string true "Gallery image ID (UUID)"
// @Success 200 {object} SuccessResponse "OK"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Gallery image not found"
// @Router /api/clinic-addresses/{id}/gallery/{imageId} [delete]
func (h *ClinicHandler) DeleteClinicAddressGalleryImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clinicAddressID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid clinic address ID format")
		return
	}

	imageID, err := uuid.Parse(vars["imageId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid gallery image ID format")
		return
	}

	currentImage, err := h.service.GetGalleryImage(imageID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Gallery image not found")
		return
	}
	if currentImage.ClinicAddressId != clinicAddressID {
		respondError(w, http.StatusNotFound, "Gallery image not found for clinic address")
		return
	}

	image, err := h.service.DeleteGalleryImage(imageID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	removeUploadedFile(image.ImageURL)
	respondJSON(w, http.StatusOK, SuccessResponse{Message: "Clinic address gallery image deleted successfully"})
}
