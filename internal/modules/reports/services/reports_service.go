package services

import (
	"bytes"
	"dental_clinic/internal/modules/reports/models"
	"dental_clinic/internal/modules/reports/repository"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

// ── Service wiring ────────────────────────────────────────────────────────────

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
	if _, err := time.Parse("2006-01-02", from); err != nil {
		return models.ReportFilters{}, errors.New("invalid from date")
	}
	if _, err := time.Parse("2006-01-02", to); err != nil {
		return models.ReportFilters{}, errors.New("invalid to date")
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

// ── CSV (unchanged) ───────────────────────────────────────────────────────────

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

// ── PDF ───────────────────────────────────────────────────────────────────────

// ToPDF generates a branded, formatted PDF for any of the four report types.
// The from/to strings are displayed in the report header.
func ToPDF(title, from, to string, data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case []models.AppointmentReportRow:
		return buildAppointmentPDF(title, from, to, v)
	case []models.RevenueReportRow:
		return buildRevenuePDF(title, from, to, v)
	case []models.DoctorPerformanceRow:
		return buildDoctorPDF(title, from, to, v)
	case []models.InventoryReportRow:
		return buildInventoryPDF(title, from, to, v)
	default:
		return nil, fmt.Errorf("unsupported report data type: %T", data)
	}
}

// ── Shared PDF layout constants & helpers ─────────────────────────────────────

const (
	pageW  = 210.0
	pageH  = 297.0
	margin = 14.0
	inner  = pageW - 2*margin

	headerH    = 44.0
	titleBarH  = 14.0
	infoBoxH   = 32.0
	rowH       = 9.0
	tableHeadH = 9.0
)

type color struct{ r, g, b int }

var (
	cDark    = color{15, 23, 42}
	cPrimary = color{37, 99, 235}
	cAccent  = color{16, 185, 129}
	cMuted   = color{100, 116, 139}
	cSubtle  = color{248, 250, 252}
	cBorder  = color{226, 232, 240}
	cWhite   = color{255, 255, 255}

	// Status badge colours for appointment report
	statusPalette = map[string]color{
		"confirmed":   {37, 99, 235},
		"completed":   {16, 185, 129},
		"cancelled":   {239, 68, 68},
		"pending":     {245, 158, 11},
		"no_show":     {156, 163, 175},
		"rescheduled": {139, 92, 246},
	}

	// Donut chart colour sequence
	chartColors = []color{
		{37, 99, 235},
		{16, 185, 129},
		{245, 158, 11},
		{239, 68, 68},
		{139, 92, 246},
		{156, 163, 175},
		{20, 184, 166},
		{249, 115, 22},
	}
)

func fill(pdf *gofpdf.Fpdf, c color) { pdf.SetFillColor(c.r, c.g, c.b) }
func text(pdf *gofpdf.Fpdf, c color) { pdf.SetTextColor(c.r, c.g, c.b) }
func draw(pdf *gofpdf.Fpdf, c color) { pdf.SetDrawColor(c.r, c.g, c.b) }

// newPDF initialises a fresh A4 portrait PDF with no margins.
// Auto page break is disabled because every element is placed with absolute
// SetXY coordinates — letting gofpdf auto-break causes spurious blank pages.
func newPDF() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	return pdf
}

