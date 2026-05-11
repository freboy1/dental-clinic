package models

import "github.com/google/uuid"

type ChatSession struct {
	Id     uuid.UUID
	UserID uuid.UUID
}

type ChatMessage struct {
	Id        uuid.UUID
	SessionID uuid.UUID
	Role      string
	Content   string
}

type BookingState struct {
	UserID          string `json:"user_id"`
	DoctorID        string `json:"doctor_id,omitempty"`
	ServiceID       string `json:"service_id,omitempty"`
	ClinicAddressID string `json:"clinic_address_id,omitempty"`
	Date            string `json:"date,omitempty"`
	Time            string `json:"time,omitempty"`
	Step            string `json:"step,omitempty"`
}

func (s BookingState) IsComplete() bool {
	return s.DoctorID != "" &&
		s.ServiceID != "" &&
		s.ClinicAddressID != "" &&
		s.Date != "" &&
		s.Time != ""
}

type ServiceOption struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ClinicOption struct {
	ClinicID        string  `json:"clinic_id"`
	ClinicAddressID string  `json:"clinic_address_id"`
	ClinicName      string  `json:"clinic_name"`
	Price           float64 `json:"price"`
	Duration        int     `json:"duration"`
}

type DoctorOption struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
}
