package services

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/clinic/repository"
)

type ClinicService struct {
	repo repository.ClinicRepository
	cfx  config.Config
}

func NewClinicService(r repository.ClinicRepository, cfx config.Config) *ClinicService {
	return &ClinicService{
		repo: r,
		cfx:  cfx,
	}
}