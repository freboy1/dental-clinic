package clinic

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/clinic/handlers"
	"dental_clinic/internal/modules/clinic/repository"
	"dental_clinic/internal/modules/clinic/services"

	addressRepository "dental_clinic/internal/modules/address/repository"
	addressServices "dental_clinic/internal/modules/address/services"


	"github.com/gorilla/mux"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewClinicRepository(db)
	
	addressRepo := addressRepository.NewAddressRepository(db)
	addressService := addressServices.NewAddressService(addressRepo, *cfg)

	service := services.NewClinicService(repo, *cfg, *addressService)
	handler := handlers.NewClinicHandler(service, *cfg)

	r.HandleFunc("/clinics", handler.GetClinics).Methods("GET")
	r.HandleFunc("/clinics/{id}", handler.GetClinic).Methods("GET")

	r.HandleFunc("/clinics/{id}/address", handler.GetClinicAddress).Methods("GET")

}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewClinicRepository(db)

	addressRepo := addressRepository.NewAddressRepository(db)
	addressService := addressServices.NewAddressService(addressRepo, *cfg)


	service := services.NewClinicService(repo, *cfg, *addressService)
	handler := handlers.NewClinicHandler(service, *cfg)

	r.HandleFunc("/clinics", handler.CreateClinic).Methods("POST")
	r.HandleFunc("/clinics/{id}", handler.UpdateClinic).Methods("PUT")
	r.HandleFunc("/clinics/{id}", handler.DeleteClinic).Methods("DELETE")

	r.HandleFunc("/clinics/{id}/address", handler.AddAddress).Methods("POST")
	r.HandleFunc("/clinics/{id}/address/{addressId}", handler.DeleteAddress).Methods("DELETE")
}
