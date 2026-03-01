package dto

type AddAddressRequest struct {
	Address_id   string    `json:"address_id"`
	Is_main      bool    `json:"is_main"`
}

type ClinicResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

type GetClinicAddressResponse struct {
	Address_id   string    `json:"address_id"`
	Is_main      bool    `json:"is_main"`
}