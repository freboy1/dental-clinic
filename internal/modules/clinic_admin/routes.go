package clinic_admin

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/clinic_admin/handlers"
	"dental_clinic/internal/modules/clinic_admin/repository"
	"dental_clinic/internal/modules/clinic_admin/services"
	userRepository "dental_clinic/internal/modules/user/repository"
	userServices "dental_clinic/internal/modules/user/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewClinicAdminRepository(db)

	userRepo := userRepository.NewUserRepository(db)
	userService := userServices.NewUserService(userRepo, *cfg)

	service := services.NewClinicAdminService(repo, *userService)
	handler := handlers.NewClinicAdminHandler(service)

	r.HandleFunc("/clinic-admins", handler.CreateClinicAdmin).Methods("POST")
	r.HandleFunc("/clinic-admins", handler.GetClinicAdmins).Methods("GET")
	r.HandleFunc("/clinic-admins/{id}", handler.GetClinicAdminByID).Methods("GET")
	r.HandleFunc("/clinic-admins/{id}", handler.UpdateClinicAdmin).Methods("PUT")
	r.HandleFunc("/clinic-admins/{id}", handler.DeleteClinicAdmin).Methods("DELETE")
}