// renderHeader draws the dark clinic header, blue title bar, and info box.
// Returns the Y position immediately below the info box.
func renderHeader(pdf *gofpdf.Fpdf, title, from, to string) float64 {
	// Dark header band
	fill(pdf, cDark)
	pdf.Rect(0, 0, pageW, headerH, "F")

	// Tooth logo (three ellipses)
	drawToothLogo(pdf, margin, 11, 22)

	// Clinic name
	pdf.SetFont("Helvetica", "B", 14)
	text(pdf, cWhite)
	pdf.SetXY(margin+26, 13)
	pdf.Cell(100, 7, "DentalCare Clinic")

	// Tagline
	pdf.SetFont("Helvetica", "", 8)
	text(pdf, color{148, 163, 184})
	pdf.SetXY(margin+26, 22)
	pdf.Cell(100, 5, "Your smile is our priority")

	// Date range (top-right)
	rangeStr := from + "  to  " + to
	pdf.SetFont("Helvetica", "", 8)
	text(pdf, color{148, 163, 184})
	pdf.SetXY(pageW-margin-70, 14)
	pdf.CellFormat(70, 5, rangeStr, "", 0, "R", false, 0, "")
	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetXY(pageW-margin-70, 21)
	pdf.CellFormat(70, 5, "Generated "+time.Now().Format("02 Jan 2006  15:04 UTC"), "", 0, "R", false, 0, "")

	// Blue title bar
	fill(pdf, cPrimary)
	pdf.Rect(0, headerH, pageW, titleBarH, "F")
	pdf.SetFont("Helvetica", "B", 10)
	text(pdf, cWhite)
	pdf.SetXY(margin, headerH+3)
	pdf.Cell(inner, 8, strings.ToUpper(title))

	// Info box
	infoY := headerH + titleBarH + 6
	fill(pdf, cSubtle)
	draw(pdf, cBorder)
	pdf.SetLineWidth(0.3)
	pdf.RoundedRect(margin, infoY, inner, infoBoxH, 3, "1234", "FD")
	fill(pdf, cAccent)
	pdf.RoundedRect(margin, infoY, 3, infoBoxH, 1.5, "1234", "F")
	pdf.Rect(margin+1.5, infoY, 1.5, infoBoxH, "F")

	pdf.SetFont("Helvetica", "B", 8.5)
	text(pdf, cDark)
	pdf.SetXY(margin+7, infoY+4)
	pdf.Cell(100, 5, "About DentalCare Clinic")

	pdf.SetFont("Helvetica", "", 7.5)
	text(pdf, cMuted)
	pdf.SetXY(margin+7, infoY+11)
	pdf.MultiCell(inner-10, 4,
		"This report is generated as part of a diploma project — a dental clinic management system. "+
			"The system covers appointment scheduling, doctor performance tracking, revenue analytics, "+
			"and inventory management across clinic locations.",
		"", "L", false)

	return infoY + infoBoxH + 8
}

// renderFooter draws the page footer at the bottom of the current page.
func renderFooter(pdf *gofpdf.Fpdf) {
	y := pageH - 11
	draw(pdf, cBorder)
	pdf.SetLineWidth(0.2)
	pdf.Line(margin, y, margin+inner, y)
	pdf.SetFont("Helvetica", "", 7)
	text(pdf, cMuted)
	pdf.SetXY(margin, y+2)
	pdf.Cell(inner/2, 5, "DentalCare Clinic  \u2014  Confidential")
	pdf.SetXY(margin+inner/2, y+2)
	pdf.CellFormat(inner/2, 5, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "R", false, 0, "")
}

// sectionHeading draws a small bold section label + divider rule. Returns Y after rule.
func sectionHeading(pdf *gofpdf.Fpdf, y float64, label string) float64 {
	pdf.SetFont("Helvetica", "B", 9)
	text(pdf, cDark)
	pdf.SetXY(margin, y)
	pdf.Cell(inner, 6, label)
	draw(pdf, cBorder)
	pdf.SetLineWidth(0.25)
	pdf.Line(margin, y+7, margin+inner, y+7)
	return y + 9
}

// tableHeader draws the dark table header row and returns Y after it.
func tableHeader(pdf *gofpdf.Fpdf, y float64, cols []tableCol) float64 {
	fill(pdf, cDark)
	pdf.Rect(margin, y, inner, tableHeadH, "F")
	pdf.SetFont("Helvetica", "B", 7.5)
	text(pdf, cWhite)
	for _, c := range cols {
		pdf.SetXY(c.x, y+1.5)
		pdf.CellFormat(c.w, tableHeadH-3, strings.ToUpper(c.label), "", 0, c.align, false, 0, "")
	}
	return y + tableHeadH
}

type tableCol struct {
	x, w  float64
	label string
	align string // "L" | "R" | "C"
}

// tableRowBand fills the alternating background for a data row.
func tableRowBand(pdf *gofpdf.Fpdf, y float64, idx int) {
	if idx%2 == 0 {
		fill(pdf, cSubtle)
	} else {
		fill(pdf, cWhite)
	}
	pdf.Rect(margin, y, inner, rowH, "F")
	draw(pdf, cBorder)
	pdf.SetLineWidth(0.1)
	pdf.Line(margin, y+rowH, margin+inner, y+rowH)
}

