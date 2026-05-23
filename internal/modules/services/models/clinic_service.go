package models

import "github.com/google/uuid"

type Clinic_Service struct {
	Id uuid.UUID

	ClinicID  uuid.UUID
	ServiceID uuid.UUID

	Price    float64
	Duration int
	IsActive bool
}
