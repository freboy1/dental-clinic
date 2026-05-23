package models

import (
	"time"

	"github.com/google/uuid"
)

type DoctorRating struct {
	Id            uuid.UUID
	AppointmentId uuid.UUID
	DoctorId      uuid.UUID
	UserId        uuid.UUID
	Rating        int
	CreatedAt     time.Time
}

type ClinicReview struct {
	Id            uuid.UUID
	AppointmentId uuid.UUID
	ClinicId      uuid.UUID
	UserId        uuid.UUID
	Rating        int
	Comment       string
	CreatedAt     time.Time
}