// progressBar draws a labelled progress bar within a table cell.
func progressBar(pdf *gofpdf.Fpdf, x, y, maxW, pct float64, c color) {
	bH := 4.0
	bY := y + (rowH-bH)/2
	fill(pdf, cBorder)
	pdf.RoundedRect(x, bY, maxW, bH, 1.5, "1234", "F")
	if pct > 0 {
		fill(pdf, c)
		w := maxW * pct / 100
		if w > maxW {
			w = maxW
		}
		pdf.RoundedRect(x, bY, w, bH, 1.5, "1234", "F")
	}
}

// badge draws a coloured pill label.
func badge(pdf *gofpdf.Fpdf, x, y float64, label string, c color) {
	fill(pdf, c)
	pdf.RoundedRect(x, y+1.8, 47, 5.5, 2, "1234", "F")
	pdf.SetFont("Helvetica", "B", 6.5)
	text(pdf, cWhite)
	pdf.SetXY(x, y+2)
	pdf.CellFormat(47, 5.5, strings.ToUpper(label), "", 0, "C", false, 0, "")
}

// kpiCard draws a summary metric card. Returns nothing — caller advances Y manually.
func kpiCard(pdf *gofpdf.Fpdf, x, y, w, h float64, value, label, sub string, stripe color) {
	fill(pdf, cSubtle)
	draw(pdf, cBorder)
	pdf.SetLineWidth(0.3)
	pdf.RoundedRect(x, y, w, h, 3, "1234", "FD")
	// top colour stripe
	fill(pdf, stripe)
	pdf.RoundedRect(x, y, w, 3, 1.5, "1234", "F")
	pdf.Rect(x, y+1.5, w, 2, "F")
	// value
	pdf.SetFont("Helvetica", "B", 16)
	text(pdf, stripe)
	pdf.SetXY(x+4, y+5)
	pdf.Cell(w-8, 8, value)
	// label
	pdf.SetFont("Helvetica", "B", 6.5)
	text(pdf, cDark)
	pdf.SetXY(x+4, y+15)
	pdf.Cell(w-8, 4, strings.ToUpper(label))
	// sub
	pdf.SetFont("Helvetica", "", 6.5)
	text(pdf, cMuted)
	pdf.SetXY(x+4, y+20)
	pdf.Cell(w-8, 4, sub)
}

// donutChart draws a donut chart + legend side by side, centred between x and x+w.
// It returns the Y position below the chart block.
func donutChart(pdf *gofpdf.Fpdf, startY float64, labels []string, values []float64, colors []color) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	if total == 0 {
		return startY
	}

	cx := margin + 36.0
	cy := startY + 36.0
	outerR := 32.0
	innerR := 18.0

	angle := -math.Pi / 2
	for i, v := range values {
		if v == 0 {
			continue
		}
		sweep := 2 * math.Pi * v / total
		c := colors[i%len(colors)]
		fill(pdf, c)
		draw(pdf, cWhite)
		pdf.SetLineWidth(0.8)
		pts := donutSegmentPoints(cx, cy, outerR, innerR, angle, angle+sweep, 36)
		if len(pts) >= 3 {
			pdf.Polygon(pts, "FD")
		}
		angle += sweep
	}
	// white hole — polygon instead of Circle to avoid cursor drift that causes extra pages
	fill(pdf, cWhite)
	draw(pdf, cWhite)
	pdf.SetLineWidth(0)
	pdf.Polygon(circlePoints(cx, cy, innerR-0.5, 48), "F")

	// total label inside donut
	pdf.SetFont("Helvetica", "B", 10)
	text(pdf, cDark)
	pdf.SetXY(cx-15, cy-4)
	pdf.CellFormat(30, 5, fmt.Sprintf("%.0f", total), "", 0, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 6.5)
	text(pdf, cMuted)
	pdf.SetXY(cx-15, cy+2)
	pdf.CellFormat(30, 4, "total", "", 0, "C", false, 0, "")

	// legend (right of donut)
	lx := margin + 80.0
	ly := startY + 4
	colW := (inner - (lx - margin)) / 2
	col := 0
	for i, lbl := range labels {
		if values[i] == 0 {
			continue
		}
		pct := values[i] / total * 100
		itemX := lx + float64(col)*colW
		fill(pdf, colors[i%len(colors)])
		draw(pdf, colors[i%len(colors)])
		pdf.SetLineWidth(0)
		pdf.Polygon(circlePoints(itemX+3, ly+3.5, 2.5, 12), "F")
		pdf.SetFont("Helvetica", "", 7.5)
		text(pdf, cDark)
		pdf.SetXY(itemX+8, ly)
		pdf.Cell(colW-10, 7, fmt.Sprintf("%s  %s (%.0f%%)", lbl, fmtFloat(values[i]), pct))
		ly += 8
		if ly > startY+68 {
			ly = startY + 4
			col++
		}
	}

	return cy + outerR + 10
}

