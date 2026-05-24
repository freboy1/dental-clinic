package handlers

import (
	"encoding/json"
	"net/http"

	"dental_clinic/internal/modules/reports/dto"
	"dental_clinic/internal/modules/reports/models"
	"dental_clinic/internal/modules/reports/services"

	"github.com/gorilla/mux"
)

type ReportsHandler struct {
	service *services.ReportsService
}

func NewReportsHandler(service *services.ReportsService) *ReportsHandler {
	return &ReportsHandler{service: service}
}

// GetRevenueReport godoc
// @Summary Get clinic revenue report
// @Description Returns clinic revenue report. Use format=csv or format=pdf to export.
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Param from query string true "Start date YYYY-MM-DD"
// @Param to query string true "End date YYYY-MM-DD"
// @Param clinic_address_id query string false "Clinic address ID"
// @Param format query string false "Export format: csv or pdf"
// @Success 200 {object} dto.ReportResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinics/{clinicId}/reports/revenue [get]
func (h *ReportsHandler) GetRevenueReport(w http.ResponseWriter, r *http.Request) {
	filters, ok := h.filters(w, r)
	if !ok {
		return
	}
	data, err := h.service.Revenue(filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondReport(w, r, "Revenue Report", filters.ClinicID, filters.ClinicAddressID, filters.From, filters.To, data)
}

// GetAppointmentReport godoc
// @Summary Get clinic appointment report
// @Description Returns appointment counts by status. Use format=csv or format=pdf to export.
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Param from query string true "Start date YYYY-MM-DD"
// @Param to query string true "End date YYYY-MM-DD"
// @Param clinic_address_id query string false "Clinic address ID"
// @Param format query string false "Export format: csv or pdf"
// @Success 200 {object} dto.ReportResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinics/{clinicId}/reports/appointments [get]
func (h *ReportsHandler) GetAppointmentReport(w http.ResponseWriter, r *http.Request) {
	filters, ok := h.filters(w, r)
	if !ok {
		return
	}
	data, err := h.service.Appointments(filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondReport(w, r, "Appointment Report", filters.ClinicID, filters.ClinicAddressID, filters.From, filters.To, data)
}

// GetDoctorPerformanceReport godoc
// @Summary Get clinic doctor performance report
// @Description Returns doctor appointment, revenue, and rating performance. Use format=csv or format=pdf to export.
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Param from query string true "Start date YYYY-MM-DD"
// @Param to query string true "End date YYYY-MM-DD"
// @Param clinic_address_id query string false "Clinic address ID"
// @Param format query string false "Export format: csv or pdf"
// @Success 200 {object} dto.ReportResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinics/{clinicId}/reports/doctors [get]
func (h *ReportsHandler) GetDoctorPerformanceReport(w http.ResponseWriter, r *http.Request) {
	filters, ok := h.filters(w, r)
	if !ok {
		return
	}
	data, err := h.service.Doctors(filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondReport(w, r, "Doctor Performance Report", filters.ClinicID, filters.ClinicAddressID, filters.From, filters.To, data)
}

// GetInventoryReport godoc
// @Summary Get clinic inventory report
// @Description Returns current stock and transaction quantities by product. Use format=csv or format=pdf to export.
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Param from query string true "Start date YYYY-MM-DD"
// @Param to query string true "End date YYYY-MM-DD"
// @Param clinic_address_id query string false "Clinic address ID"
// @Param format query string false "Export format: csv or pdf"
// @Success 200 {object} dto.ReportResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinics/{clinicId}/reports/inventory [get]
func (h *ReportsHandler) GetInventoryReport(w http.ResponseWriter, r *http.Request) {
	filters, ok := h.filters(w, r)
	if !ok {
		return
	}
	data, err := h.service.Inventory(filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondReport(w, r, "Inventory Report", filters.ClinicID, filters.ClinicAddressID, filters.From, filters.To, data)
}

func (h *ReportsHandler) filters(w http.ResponseWriter, r *http.Request) (models.ReportFilters, bool) {
	vars := mux.Vars(r)
	query := r.URL.Query()
	filters, err := h.service.BuildFilters(vars["clinicId"], query.Get("clinic_address_id"), query.Get("from"), query.Get("to"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return models.ReportFilters{}, false
	}
	return filters, true
}

func (h *ReportsHandler) respondReport(w http.ResponseWriter, r *http.Request, title, clinicID, clinicAddressID, from, to string, data interface{}) {
	switch r.URL.Query().Get("format") {
	case "csv":
		body, err := services.ToCSV(data)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="report.csv"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	case "pdf":
		body, err := services.ToPDF(title, data)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="report.pdf"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	default:
		respondJSON(w, http.StatusOK, dto.ReportResponse{
			ClinicID:        clinicID,
			ClinicAddressID: clinicAddressID,
			From:            from,
			To:              to,
			Data:            data,
		})
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
