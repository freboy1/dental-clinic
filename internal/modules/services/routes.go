package services

import (
	"dental_clinic/internal/modules/services/handlers"
	"dental_clinic/internal/modules/services/repository"
	"dental_clinic/internal/modules/services/services"

	"dental_clinic/internal/config"

	clinicRepository "dental_clinic/internal/modules/clinic/repository"
	clinicServices "dental_clinic/internal/modules/clinic/services"

	addressRepository "dental_clinic/internal/modules/address/repository"
	addressServices "dental_clinic/internal/modules/address/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewServiceRepository(db)

	addressRepo := addressRepository.NewAddressRepository(db)
	addressService := addressServices.NewAddressService(addressRepo, *cfg)

	clinicRepo := clinicRepository.NewClinicRepository(db)
	clinicService := clinicServices.NewClinicService(clinicRepo, *cfg, *addressService)

	service := services.NewServiceService(repo, *clinicService)

	handler := handlers.NewServiceHandler(service)

	r.HandleFunc("/services", handler.GetAllServices).Methods("GET")
	r.HandleFunc("/services/{id}", handler.GetServiceByID).Methods("GET")
	r.HandleFunc("/clinics/{clinic_id}/services", handler.GetServicesByClinic).Methods("GET")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewServiceRepository(db)

	addressRepo := addressRepository.NewAddressRepository(db)
	addressService := addressServices.NewAddressService(addressRepo, *cfg)

	clinicRepo := clinicRepository.NewClinicRepository(db)
	clinicService := clinicServices.NewClinicService(clinicRepo, *cfg, *addressService)

	service := services.NewServiceService(repo, *clinicService)
	handler := handlers.NewServiceHandler(service)

	r.HandleFunc("/services", handler.CreateService).Methods("POST")
	r.HandleFunc("/services/{id}", handler.UpdateService).Methods("PUT")
	r.HandleFunc("/services/{id}", handler.DeleteService).Methods("DELETE")
}