// circlePoints returns polygon points approximating a circle (used to avoid
// pdf.Circle which advances the internal cursor and triggers spurious page breaks).
func circlePoints(cx, cy, r float64, steps int) []gofpdf.PointType {
	pts := make([]gofpdf.PointType, steps)
	for i := 0; i < steps; i++ {
		a := 2 * math.Pi * float64(i) / float64(steps)
		pts[i] = gofpdf.PointType{X: cx + r*math.Cos(a), Y: cy + r*math.Sin(a)}
	}
	return pts
}

func donutSegmentPoints(cx, cy, outerR, innerR, start, end float64, steps int) []gofpdf.PointType {
	pts := make([]gofpdf.PointType, 0, steps*2+2)
	for i := 0; i <= steps; i++ {
		a := start + float64(i)*(end-start)/float64(steps)
		pts = append(pts, gofpdf.PointType{X: cx + outerR*math.Cos(a), Y: cy + outerR*math.Sin(a)})
	}
	for i := steps; i >= 0; i-- {
		a := start + float64(i)*(end-start)/float64(steps)
		pts = append(pts, gofpdf.PointType{X: cx + innerR*math.Cos(a), Y: cy + innerR*math.Sin(a)})
	}
	return pts
}

// drawToothLogo draws a simplified tooth shape at position (x,y) with given size.
func drawToothLogo(pdf *gofpdf.Fpdf, x, y, size float64) {
	cx := x + size/2
	cy := y + size*0.38
	fill(pdf, cWhite)
	draw(pdf, color{180, 210, 255})
	pdf.SetLineWidth(0.7)
	pdf.Ellipse(cx, cy, size*0.42, size*0.30, 0, "FD")
	pdf.Ellipse(cx-size*0.13, cy+size*0.28, size*0.12, size*0.17, 0, "FD")
	pdf.Ellipse(cx+size*0.13, cy+size*0.28, size*0.12, size*0.17, 0, "FD")
	// shine
	fill(pdf, color{200, 225, 255})
	pdf.Ellipse(cx-size*0.10, cy-size*0.08, size*0.07, size*0.05, 0, "F")
}

// starRating draws filled/empty star characters to represent a 0-5 rating.
func starRating(pdf *gofpdf.Fpdf, x, y, rating float64) {
	pdf.SetFont("Helvetica", "", 7)
	full := int(math.Round(rating))
	for i := 0; i < 5; i++ {
		if i < full {
			text(pdf, color{245, 158, 11})
			pdf.SetXY(x+float64(i)*5, y)
			pdf.Cell(5, rowH, "*")
		} else {
			text(pdf, cBorder)
			pdf.SetXY(x+float64(i)*5, y)
			pdf.Cell(5, rowH, "-")
		}
	}
}

// fmtFloat formats a float for display: no decimals if whole, 2 decimals otherwise.
func fmtFloat(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.2f", v)
}

func statusColor(s string) color {
	if c, ok := statusPalette[strings.ToLower(s)]; ok {
		return c
	}
	return cMuted
}

func humanStatus(s string) string {
	if s == "" {
		return "unknown"
	}
	return strings.ReplaceAll(s, "_", " ")
}

// ── Appointment report ────────────────────────────────────────────────────────

