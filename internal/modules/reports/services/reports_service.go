package services

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"dental_clinic/internal/modules/reports/models"
	"dental_clinic/internal/modules/reports/repository"

	"github.com/google/uuid"
)

type ReportsService struct {
	repo repository.ReportsRepository
}

func NewReportsService(repo repository.ReportsRepository) *ReportsService {
	return &ReportsService{repo: repo}
}

func (s *ReportsService) BuildFilters(clinicID, clinicAddressID, from, to string) (models.ReportFilters, error) {
	if _, err := uuid.Parse(clinicID); err != nil {
		return models.ReportFilters{}, errors.New("invalid clinic id")
	}
	if clinicAddressID != "" {
		if _, err := uuid.Parse(clinicAddressID); err != nil {
			return models.ReportFilters{}, errors.New("invalid clinic_address_id")
		}
	}
	if from == "" {
		return models.ReportFilters{}, errors.New("from date is required")
	}
	if to == "" {
		return models.ReportFilters{}, errors.New("to date is required")
	}
	if from != "" {
		if _, err := time.Parse("2006-01-02", from); err != nil {
			return models.ReportFilters{}, errors.New("invalid from date")
		}
	}
	if to != "" {
		if _, err := time.Parse("2006-01-02", to); err != nil {
			return models.ReportFilters{}, errors.New("invalid to date")
		}
	}
	return models.ReportFilters{
		ClinicID:        clinicID,
		ClinicAddressID: clinicAddressID,
		From:            from,
		To:              to,
	}, nil
}

func (s *ReportsService) Revenue(filters models.ReportFilters) ([]models.RevenueReportRow, error) {
	return s.repo.GetRevenueReport(filters)
}

func (s *ReportsService) Appointments(filters models.ReportFilters) ([]models.AppointmentReportRow, error) {
	return s.repo.GetAppointmentReport(filters)
}

func (s *ReportsService) Doctors(filters models.ReportFilters) ([]models.DoctorPerformanceRow, error) {
	return s.repo.GetDoctorPerformanceReport(filters)
}

func (s *ReportsService) Inventory(filters models.ReportFilters) ([]models.InventoryReportRow, error) {
	return s.repo.GetInventoryReport(filters)
}

func ToCSV(data interface{}) ([]byte, error) {
	rows, err := tableRows(data)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.WriteAll(rows); err != nil {
		return nil, err
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func ToPDF(title string, data interface{}) ([]byte, error) {
	rows, err := tableRows(data)
	if err != nil {
		return nil, err
	}

	lines := []string{title}
	for _, row := range rows {
		lines = append(lines, strings.Join(row, " | "))
	}

	var content strings.Builder
	content.WriteString("BT\n/F1 9 Tf\n50 780 Td\n")
	for i, line := range lines {
		if i > 0 {
			content.WriteString("0 -14 Td\n")
		}
		content.WriteString("(")
		content.WriteString(escapePDF(line))
		content.WriteString(") Tj\n")
	}
	content.WriteString("ET")

	stream := content.String()
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(stream), stream),
	}

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	offsets := make([]int, 0, len(objects)+1)
	offsets = append(offsets, 0)
	for i, object := range objects {
		offsets = append(offsets, buf.Len())
		buf.WriteString(strconv.Itoa(i + 1))
		buf.WriteString(" 0 obj\n")
		buf.WriteString(object)
		buf.WriteString("\nendobj\n")
	}
	xref := buf.Len()
	buf.WriteString("xref\n0 ")
	buf.WriteString(strconv.Itoa(len(objects) + 1))
	buf.WriteString("\n0000000000 65535 f \n")
	for i := 1; i < len(offsets); i++ {
		buf.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}
	buf.WriteString("trailer\n<< /Size ")
	buf.WriteString(strconv.Itoa(len(objects) + 1))
	buf.WriteString(" /Root 1 0 R >>\nstartxref\n")
	buf.WriteString(strconv.Itoa(xref))
	buf.WriteString("\n%%EOF")
	return buf.Bytes(), nil
}

func tableRows(data interface{}) ([][]string, error) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		return nil, errors.New("report data must be a slice")
	}
	if value.Len() == 0 {
		return [][]string{{"empty"}}, nil
	}

	elemType := value.Type().Elem()
	headers := make([]string, 0, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		headers = append(headers, jsonName(elemType.Field(i)))
	}

	rows := [][]string{headers}
	for i := 0; i < value.Len(); i++ {
		row := make([]string, 0, elemType.NumField())
		elem := value.Index(i)
		for j := 0; j < elem.NumField(); j++ {
			row = append(row, fmt.Sprint(elem.Field(j).Interface()))
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func jsonName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}
	parts := strings.Split(tag, ",")
	if parts[0] == "" {
		return field.Name
	}
	return parts[0]
}

func escapePDF(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "(", "\\(")
	value = strings.ReplaceAll(value, ")", "\\)")
	return value
}
