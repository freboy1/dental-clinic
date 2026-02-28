package dto


type CreateRequest struct {
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Street    string    `json:"street"`
	Building  string    `json:"building"`
	Latitude  float64   `json:"latitude"`
	Longitude float64  `json:"longitude"`
}

type CreateResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Address_id  string `json:"address_id"`
}

type AddressResponse struct {
	ID        string `json:"id"` 
	Country   string `json:"country"` 
	City      string `json:"city"` 
	Street    string `json:"street"` 
	Building  string `json:"building"` 
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"` 
}

type Response struct {
	Success string `json:"success"`
	Message string `json:"message"`
}