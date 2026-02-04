package clinic

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/clinic/handlers"
	"dental_clinic/internal/modules/clinic/repository"
	"dental_clinic/internal/modules/clinic/services"
	"github.com/gorilla/mux"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewClinicRepository(db)
	service := services.NewClinicService(repo, *cfg)
	handler := handlers.NewClinicHandler(service, *cfg)

	r.HandleFunc("/clinics", handler.GetClinics).Methods("GET")
	r.HandleFunc("/clinics/{id}", handler.GetClinic).Methods("GET")

}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewClinicRepository(db)
	service := services.NewClinicService(repo, *cfg)
	handler := handlers.NewClinicHandler(service, *cfg)

	r.HandleFunc("/clinics", handler.CreateClinic).Methods("POST")
	r.HandleFunc("/clinics/{id}", handler.UpdateClinic).Methods("PUT")
	r.HandleFunc("/clinics/{id}", handler.DeleteClinic).Methods("DELETE")

}
