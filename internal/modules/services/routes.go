package services

import (
	"dental_clinic/internal/modules/services/handlers"
	"dental_clinic/internal/modules/services/repository"
	"dental_clinic/internal/modules/services/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := repository.NewServiceRepository(db)
	service := services.NewServiceService(repo)
	handler := handlers.NewServiceHandler(service)

	r.HandleFunc("/services", handler.GetAllServices).Methods("GET")
	r.HandleFunc("/services/{id}", handler.GetServiceByID).Methods("GET")
	r.HandleFunc("/clinics/{clinic_id}/services", handler.GetServicesByClinic).Methods("GET")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := repository.NewServiceRepository(db)
	service := services.NewServiceService(repo)
	handler := handlers.NewServiceHandler(service)

	r.HandleFunc("/services", handler.CreateService).Methods("POST")
	r.HandleFunc("/services/{id}", handler.UpdateService).Methods("PUT")
	r.HandleFunc("/services/{id}", handler.DeleteService).Methods("DELETE")
}