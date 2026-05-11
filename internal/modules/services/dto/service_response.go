package dto

type CreateServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ServiceResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ServiceResponseWithName struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`
	ClinicID    string  `json:"clinic_id"`
	ClinicName  string  `json:"clinic_name"`
	IsActive    bool    `json:"is_active"`
}

type ServiceActionResponse struct {
	Success   string `json:"success"`
	Message   string `json:"message"`
	ServiceID string `json:"service_id"`
}
