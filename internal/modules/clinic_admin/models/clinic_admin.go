package models

import (
	"time"

	"github.com/google/uuid"
)

type ClinicAdmin struct {
	Id        uuid.UUID
	ClinicID  uuid.UUID
	UserID    uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
}
