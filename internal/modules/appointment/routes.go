package appointment

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/appointment/repository"
	"dental_clinic/internal/modules/appointment/services"
	"dental_clinic/internal/modules/appointment/handlers"

	scheduleRepository "dental_clinic/internal/modules/schedule/repository"
	scheduleServices "dental_clinic/internal/modules/schedule/services"

	serviceRepository "dental_clinic/internal/modules/services/repository"
	serviceServices "dental_clinic/internal/modules/services/services"


	"github.com/gorilla/mux"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewAppointmentRepository(db)

	serviceRepo := serviceRepository.NewServiceRepository(db)
	serviceService := serviceServices.NewServiceService(serviceRepo)


	scheduleRepo := scheduleRepository.NewScheduleRepository(db)
	scheduleService := scheduleServices.NewScheduleService(scheduleRepo, *cfg, *serviceService)



	service := services.NewAppointmentService(repo, *cfg, *scheduleService, *serviceService)
	handler := handlers.NewAppointmentHandler(service, *cfg)

	r.HandleFunc("/appointment", handler.CreateAppointment).Methods("POST")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewAppointmentRepository(db)

	serviceRepo := serviceRepository.NewServiceRepository(db)
	serviceService := serviceServices.NewServiceService(serviceRepo)


	scheduleRepo := scheduleRepository.NewScheduleRepository(db)
	scheduleService := scheduleServices.NewScheduleService(scheduleRepo, *cfg, *serviceService)



	service := services.NewAppointmentService(repo, *cfg, *scheduleService, *serviceService)
	
	handler := handlers.NewAppointmentHandler(service, *cfg)

	r.HandleFunc("/appointment", handler.GetAllAppointments).Methods("GET")

}
