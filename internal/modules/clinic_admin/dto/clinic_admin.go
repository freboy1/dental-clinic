package dto

type CreateClinicAdminRequest struct {
	ClinicID string `json:"clinic_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}

type UpdateClinicAdminRequest struct {
	ClinicID    string `json:"clinic_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
	IsActive    bool   `json:"is_active"`
}

type ClinicAdminResponse struct {
	Id        string `json:"id"`
	ClinicID  string `json:"clinic_id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type ClinicAdminActionResponse struct {
	Success       string `json:"success"`
	Message       string `json:"message"`
	ClinicAdminID string `json:"clinic_admin_id,omitempty"`
}
