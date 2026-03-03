package dto

type CreateDoctorRequest struct {
	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
	ClinicID       string `json:"clinic_id"`
	Bio            string `json:"bio"`
	IsAvailable    bool   `json:"is_available"`
}

type UpdateDoctorRequest struct {
	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
	ClinicID       string `json:"clinic_id"`
	Bio            string `json:"bio"`
	IsAvailable    bool   `json:"is_available"`
}

type DoctorResponse struct {
	Id             string `json:"id"`
	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
	ClinicID       string `json:"clinic_id"`
	Bio            string `json:"bio"`
	IsAvailable    bool   `json:"is_available"`
}

type DoctorActionResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	DoctorID string `json:"doctor_id"`
}