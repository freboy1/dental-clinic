package ai_assistant

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/ai_assistant/handlers"
	aiRepository "dental_clinic/internal/modules/ai_assistant/repository"
	"dental_clinic/internal/modules/ai_assistant/services"

	addressRepository "dental_clinic/internal/modules/address/repository"
	addressServices "dental_clinic/internal/modules/address/services"
	appointmentRepository "dental_clinic/internal/modules/appointment/repository"
	appointmentServices "dental_clinic/internal/modules/appointment/services"
	clinicRepository "dental_clinic/internal/modules/clinic/repository"
	clinicServices "dental_clinic/internal/modules/clinic/services"
	medicalRecordRepository "dental_clinic/internal/modules/medical_record/repository"
	medicalRecordServices "dental_clinic/internal/modules/medical_record/services"
	reviewRepository "dental_clinic/internal/modules/reviews/repository"
	reviewServices "dental_clinic/internal/modules/reviews/services"
	scheduleRepository "dental_clinic/internal/modules/schedule/repository"
	scheduleServices "dental_clinic/internal/modules/schedule/services"
	serviceRepository "dental_clinic/internal/modules/services/repository"
	serviceServices "dental_clinic/internal/modules/services/services"
	userRepository "dental_clinic/internal/modules/user/repository"
	userServices "dental_clinic/internal/modules/user/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	addressRepo := addressRepository.NewAddressRepository(db)
	addressService := addressServices.NewAddressService(addressRepo, *cfg)

	clinicRepo := clinicRepository.NewClinicRepository(db)
	clinicService := clinicServices.NewClinicService(clinicRepo, *cfg, *addressService)

	serviceRepo := serviceRepository.NewServiceRepository(db)
	serviceService := serviceServices.NewServiceService(serviceRepo, *clinicService)

	scheduleRepo := scheduleRepository.NewScheduleRepository(db)
	scheduleService := scheduleServices.NewScheduleService(scheduleRepo, *cfg, *serviceService, *clinicService)

	medicalRecordRepo := medicalRecordRepository.NewMedicalRecordRepository(db)
	medicalRecordService := medicalRecordServices.NewMedicalRecordService(medicalRecordRepo)

	reviewRepo := reviewRepository.NewReviewRepository(db)
	reviewService := reviewServices.NewReviewService(reviewRepo)

	appointmentRepo := appointmentRepository.NewAppointmentRepository(db)
	appointmentService := appointmentServices.NewAppointmentService(appointmentRepo, db, *cfg, *scheduleService, *serviceService, *medicalRecordService, *clinicService, reviewService)

	userRepo := userRepository.NewUserRepository(db)
	userService := userServices.NewUserService(userRepo, *cfg)

	assistantRepo := aiRepository.NewAIAssistantRepository(db)
	llmClient := services.NewOpenAIClient(*cfg)
	assistantService := services.NewAIAssistantService(*cfg, assistantRepo, llmClient, *appointmentService, *scheduleService, *userService)
	handler := handlers.NewAIAssistantHandler(assistantService, *cfg)

	r.HandleFunc("/ai/chat", handler.Chat).Methods("POST")
	r.HandleFunc("/ai/chat/reset", handler.Reset).Methods("POST")
}