func buildAppointmentPDF(title, from, to string, rows []models.AppointmentReportRow) ([]byte, error) {
	pdf := newPDF()
	y := renderHeader(pdf, title, from, to)

	// KPI cards
	total := 0
	completed := 0
	cancelled := 0
	for _, r := range rows {
		total += r.AppointmentCount
		switch strings.ToLower(r.Status) {
		case "completed":
			completed += r.AppointmentCount
		case "cancelled", "no_show":
			cancelled += r.AppointmentCount
		}
	}
	compRate := 0.0
	if total > 0 {
		compRate = float64(completed) / float64(total) * 100
	}
	cardW := (inner - 6) / 3
	kpiCard(pdf, margin, y, cardW, 28, fmt.Sprintf("%d", total), "Total Appointments", "in selected period", cPrimary)
	kpiCard(pdf, margin+cardW+3, y, cardW, 28, fmt.Sprintf("%.0f%%", compRate), "Completion Rate", fmt.Sprintf("%d completed", completed), cAccent)
	kpiCard(pdf, margin+2*(cardW+3), y, cardW, 28, fmt.Sprintf("%d", cancelled), "Cancelled / No-show", fmt.Sprintf("%.0f%% of total", safePct(cancelled, total)), color{239, 68, 68})
	y += 36

	// Breakdown table
	y = sectionHeading(pdf, y, "Appointment Breakdown by Status")
	barMaxW := 44.0
	cols := []tableCol{
		{margin + 2, 50, "Status", "L"},
		{margin + 56, 20, "Count", "R"},
		{margin + 80, 22, "Share (%)", "R"},
		{margin + 108, barMaxW + 2, "Visual Distribution", "L"},
	}
	y = tableHeader(pdf, y, cols)

	for i, r := range rows {
		tableRowBand(pdf, y, i)
		pct := safePct(r.AppointmentCount, total)
		sc := statusColor(r.Status)
		badge(pdf, margin+2, y, humanStatus(r.Status), sc)
		pdf.SetFont("Helvetica", "B", 8.5)
		text(pdf, cDark)
		pdf.SetXY(cols[1].x, y)
		pdf.CellFormat(cols[1].w, rowH, fmt.Sprintf("%d", r.AppointmentCount), "", 0, "R", false, 0, "")
		pdf.SetFont("Helvetica", "", 8)
		text(pdf, cMuted)
		pdf.SetXY(cols[2].x, y)
		pdf.CellFormat(cols[2].w, rowH, fmt.Sprintf("%.1f%%", pct), "", 0, "R", false, 0, "")
		progressBar(pdf, cols[3].x, y, barMaxW, pct, sc)
		y += rowH
	}
	y += 10

	// Donut chart
	y = sectionHeading(pdf, y, "Status Distribution")
	labels := make([]string, len(rows))
	values := make([]float64, len(rows))
	colors := make([]color, len(rows))
	for i, r := range rows {
		labels[i] = humanStatus(r.Status)
		values[i] = float64(r.AppointmentCount)
		colors[i] = statusColor(r.Status)
	}
	donutChart(pdf, y, labels, values, colors)

	renderFooter(pdf)
	return pdfBytes(pdf)
}

// ── Revenue report ────────────────────────────────────────────────────────────

