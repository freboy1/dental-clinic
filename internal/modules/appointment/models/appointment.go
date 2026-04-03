package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	Id                uuid.UUID
	Doctor_id         uuid.UUID
	Clinic_address_id uuid.UUID
	Service_id        uuid.UUID
	User_id           uuid.UUID

	Start_time time.Time
	End_time   time.Time

	Status     string
	Created_at time.Time
	Name       string
	Email      string
}
