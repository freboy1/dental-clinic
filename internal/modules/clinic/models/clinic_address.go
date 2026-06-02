package models

import (
	"github.com/google/uuid"
)

type ClinicAddress struct {
	Id            uuid.UUID `json:"id"`
	ClinicId      uuid.UUID `json:"clinic_id"`
	AddressId     uuid.UUID `json:"address_id"`
	IsMain        bool      `json:"is_main"`
	CoverImageURL string    `json:"cover_image_url"`
}

type ClinicAddressGalleryImage struct {
	Id              uuid.UUID `json:"id"`
	ClinicAddressId uuid.UUID `json:"clinic_address_id"`
	ImageURL        string    `json:"image_url"`
}

type ClinicAddressWithNames struct {
	Id              uuid.UUID                   `json:"id"`
	ClinicId        uuid.UUID                   `json:"clinic_id"`
	AddressId       uuid.UUID                   `json:"address_id"`
	AddressName     string                      `json:"address_name"`
	AddressBuilding string                      `json:"address_building"`
	IsMain          bool                        `json:"is_main"`
	CoverImageURL   string                      `json:"cover_image_url"`
	Gallery         []ClinicAddressGalleryImage `json:"gallery"`
}
