package dto

type UpdateMedicalRecordRequest struct {
	Diagnosis  string `json:"diagnosis"`
	Notes      string `json:"notes"`
	Is_checked bool   `json:"is_checked"`
}

type MedicalRecordResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}
