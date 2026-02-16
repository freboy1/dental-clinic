package address

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/address/repository"
	"dental_clinic/internal/modules/address/services"
	"dental_clinic/internal/modules/address/handlers"
	"github.com/gorilla/mux"
)

// func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
// 	repo := repository.NewAddressRepository(db)
// 	service := services.NewAddressService(repo, *cfg)
// 	handler := handlers.NewAddressHandler(service, *cfg)

	// r.HandleFunc("/address", handler.Register).Methods("POST")
// }

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewAddressRepository(db)
	service := services.NewAddressService(repo, *cfg)
	handler := handlers.NewAddressHandler(service, *cfg)

	r.HandleFunc("/address", handler.CreateAddress).Methods("POST")
	r.HandleFunc("/address", handler.GetAllAddresss).Methods("GET")
	r.HandleFunc("/address/{id}", handler.GetAddressByID).Methods("GET")
	r.HandleFunc("/address/{id}", handler.DeleteAddress).Methods("DELETE")
}
