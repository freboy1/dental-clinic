package handlers

import (
	"database/sql"
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/address/services"
	"dental_clinic/internal/modules/address/dto"
	"encoding/json"
	"errors"
	"net/http"
	"github.com/gorilla/mux"
	"strings"
)

type AddressHandler struct {
	service *services.AddressService
	cfg config.Config
}

func NewAddressHandler(s *services.AddressService, cfg config.Config) *AddressHandler {
	return &AddressHandler{
		service: s,
		cfg: cfg,
	}
}
// CreateAddress godoc
// @Summary Create new address
// @Description Creates a new address
// @Tags Addresss
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body dto.CreateRequest true "Address registration data"
// @Success 200 {object} dto.CreateResponse
// @Failure 400 {object} dto.CreateResponse
// @Router /api/address [post]
func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	response := dto.CreateResponse{
		Success: "0",
		Message: "",
		Address_id:  "",
	}

	var req dto.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	address, err := h.service.CreateAddress(req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Success = "1"
	response.Message = "successfully created"
	response.Address_id = address.ID.String()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}


// GetAllAddresss godoc
// @Summary Get all addresss
// @Description Returns a list of all addresss
// @Tags Addresss
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.AddressResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/address [get]
func (h *AddressHandler) GetAllAddresss(w http.ResponseWriter, r *http.Request) {
	_ = r
	tokenStr := getToken(r)
	addresss, err := h.service.GetAllAddresss(tokenStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(services.ToAddressResponseList(addresss))
}

// GetAddress godoc
// @Summary Get address by ID
// @Description Returns a single address by its UUID
// @Tags Addresss
// @Security BearerAuth
// @Produce json
// @Param id path string true "Address ID (UUID)"
// @Success 200 {object} dto.AddressResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/address/{id} [get]
func (h *AddressHandler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	address, err := h.service.GetAddressByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.ToAddressResponse(*address))
}

// UpdateAddresss godoc
// @Summary Update addresss
// @Description Update & return status
// @Tags Addresss
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Address ID (UUID)"
// @Param request body dto.CreateRequest true "Address update data"
// @Success 200 {array} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/address/{id} [put]
func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	response := dto.Response{
		Success: "0",
		Message: "",
	}
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err := h.service.UpdateAddress(id, req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "successfully updated"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteAddress godoc
// @Summary Delete Address
// @Description Deletes an address by UUID
// @Tags Addresss
// @Security BearerAuth
// @Produce json
// @Param id path string true "Address ID (UUID)"
// @Success 200 {array} dto.Response
// @Failure 404 {array} dto.Response
// @Router /api/address/{id} [delete]
func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	response := dto.Response{
		Success: "0",
		Message: "",
	}
	vars := mux.Vars(r)
	id := vars["id"]
	tokenStr := getToken(r)
	err := h.service.DeleteAddress(id, tokenStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "address not found" {
			response.Message = "Address not found"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		response.Message = "Internal server error"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = "1"
	response.Message = "successfully deleted"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenStr := parts[1]
	return tokenStr
}