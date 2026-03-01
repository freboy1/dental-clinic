package models

import (
	"github.com/google/uuid"
)

type ClinicAddress struct {
	Id        uuid.UUID 
	ClinicId uuid.UUID 
	AddressId uuid.UUID 
	IsMain    bool      
}
