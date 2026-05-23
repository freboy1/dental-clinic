package dto

import (
	"dental_clinic/internal/modules/ai_assistant/models"
	"time"
)

type ChatRequest struct {
	Message    string `json:"message"`
	ChoiceType string `json:"choice_type"`
	ChoiceID   string `json:"choice_id"`
}

type ChatResponse struct {
	Reply          string                 `json:"reply"`
	SessionID      string                 `json:"session_id"`
	ChoiceRequired bool                   `json:"choice_required"`
	ChoiceType     string                 `json:"choice_type,omitempty"`
	State          models.BookingState    `json:"state"`
	Services       []models.ServiceOption `json:"services,omitempty"`
	Clinics        []models.ClinicOption  `json:"clinics,omitempty"`
	Doctors        []models.DoctorOption  `json:"doctors,omitempty"`
	AvailableSlots []SlotResponse         `json:"available_slots,omitempty"`
	AppointmentID  string                 `json:"appointment_id,omitempty"`
}

type SlotResponse struct {
	Id        string    `json:"id"`
	SlotStart time.Time `json:"slot_start"`
	SlotEnd   time.Time `json:"slot_end"`
	Status    string    `json:"status"`
}