func buildRevenuePDF(title, from, to string, rows []models.RevenueReportRow) ([]byte, error) {
	pdf := newPDF()
	y := renderHeader(pdf, title, from, to)

	totalRevenue := 0.0
	totalAppts := 0
	for _, r := range rows {
		totalRevenue += r.TotalRevenue
		totalAppts += r.AppointmentCount
	}
	avgRevPerAppt := 0.0
	if totalAppts > 0 {
		avgRevPerAppt = totalRevenue / float64(totalAppts)
	}

	cardW := (inner - 6) / 3
	kpiCard(pdf, margin, y, cardW, 28, fmt.Sprintf("$%.0f", totalRevenue), "Total Revenue", "across all services", cPrimary)
	kpiCard(pdf, margin+cardW+3, y, cardW, 28, fmt.Sprintf("%d", totalAppts), "Appointments", "in selected period", cAccent)
	kpiCard(pdf, margin+2*(cardW+3), y, cardW, 28, fmt.Sprintf("$%.0f", avgRevPerAppt), "Avg Revenue / Appt", "across services", color{139, 92, 246})
	y += 36

	y = sectionHeading(pdf, y, "Revenue by Service")
	barMaxW := 35.0
	cols := []tableCol{
		{margin + 2, 62, "Service", "L"},
		{margin + 68, 22, "Appts", "R"},
		{margin + 94, 24, "Unit Price", "R"},
		{margin + 122, 26, "Total ($)", "R"},
		{margin + 152, barMaxW, "Share", "L"},
	}
	y = tableHeader(pdf, y, cols)

	for i, r := range rows {
		tableRowBand(pdf, y, i)
		pct := safePctF(r.TotalRevenue, totalRevenue)
		c := chartColors[i%len(chartColors)]

		pdf.SetFont("Helvetica", "", 8)
		text(pdf, cDark)
		pdf.SetXY(cols[0].x, y)
		pdf.CellFormat(cols[0].w, rowH, truncate(utf8safe(r.ServiceName), 30), "", 0, "L", false, 0, "")

		pdf.SetXY(cols[1].x, y)
		pdf.CellFormat(cols[1].w, rowH, fmt.Sprintf("%d", r.AppointmentCount), "", 0, "R", false, 0, "")

		text(pdf, cMuted)
		pdf.SetXY(cols[2].x, y)
		pdf.CellFormat(cols[2].w, rowH, fmt.Sprintf("$%.2f", r.UnitPrice), "", 0, "R", false, 0, "")

		pdf.SetFont("Helvetica", "B", 8)
		text(pdf, cDark)
		pdf.SetXY(cols[3].x, y)
		pdf.CellFormat(cols[3].w, rowH, fmt.Sprintf("$%.2f", r.TotalRevenue), "", 0, "R", false, 0, "")

		progressBar(pdf, cols[4].x, y, barMaxW, pct, c)
		y += rowH
	}
	y += 10

	// Donut
	y = sectionHeading(pdf, y, "Revenue Distribution by Service")
	labels := make([]string, len(rows))
	values := make([]float64, len(rows))
	colors := make([]color, len(rows))
	for i, r := range rows {
		labels[i] = truncate(utf8safe(r.ServiceName), 18)
		values[i] = r.TotalRevenue
		colors[i] = chartColors[i%len(chartColors)]
	}
	donutChart(pdf, y, labels, values, colors)

	renderFooter(pdf)
	return pdfBytes(pdf)
}

// ── Doctor performance report ─────────────────────────────────────────────────

