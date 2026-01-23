package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/user/repository"
	"dental_clinic/internal/modules/user/services"
	"dental_clinic/internal/modules/user/handlers"
	"github.com/gorilla/mux"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewUserRepository(db)
	service := services.NewUserService(repo, *cfg)
	handler := handlers.NewUserHandler(service, *cfg)

	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/verify", handler.VerifyAccountByLink).Methods("GET")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewUserRepository(db)
	service := services.NewUserService(repo, *cfg)
	handler := handlers.NewUserHandler(service, *cfg)

	r.HandleFunc("/users/{id}", handler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users/update-password", handler.UpdatePassword).Methods("POST")
	r.HandleFunc("/users/update-email", handler.UpdateEmail).Methods("POST")
	r.HandleFunc("/users/verify-email", handler.VerifyNewEmail).Methods("GET")
	r.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
}
