package doctor

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/doctor/handlers"
	"dental_clinic/internal/modules/doctor/repository"
	"dental_clinic/internal/modules/doctor/services"

	userRepository "dental_clinic/internal/modules/user/repository"
	userServices "dental_clinic/internal/modules/user/services"

	medical_recordRepository "dental_clinic/internal/modules/medical_record/repository"
	medical_recordServices "dental_clinic/internal/modules/medical_record/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewDoctorRepository(db)

	userRepo := userRepository.NewUserRepository(db)
	userService := userServices.NewUserService(userRepo, *cfg)

	medical_recordRepo := medical_recordRepository.NewMedicalRecordRepository(db)
	medical_recordService := medical_recordServices.NewMedicalRecordService(medical_recordRepo)

	service := services.NewDoctorService(repo, *userService, *medical_recordService)
	handler := handlers.NewDoctorHandler(service)

	r.HandleFunc("/doctors", handler.GetAllDoctors).Methods("GET")
	r.HandleFunc("/doctors/{id}", handler.GetDoctorByID).Methods("GET")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewDoctorRepository(db)

	userRepo := userRepository.NewUserRepository(db)
	userService := userServices.NewUserService(userRepo, *cfg)

	medical_recordRepo := medical_recordRepository.NewMedicalRecordRepository(db)
	medical_recordService := medical_recordServices.NewMedicalRecordService(medical_recordRepo)

	service := services.NewDoctorService(repo, *userService, *medical_recordService)
	handler := handlers.NewDoctorHandler(service)

	r.HandleFunc("/doctors", handler.CreateDoctor).Methods("POST")
	r.HandleFunc("/doctors/medical-records/{id}", handler.GetDoctorByIdMedicalRecords).Methods("GET")
	r.HandleFunc("/doctors/{id}", handler.UpdateDoctor).Methods("PUT")
	r.HandleFunc("/doctors/{id}", handler.DeleteDoctor).Methods("DELETE")
}
