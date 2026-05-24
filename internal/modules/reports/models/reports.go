package models

type ReportFilters struct {
	ClinicID        string
	ClinicAddressID string
	From            string
	To              string
}

type RevenueReportRow struct {
	ServiceID        string  `json:"service_id"`
	ServiceName      string  `json:"service_name"`
	AppointmentCount int     `json:"appointment_count"`
	UnitPrice        float64 `json:"unit_price"`
	TotalRevenue     float64 `json:"total_revenue"`
}

type AppointmentReportRow struct {
	Status           string `json:"status"`
	AppointmentCount int    `json:"appointment_count"`
}

type DoctorPerformanceRow struct {
	DoctorID         string  `json:"doctor_id"`
	DoctorName       string  `json:"doctor_name"`
	Specialization   string  `json:"specialization"`
	AppointmentCount int     `json:"appointment_count"`
	CompletedCount   int     `json:"completed_count"`
	Revenue          float64 `json:"revenue"`
	AverageRating    float64 `json:"average_rating"`
}

type InventoryReportRow struct {
	ProductID          string  `json:"product_id"`
	ProductName        string  `json:"product_name"`
	Unit               string  `json:"unit"`
	CurrentQuantity    float64 `json:"current_quantity"`
	RestockedQuantity  float64 `json:"restocked_quantity"`
	UsedQuantity       float64 `json:"used_quantity"`
	AdjustmentQuantity float64 `json:"adjustment_quantity"`
}
