package models

import (
	"github.com/google/uuid"
)

type ClinicAddress struct {
	Id        uuid.UUID `json:"id"`
	ClinicId  uuid.UUID `json:"clinic_id"`
	AddressId uuid.UUID `json:"address_id"`
	IsMain    bool      `json:"is_main"`
}