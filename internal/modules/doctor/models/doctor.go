package models

import "github.com/google/uuid"

type Doctor struct {
	Id             uuid.UUID
	Specialization string
	Experience     int
	ClinicID       uuid.UUID
	Bio            string
	IsAvailable    bool
}