func buildDoctorPDF(title, from, to string, rows []models.DoctorPerformanceRow) ([]byte, error) {
	pdf := newPDF()
	y := renderHeader(pdf, title, from, to)

	totalRevenue := 0.0
	totalAppts := 0
	avgRating := 0.0
	for _, r := range rows {
		totalRevenue += r.Revenue
		totalAppts += r.AppointmentCount
		avgRating += r.AverageRating
	}
	if len(rows) > 0 {
		avgRating /= float64(len(rows))
	}

	cardW := (inner - 6) / 3
	kpiCard(pdf, margin, y, cardW, 28, fmt.Sprintf("%d", len(rows)), "Doctors", "in this period", cPrimary)
	kpiCard(pdf, margin+cardW+3, y, cardW, 28, fmt.Sprintf("$%.0f", totalRevenue), "Total Revenue", fmt.Sprintf("%d appointments", totalAppts), cAccent)
	kpiCard(pdf, margin+2*(cardW+3), y, cardW, 28, fmt.Sprintf("%.2f / 5", avgRating), "Avg Rating", "across all doctors", color{245, 158, 11})
	y += 36

	y = sectionHeading(pdf, y, "Doctor Performance Breakdown")
	cols := []tableCol{
		{margin + 2, 44, "Doctor", "L"},
		{margin + 50, 34, "Specialization", "L"},
		{margin + 88, 18, "Appts", "R"},
		{margin + 110, 20, "Completed", "R"},
		{margin + 134, 24, "Revenue ($)", "R"},
		{margin + 162, 27, "Rating", "L"},
	}
	y = tableHeader(pdf, y, cols)

	maxRevenue := 0.0
	for _, r := range rows {
		if r.Revenue > maxRevenue {
			maxRevenue = r.Revenue
		}
	}

	for i, r := range rows {
		tableRowBand(pdf, y, i)

		compPct := safePct(r.CompletedCount, r.AppointmentCount)
		compColor := cAccent
		if compPct < 50 {
			compColor = color{239, 68, 68}
		} else if compPct < 75 {
			compColor = color{245, 158, 11}
		}

		pdf.SetFont("Helvetica", "B", 7.5)
		text(pdf, cDark)
		pdf.SetXY(cols[0].x, y)
		pdf.CellFormat(cols[0].w, rowH, truncate(utf8safe(r.DoctorName), 18), "", 0, "L", false, 0, "")

		pdf.SetFont("Helvetica", "", 7.5)
		text(pdf, cMuted)
		pdf.SetXY(cols[1].x, y)
		pdf.CellFormat(cols[1].w, rowH, truncate(utf8safe(r.Specialization), 16), "", 0, "L", false, 0, "")

		pdf.SetFont("Helvetica", "", 8)
		text(pdf, cDark)
		pdf.SetXY(cols[2].x, y)
		pdf.CellFormat(cols[2].w, rowH, fmt.Sprintf("%d", r.AppointmentCount), "", 0, "R", false, 0, "")

		text(pdf, compColor)
		pdf.SetXY(cols[3].x, y)
		pdf.CellFormat(cols[3].w, rowH, fmt.Sprintf("%d (%.0f%%)", r.CompletedCount, compPct), "", 0, "R", false, 0, "")

		pdf.SetFont("Helvetica", "B", 8)
		text(pdf, cDark)
		pdf.SetXY(cols[4].x, y)
		pdf.CellFormat(cols[4].w, rowH, fmt.Sprintf("$%.0f", r.Revenue), "", 0, "R", false, 0, "")

		starRating(pdf, cols[5].x, y, r.AverageRating)
		pdf.SetFont("Helvetica", "", 7)
		text(pdf, cMuted)
		pdf.SetXY(cols[5].x+26, y)
		pdf.Cell(10, rowH, fmt.Sprintf("%.1f", r.AverageRating))

		y += rowH
	}
	y += 10

	// Revenue bar chart (simple horizontal bars)
	if len(rows) > 0 {
		y = sectionHeading(pdf, y, "Revenue by Doctor")
		barTrackW := inner - 60.0
		for i, r := range rows {
			if maxRevenue == 0 {
				break
			}
			barFillW := barTrackW * r.Revenue / maxRevenue
			c := chartColors[i%len(chartColors)]

			pdf.SetFont("Helvetica", "", 7.5)
			text(pdf, cDark)
			pdf.SetXY(margin, y)
			pdf.CellFormat(56, 8, truncate(utf8safe(r.DoctorName), 22), "", 0, "L", false, 0, "")

			fill(pdf, cBorder)
			pdf.RoundedRect(margin+58, y+2, barTrackW, 4, 1.5, "1234", "F")
			if barFillW > 0 {
				fill(pdf, c)
				pdf.RoundedRect(margin+58, y+2, barFillW, 4, 1.5, "1234", "F")
			}
			text(pdf, cMuted)
			pdf.SetXY(margin+58+barTrackW+2, y)
			pdf.Cell(20, 8, fmt.Sprintf("$%.0f", r.Revenue))
			y += 8
		}
	}

	renderFooter(pdf)
	return pdfBytes(pdf)
}

// ── Inventory report ──────────────────────────────────────────────────────────

