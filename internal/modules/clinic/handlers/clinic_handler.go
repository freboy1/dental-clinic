package handlers

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/clinic/services"
	"net/http"
)

type ClinicHandler struct {
	service *services.ClinicService
	cfg config.Config
}

func NewClinicHandler(s *services.ClinicService, cfg config.Config) *ClinicHandler {
	return &ClinicHandler{
		service: s,
		cfg: cfg,
	}
}
func (h *ClinicHandler) GetClinics(w http.ResponseWriter, r *http.Request) {

}
func (h *ClinicHandler) GetClinic(w http.ResponseWriter, r *http.Request) {

}
func (h *ClinicHandler) CreateClinic(w http.ResponseWriter, r *http.Request) {

}
func (h *ClinicHandler) UpdateClinic(w http.ResponseWriter, r *http.Request) {

}
func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {

}