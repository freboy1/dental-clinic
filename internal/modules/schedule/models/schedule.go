package models

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	Id                uuid.UUID
	Doctor_id         uuid.UUID
	Clinic_address_id uuid.UUID
	Day_of_week       int
	Start_time        string
	End_time          string
}

type Slot struct {
	Id                uuid.UUID
	Doctor_id         uuid.UUID
	Clinic_address_id uuid.UUID
	Slot_start        time.Time
	Slot_end          time.Time
	Status            string
	Created_at        time.Time
}
