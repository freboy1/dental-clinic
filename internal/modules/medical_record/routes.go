package medical_record

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/medical_record/handlers"
	"dental_clinic/internal/modules/medical_record/repository"
	"dental_clinic/internal/modules/medical_record/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool) {
	//repo := repository.NewDoctorRepository(db)
	//service := services.NewDoctorService(repo)
	//handler := handlers.NewDoctorHandler(service)
	//
	//r.HandleFunc("/doctors", handler.GetAllDoctors).Methods("GET")
	//r.HandleFunc("/doctors/{id}", handler.GetDoctorByID).Methods("GET")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := repository.NewMedicalRecordRepository(db)
	service := services.NewMedicalRecordService(repo)
	handler := handlers.NewMedicalRecordHandler(service)

	r.HandleFunc("/medical-records/{id}", handler.GetMedicalRecord).Methods("GET")
}

func RegisterDoctorRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewMedicalRecordRepository(db)
	service := services.NewMedicalRecordService(repo)
	handler := handlers.NewMedicalRecordHandler(service)

	//r.HandleFunc("/doctors", handler.CreateDoctor).Methods("POST")
	r.HandleFunc("/medical-records/{id}", handler.UpdateMedicalRecord).Methods("PUT")
	//r.HandleFunc("/doctors/{id}", handler.DeleteDoctor).Methods("DELETE")
}
