package dto

import "mime/multipart"

type UpdateMedicalRecordRequest struct {
	Diagnosis  string                  `form:"diagnosis"`
	Notes      string                  `form:"notes"`
	Is_checked bool                    `form:"is_checked"`
	Files      []*multipart.FileHeader `form:"files"`
}

type MedicalRecordResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

type GetMedicalRecordResponse struct {
	Status     string                `json:"status"`
	Message    string                `json:"message"`
	Diagnosis  string                `json:"diagnosis"`
	Notes      string                `json:"notes"`
	Is_checked bool                  `json:"is_checked"`
	Files      []MedicalFileResponse `json:"files"`
}

type MedicalFileResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
}

type UpdateMedicalRecordResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}
