package dto

type ProductRequest struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type ProductResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Unit      string `json:"unit"`
	CreatedAt string `json:"created_at"`
}

type InventoryQuantityRequest struct {
	ProductId string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

type UpdateInventoryRequest struct {
	Quantity float64 `json:"quantity"`
}

type InventoryResponse struct {
	Id              string  `json:"id"`
	ClinicAddressId string  `json:"clinic_address_id"`
	ProductId       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductUnit     string  `json:"product_unit"`
	Quantity        float64 `json:"quantity"`
	UpdatedAt       string  `json:"updated_at"`
}

type AttachMaterialRequest struct {
	ProductId        string  `json:"product_id"`
	QuantityRequired float64 `json:"quantity_required"`
}

type ServiceMaterialResponse struct {
	Id               string  `json:"id"`
	ClinicServiceId  string  `json:"clinic_service_id"`
	ProductId        string  `json:"product_id"`
	ProductName      string  `json:"product_name"`
	ProductUnit      string  `json:"product_unit"`
	QuantityRequired float64 `json:"quantity_required"`
}

type InventoryTransactionResponse struct {
	Id              string  `json:"id"`
	ClinicAddressId string  `json:"clinic_address_id"`
	ProductId       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	Quantity        float64 `json:"quantity"`
	TransactionType string  `json:"transaction_type"`
	AppointmentId   string  `json:"appointment_id,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

type ActionResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}
