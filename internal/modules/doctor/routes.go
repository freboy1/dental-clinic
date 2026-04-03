package doctor

import (
	"dental_clinic/internal/modules/doctor/handlers"
	"dental_clinic/internal/modules/doctor/repository"
	"dental_clinic/internal/modules/doctor/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := repository.NewDoctorRepository(db)
	service := services.NewDoctorService(repo)
	handler := handlers.NewDoctorHandler(service)

	r.HandleFunc("/doctors", handler.GetAllDoctors).Methods("GET")
	r.HandleFunc("/doctors/{id}", handler.GetDoctorByID).Methods("GET")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := repository.NewDoctorRepository(db)
	service := services.NewDoctorService(repo)
	handler := handlers.NewDoctorHandler(service)

	r.HandleFunc("/doctors", handler.CreateDoctor).Methods("POST")
	r.HandleFunc("/doctors/{id}", handler.UpdateDoctor).Methods("PUT")
	r.HandleFunc("/doctors/{id}", handler.DeleteDoctor).Methods("DELETE")
}
