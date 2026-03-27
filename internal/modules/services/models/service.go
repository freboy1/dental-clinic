package models

import "github.com/google/uuid"

type Service struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       float64
	Duration    int
	ClinicID    uuid.UUID
	IsActive    bool
}

type ServiceWithClinicName struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       float64
	Duration    int
	ClinicID    uuid.UUID
	IsActive    bool
	ClinicName    string
}