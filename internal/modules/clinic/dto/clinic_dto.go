package dto

type CreateClinicRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	IsActive    bool   `json:"is_active"`
}

type AddAddressRequest struct {
	Address_id string `json:"address_id"`
	Is_main    bool   `json:"is_main"`
}

type ClinicResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

type GetClinicAddressResponse struct {
	Id               string                       `json:"id"`
	Address_id       string                       `json:"address_id"`
	Address_name     string                       `json:"address_name"`
	Address_building string                       `json:"address_building"`
	CoverImageURL    string                       `json:"cover_image_url"`
	Gallery          []ClinicAddressImageResponse `json:"gallery"`
	Is_main          bool                         `json:"is_main"`
}

type ClinicAddressImageResponse struct {
	Id       string `json:"id"`
	ImageURL string `json:"image_url"`
}
