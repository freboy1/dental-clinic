package dto

import "mime/multipart"

type UpdateMedicalRecordRequest struct {
	Diagnosis  string                  `json:"diagnosis"`
	Notes      string                  `json:"notes"`
	Is_checked bool                    `json:"is_checked"`
	Files      []*multipart.FileHeader `form:"files"`
}

type MedicalRecordResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

type GetMedicalRecordResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Diagnosis  string `json:"diagnosis"`
	Notes      string `json:"notes"`
	Is_checked bool   `json:"is_checked"`
}
