package reports

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/reports/handlers"
	"dental_clinic/internal/modules/reports/repository"
	"dental_clinic/internal/modules/reports/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	_ = cfg
	repo := repository.NewReportsRepository(db)
	service := services.NewReportsService(repo)
	handler := handlers.NewReportsHandler(service)

	r.HandleFunc("/clinics/{clinicId}/reports/revenue", handler.GetRevenueReport).Methods("GET")
	r.HandleFunc("/clinics/{clinicId}/reports/appointments", handler.GetAppointmentReport).Methods("GET")
	r.HandleFunc("/clinics/{clinicId}/reports/doctors", handler.GetDoctorPerformanceReport).Methods("GET")
	r.HandleFunc("/clinics/{clinicId}/reports/inventory", handler.GetInventoryReport).Methods("GET")
}
