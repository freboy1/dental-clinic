package repository

import (
	"context"

	"dental_clinic/internal/modules/reports/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportsRepository interface {
	GetRevenueReport(filters models.ReportFilters) ([]models.RevenueReportRow, error)
	GetAppointmentReport(filters models.ReportFilters) ([]models.AppointmentReportRow, error)
	GetDoctorPerformanceReport(filters models.ReportFilters) ([]models.DoctorPerformanceRow, error)
	GetInventoryReport(filters models.ReportFilters) ([]models.InventoryReportRow, error)
}

type reportsRepo struct {
	db *pgxpool.Pool
}

func NewReportsRepository(db *pgxpool.Pool) ReportsRepository {
	return &reportsRepo{db: db}
}

func (r *reportsRepo) GetRevenueReport(filters models.ReportFilters) ([]models.RevenueReportRow, error) {
	query := `
		SELECT
			s.id::text,
			s.name,
			COUNT(a.id)::int AS appointment_count,
			COALESCE(cs.price, 0)::float8 AS unit_price,
			COALESCE(SUM(cs.price), 0)::float8 AS total_revenue
		FROM appointments a
		JOIN clinic_addresses ca ON ca.id = a.clinic_address_id
		JOIN services s ON s.id = a.service_id
		LEFT JOIN clinic_services cs ON cs.clinic_id = ca.clinic_id AND cs.service_id = a.service_id
		WHERE ca.clinic_id = $1::uuid
			AND ($2 = '' OR a.clinic_address_id = $2::uuid)
			AND a.start_time >= $3::date
			AND a.start_time < ($4::date + INTERVAL '1 day')
		GROUP BY s.id, s.name, cs.price
		ORDER BY total_revenue DESC, s.name
	`
	rows, err := r.db.Query(context.Background(), query, filters.ClinicID, filters.ClinicAddressID, filters.From, filters.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.RevenueReportRow, 0)
	for rows.Next() {
		var row models.RevenueReportRow
		if err := rows.Scan(&row.ServiceID, &row.ServiceName, &row.AppointmentCount, &row.UnitPrice, &row.TotalRevenue); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func (r *reportsRepo) GetAppointmentReport(filters models.ReportFilters) ([]models.AppointmentReportRow, error) {
	query := `
		SELECT COALESCE(a.status, '') AS status, COUNT(a.id)::int AS appointment_count
		FROM appointments a
		JOIN clinic_addresses ca ON ca.id = a.clinic_address_id
		WHERE ca.clinic_id = $1::uuid
			AND ($2 = '' OR a.clinic_address_id = $2::uuid)
			AND a.start_time >= $3::date
			AND a.start_time < ($4::date + INTERVAL '1 day')
		GROUP BY a.status
		ORDER BY appointment_count DESC, a.status
	`
	rows, err := r.db.Query(context.Background(), query, filters.ClinicID, filters.ClinicAddressID, filters.From, filters.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.AppointmentReportRow, 0)
	for rows.Next() {
		var row models.AppointmentReportRow
		if err := rows.Scan(&row.Status, &row.AppointmentCount); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func (r *reportsRepo) GetDoctorPerformanceReport(filters models.ReportFilters) ([]models.DoctorPerformanceRow, error) {
	query := `
		SELECT
			d.id::text,
			d.name,
			d.specialization,
			COUNT(DISTINCT a.id)::int AS appointment_count,
			COUNT(DISTINCT a.id) FILTER (WHERE a.status = 'completed')::int AS completed_count,
			COALESCE(SUM(cs.price), 0)::float8 AS revenue,
			COALESCE(ratings.average_rating, 0)::float8 AS average_rating
		FROM doctors d
		JOIN clinic_addresses ca ON ca.clinic_id = d.clinic_id
		LEFT JOIN appointments a ON a.doctor_id = d.id AND a.clinic_address_id = ca.id
			AND a.start_time >= $3::date
			AND a.start_time < ($4::date + INTERVAL '1 day')
		LEFT JOIN clinic_services cs ON cs.clinic_id = ca.clinic_id AND cs.service_id = a.service_id
		LEFT JOIN (
			SELECT doctor_id, ROUND(AVG(rating)::numeric, 2)::float8 AS average_rating
			FROM doctor_ratings
			GROUP BY doctor_id
		) ratings ON ratings.doctor_id = d.id
		WHERE d.clinic_id = $1::uuid
			AND d.is_deleted = 0
			AND ($2 = '' OR ca.id = $2::uuid)
		GROUP BY d.id, d.name, d.specialization, ratings.average_rating
		ORDER BY revenue DESC, appointment_count DESC, d.name
	`
	rows, err := r.db.Query(context.Background(), query, filters.ClinicID, filters.ClinicAddressID, filters.From, filters.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.DoctorPerformanceRow, 0)
	for rows.Next() {
		var row models.DoctorPerformanceRow
		if err := rows.Scan(&row.DoctorID, &row.DoctorName, &row.Specialization, &row.AppointmentCount, &row.CompletedCount, &row.Revenue, &row.AverageRating); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func (r *reportsRepo) GetInventoryReport(filters models.ReportFilters) ([]models.InventoryReportRow, error) {
	query := `
		SELECT
			p.id::text,
			p.name,
			p.unit,
			ABS(COALESCE(inv.current_quantity, 0))::float8,
			ABS(COALESCE(tx.restocked_quantity, 0))::float8,
			ABS(COALESCE(tx.used_quantity, 0))::float8,
			COALESCE(tx.adjustment_quantity, 0)::float8
		FROM products p
		JOIN (
			SELECT ai.product_id, SUM(ai.quantity)::float8 AS current_quantity
			FROM address_inventory ai
			JOIN clinic_addresses ca ON ca.id = ai.clinic_address_id
			WHERE ca.clinic_id = $1::uuid
				AND ($2 = '' OR ca.id = $2::uuid)
			GROUP BY ai.product_id
		) inv ON inv.product_id = p.id
		LEFT JOIN (
			SELECT
				it.product_id,
				SUM(it.quantity) FILTER (WHERE it.transaction_type = 'restocked')::float8 AS restocked_quantity,
				SUM(it.quantity) FILTER (WHERE it.transaction_type = 'used')::float8 AS used_quantity,
				SUM(it.quantity) FILTER (WHERE it.transaction_type = 'manual_adjustment')::float8 AS adjustment_quantity
			FROM inventory_transactions it
			JOIN clinic_addresses ca ON ca.id = it.clinic_address_id
			WHERE ca.clinic_id = $1::uuid
				AND ($2 = '' OR ca.id = $2::uuid)
				AND it.created_at >= $3::date
				AND it.created_at < ($4::date + INTERVAL '1 day')
			GROUP BY it.product_id
		) tx ON tx.product_id = p.id
		ORDER BY p.name
	`
	rows, err := r.db.Query(context.Background(), query, filters.ClinicID, filters.ClinicAddressID, filters.From, filters.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.InventoryReportRow, 0)
	for rows.Next() {
		var row models.InventoryReportRow
		if err := rows.Scan(&row.ProductID, &row.ProductName, &row.Unit, &row.CurrentQuantity, &row.RestockedQuantity, &row.UsedQuantity, &row.AdjustmentQuantity); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
