package dto

type CreateServiceRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`
	ClinicID    string  `json:"clinic_id"`
	IsActive    bool    `json:"is_active"`
}

type UpdateServiceRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`
	ClinicID    string  `json:"clinic_id"`
	IsActive    bool    `json:"is_active"`
}

type ServiceResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`
	ClinicID    string  `json:"clinic_id"`
	IsActive    bool    `json:"is_active"`
}

type ServiceActionResponse struct {
	Success   string `json:"success"`
	Message   string `json:"message"`
	ServiceID string `json:"service_id"`
}