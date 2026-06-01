package dto

type CreateDoctorRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
	ClinicID       string `json:"clinic_id"`
	Bio            string `json:"bio"`
	IsAvailable    bool   `json:"is_available"`
	Password       string `json:"password"`
	Is_active      bool   `json:"is_active"`
}

type UpdateDoctorRequest struct {
	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
	// ClinicID       string `json:"clinic_id"`
	Bio         string `json:"bio"`
	IsAvailable bool   `json:"is_available"`
	NewPassword string `json:"new_password"`
	Is_active   bool   `json:"is_active"`
}

type DoctorResponse struct {
	Id             string  `json:"id"`
	Specialization string  `json:"specialization"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Experience     int     `json:"experience"`
	ClinicID       string  `json:"clinic_id"`
	Bio            string  `json:"bio"`
	IsAvailable    bool    `json:"is_available"`
	Rating         float64 `json:"rating"`
	PhotoURL       string  `json:"photo_url"`
}

type DoctorActionResponse struct {
	Success          string `json:"success"`
	Message          string `json:"message"`
	DoctorID         string `json:"doctor_id"`
	ConfirmationCode string `json:"confirmation_code,omitempty"`
}

type GetMedicalRecordDoctorResponse struct {
	Id         string `json:"id"`
	Diagnosis  string `json:"diagnosis"`
	Notes      string `json:"notes"`
	Is_checked bool   `json:"is_checked"`
	Created_at string `json:"created_at"`
}

type DoctorPhotoRequest struct {
	PhotoURL string `json:"photo_url"`
}
