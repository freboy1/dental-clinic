package schedule

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"dental_clinic/internal/modules/schedule/repository"
	"dental_clinic/internal/modules/schedule/services"
	"dental_clinic/internal/modules/schedule/handlers"

	"dental_clinic/internal/config"

	"github.com/gorilla/mux"
)

func RegisterPublicRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	// repo := repository.NewScheduleRepository(db)
	// service := services.NewScheduleService(repo, *cfg)
	// handler := handlers.NewScheduleHandler(service, *cfg)

	// scheduleRouter := r.PathPrefix("/schedule").Subrouter()

	// scheduleRouter.HandleFunc("/available-slots/", handler.GetAvailableSlots).Methods("POST")
}

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	repo := repository.NewScheduleRepository(db)
	service := services.NewScheduleService(repo, *cfg)
	handler := handlers.NewScheduleHandler(service, *cfg)

	scheduleRouter := r.PathPrefix("/schedule").Subrouter()

	scheduleRouter.HandleFunc("/doctors/{doctorId}/working-hours", handler.CreateDoctorSchedule).Methods("POST")
	scheduleRouter.HandleFunc("/generate", handler.GenerateSlots).Methods("POST")
	// scheduleRouter.HandleFunc("/doctors/{doctorId}/working-hours", handler.GetDoctorSchedule).Methods("GET")

	// scheduleRouter.HandleFunc("/working-hours/{id}", handler.UpdateDoctorSchedule).Methods("PUT")
	// scheduleRouter.HandleFunc("/working-hours/{id}", handler.DeleteDoctorSchedule).Methods("DELETE")

	// scheduleRouter.HandleFunc("/doctors/{id}/slots", handler.GetSlots).Methods("GET")
}