func buildInventoryPDF(title, from, to string, rows []models.InventoryReportRow) ([]byte, error) {
	pdf := newPDF()
	y := renderHeader(pdf, title, from, to)

	totalProducts := len(rows)
	totalStock := 0.0
	totalRestocked := 0.0
	totalUsed := 0.0
	for _, r := range rows {
		totalStock += r.CurrentQuantity
		totalRestocked += r.RestockedQuantity
		totalUsed += r.UsedQuantity
	}

	cardW := (inner - 6) / 3
	kpiCard(pdf, margin, y, cardW, 28, fmt.Sprintf("%d", totalProducts), "Products Tracked", "in inventory", cPrimary)
	kpiCard(pdf, margin+cardW+3, y, cardW, 28, fmtFloat(totalRestocked), "Units Restocked", "in period", cAccent)
	kpiCard(pdf, margin+2*(cardW+3), y, cardW, 28, fmtFloat(totalUsed), "Units Used", "in period", color{245, 158, 11})
	y += 36

	y = sectionHeading(pdf, y, "Inventory Status by Product")
	cols := []tableCol{
		{margin + 2, 50, "Product", "L"},
		{margin + 56, 14, "Unit", "C"},
		{margin + 74, 26, "Current Stock", "R"},
		{margin + 104, 24, "Restocked", "R"},
		{margin + 132, 20, "Used", "R"},
		{margin + 156, 24, "Adjustment", "R"},
	}
	y = tableHeader(pdf, y, cols)

	maxStock := 0.0
	for _, r := range rows {
		if r.CurrentQuantity > maxStock {
			maxStock = r.CurrentQuantity
		}
	}

	for i, r := range rows {
		tableRowBand(pdf, y, i)

		// Stock level indicator colour
		stockColor := cAccent
		if maxStock > 0 {
			ratio := r.CurrentQuantity / maxStock
			if ratio < 0.2 {
				stockColor = color{239, 68, 68}
			} else if ratio < 0.5 {
				stockColor = color{245, 158, 11}
			}
		}

		pdf.SetFont("Helvetica", "B", 7.5)
		text(pdf, cDark)
		pdf.SetXY(cols[0].x, y)
		pdf.CellFormat(cols[0].w, rowH, truncate(utf8safe(r.ProductName), 22), "", 0, "L", false, 0, "")

		pdf.SetFont("Helvetica", "", 7.5)
		text(pdf, cMuted)
		pdf.SetXY(cols[1].x, y)
		pdf.CellFormat(cols[1].w, rowH, utf8safe(r.Unit), "", 0, "C", false, 0, "")

		pdf.SetFont("Helvetica", "B", 8)
		text(pdf, stockColor)
		pdf.SetXY(cols[2].x, y)
		pdf.CellFormat(cols[2].w, rowH, fmtFloat(r.CurrentQuantity), "", 0, "R", false, 0, "")

		pdf.SetFont("Helvetica", "", 8)
		text(pdf, cAccent)
		pdf.SetXY(cols[3].x, y)
		pdf.CellFormat(cols[3].w, rowH, "+"+fmtFloat(r.RestockedQuantity), "", 0, "R", false, 0, "")

		text(pdf, color{239, 68, 68})
		pdf.SetXY(cols[4].x, y)
		pdf.CellFormat(cols[4].w, rowH, "-"+fmtFloat(r.UsedQuantity), "", 0, "R", false, 0, "")

		text(pdf, cMuted)
		pdf.SetXY(cols[5].x, y)
		pdf.CellFormat(cols[5].w, rowH, fmtFloat(r.AdjustmentQuantity), "", 0, "R", false, 0, "")

		y += rowH
	}
	y += 10

	// Stock level overview bar chart
	if len(rows) > 0 && maxStock > 0 {
		y = sectionHeading(pdf, y, "Current Stock Levels")
		barTrackW := inner - 72.0
		for i, r := range rows {
			barFillW := barTrackW * r.CurrentQuantity / maxStock
			ratio := r.CurrentQuantity / maxStock
			c := cAccent
			if ratio < 0.2 {
				c = color{239, 68, 68}
			} else if ratio < 0.5 {
				c = color{245, 158, 11}
			}

			pdf.SetFont("Helvetica", "", 7.5)
			text(pdf, cDark)
			pdf.SetXY(margin, y)
			pdf.CellFormat(68, 8, truncate(utf8safe(r.ProductName)+" ("+utf8safe(r.Unit)+")", 28), "", 0, "L", false, 0, "")

			fill(pdf, cBorder)
			pdf.RoundedRect(margin+70, y+2, barTrackW, 4, 1.5, "1234", "F")
			if barFillW > 0 {
				fill(pdf, c)
				pdf.RoundedRect(margin+70, y+2, barFillW, 4, 1.5, "1234", "F")
			}
			text(pdf, cMuted)
			pdf.SetXY(margin+70+barTrackW+2, y)
			pdf.Cell(16, 8, fmtFloat(r.CurrentQuantity))
			y += 8
			_ = i
		}
	}

	renderFooter(pdf)
	return pdfBytes(pdf)
}

// ── Utilities ─────────────────────────────────────────────────────────────────

func pdfBytes(pdf *gofpdf.Fpdf) ([]byte, error) {
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func safePct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(n) / float64(total) * 100
}

func safePctF(v, total float64) float64 {
	if total == 0 {
		return 0
	}
	return v / total * 100
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "\u2026"
}

// ── tableRows (for CSV, unchanged) ───────────────────────────────────────────

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
