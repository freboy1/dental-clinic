package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id        uuid.UUID
	Name      string
	Unit      string
	CreatedAt time.Time
}

type AddressInventory struct {
	Id              uuid.UUID
	ClinicAddressId uuid.UUID
	ProductId       uuid.UUID
	ProductName     string
	ProductUnit     string
	Quantity        float64
	UpdatedAt       time.Time
}

type InventoryTransaction struct {
	Id              uuid.UUID
	ClinicAddressId uuid.UUID
	ProductId       uuid.UUID
	ProductName     string
	Quantity        float64
	TransactionType string
	AppointmentId   uuid.UUID
	CreatedAt       time.Time
}

type ServiceMaterial struct {
	Id               uuid.UUID
	ClinicServiceId  uuid.UUID
	ProductId        uuid.UUID
	ProductName      string
	ProductUnit      string
	QuantityRequired float64
}
