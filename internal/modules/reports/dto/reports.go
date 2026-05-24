package dto

type ReportResponse struct {
	ClinicID        string      `json:"clinic_id"`
	ClinicAddressID string      `json:"clinic_address_id,omitempty"`
	From            string      `json:"from,omitempty"`
	To              string      `json:"to,omitempty"`
	Data            interface{} `json:"data"`
}
