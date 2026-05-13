package models

import (
	"time"

	"github.com/google/uuid"
)

type MedicalRecord struct {
	Id             uuid.UUID
	Appointment_id uuid.UUID
	Doctor_id      uuid.UUID
	Patient_id     uuid.UUID
	Diagnosis      string
	Notes          string
	Is_checked     bool

	Created_at time.Time
	Updated_at time.Time
}
