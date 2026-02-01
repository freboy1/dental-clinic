package clinic

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"dental_clinic/internal/config"
	"github.com/gorilla/mux"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewUserRepository(db)
	service := services.NewUserService(repo, *cfg)
	handler := handlers.NewUserHandler(service, *cfg)

	// r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/clinics", handler.GetClinics).Methods("GET")
	r.HandleFunc("/clinics/{id}", handler.GetClinic).Methods("GET")

}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewUserRepository(db)
	service := services.NewUserService(repo, *cfg)
	handler := handlers.NewUserHandler(service, *cfg)

	r.HandleFunc("/clinics", handler.CreateClinic).Methods("POST")
	r.HandleFunc("/clinics/{id}", handler.UpdateClinic).Methods("PUT")
	r.HandleFunc("/clinics/{id}", handler.DeleteClinic).Methods("DELETE")

}
