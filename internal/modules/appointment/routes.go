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

	clinicRepository "dental_clinic/internal/modules/clinic/repository"
	clinicServices "dental_clinic/internal/modules/clinic/services"

	addressRepository "dental_clinic/internal/modules/address/repository"
	addressServices "dental_clinic/internal/modules/address/services"



	"github.com/gorilla/mux"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewAppointmentRepository(db)



	addressRepo := addressRepository.NewAddressRepository(db)
	addressService := addressServices.NewAddressService(addressRepo, *cfg)

	clinicRepo := clinicRepository.NewClinicRepository(db)
	clinicService := clinicServices.NewClinicService(clinicRepo, *cfg, *addressService)


	serviceRepo := serviceRepository.NewServiceRepository(db)
	serviceService := serviceServices.NewServiceService(serviceRepo, *clinicService)


	scheduleRepo := scheduleRepository.NewScheduleRepository(db)
	scheduleService := scheduleServices.NewScheduleService(scheduleRepo, *cfg, *serviceService)



	service := services.NewAppointmentService(repo, *cfg, *scheduleService, *serviceService)
	handler := handlers.NewAppointmentHandler(service, *cfg)

	r.HandleFunc("/appointment", handler.CreateAppointment).Methods("POST")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewAppointmentRepository(db)



	addressRepo := addressRepository.NewAddressRepository(db)
	addressService := addressServices.NewAddressService(addressRepo, *cfg)

	clinicRepo := clinicRepository.NewClinicRepository(db)
	clinicService := clinicServices.NewClinicService(clinicRepo, *cfg, *addressService)



	serviceRepo := serviceRepository.NewServiceRepository(db)
	serviceService := serviceServices.NewServiceService(serviceRepo, *clinicService)


	scheduleRepo := scheduleRepository.NewScheduleRepository(db)
	scheduleService := scheduleServices.NewScheduleService(scheduleRepo, *cfg, *serviceService)



	service := services.NewAppointmentService(repo, *cfg, *scheduleService, *serviceService)

	handler := handlers.NewAppointmentHandler(service, *cfg)

	r.HandleFunc("/appointment", handler.GetAllAppointments).Methods("GET")
	r.HandleFunc("/appointment/{id}", handler.GetAppointmentByID).Methods("GET")
	r.HandleFunc("/appointment/{id}", handler.UpdateAppointment).Methods("PUT")
	r.HandleFunc("/appointment/{id}", handler.DeleteAppointment).Methods("DELETE")
}