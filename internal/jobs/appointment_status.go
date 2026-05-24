package jobs

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type completedAppointment struct {
	id              uuid.UUID
	clinicAddressId uuid.UUID
	clinicServiceId uuid.UUID
}

type serviceMaterial struct {
	productId uuid.UUID
	quantity  float64
}

func StartAppointmentStatusCron(ctx context.Context, db *pgxpool.Pool, interval time.Duration) {
	if db == nil {
		return
	}
	if interval <= 0 {
		interval = time.Minute
	}

	go func() {
		runAppointmentStatusJob(ctx, db)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runAppointmentStatusJob(ctx, db)
			}
		}
	}()
}

func runAppointmentStatusJob(ctx context.Context, db *pgxpool.Pool) {
	updated, err := completeExpiredAppointments(ctx, db)
	if err != nil {
		log.Printf("appointment status cron failed: %v", err)
		return
	}
	if updated > 0 {
		log.Printf("appointment status cron completed %d appointment(s)", updated)
	}
}

func completeExpiredAppointments(ctx context.Context, db *pgxpool.Pool) (int64, error) {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	appointments, err := getExpiredBookedAppointments(ctx, tx)
	if err != nil {
		return 0, err
	}

	for _, appointment := range appointments {
		if appointment.clinicServiceId != uuid.Nil {
			materials, err := getServiceMaterials(ctx, tx, appointment.clinicServiceId)
			if err != nil {
				return 0, err
			}

			for _, material := range materials {
				if err := subtractInventoryMaterial(ctx, tx, appointment, material); err != nil {
					return 0, err
				}
			}
		}

		if err := markAppointmentCompleted(ctx, tx, appointment.id); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return int64(len(appointments)), nil
}

func getExpiredBookedAppointments(ctx context.Context, tx pgx.Tx) ([]completedAppointment, error) {
	query := `
		SELECT
			a.id,
			a.clinic_address_id,
			COALESCE(cs.id, '00000000-0000-0000-0000-000000000000'::uuid)
		FROM appointments a
		LEFT JOIN clinic_addresses ca ON ca.id = a.clinic_address_id
		LEFT JOIN clinic_services cs ON cs.clinic_id = ca.clinic_id
			AND cs.service_id = a.service_id
			AND cs.is_active = true
		WHERE a.status = 'booked'
			AND a.end_time <= NOW()
		FOR UPDATE OF a SKIP LOCKED
	`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := make([]completedAppointment, 0)
	for rows.Next() {
		var appointment completedAppointment
		if err := rows.Scan(&appointment.id, &appointment.clinicAddressId, &appointment.clinicServiceId); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}
	return appointments, rows.Err()
}

func getServiceMaterials(ctx context.Context, tx pgx.Tx, clinicServiceId uuid.UUID) ([]serviceMaterial, error) {
	query := `
		SELECT product_id, quantity_required
		FROM service_materials
		WHERE service_id = $1
	`
	rows, err := tx.Query(ctx, query, clinicServiceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	materials := make([]serviceMaterial, 0)
	for rows.Next() {
		var material serviceMaterial
		if err := rows.Scan(&material.productId, &material.quantity); err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}
	return materials, rows.Err()
}

func subtractInventoryMaterial(ctx context.Context, tx pgx.Tx, appointment completedAppointment, material serviceMaterial) error {
	query := `
		UPDATE address_inventory
		SET quantity = quantity - $3,
			updated_at = NOW()
		WHERE clinic_address_id = $1
			AND product_id = $2
	`
	result, err := tx.Exec(ctx, query, appointment.clinicAddressId, material.productId, material.quantity)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		insertQuery := `
			INSERT INTO address_inventory (id, clinic_address_id, product_id, quantity, updated_at)
			VALUES ($1, $2, $3, $4, NOW())
		`
		if _, err := tx.Exec(ctx, insertQuery, uuid.New(), appointment.clinicAddressId, material.productId, -material.quantity); err != nil {
			return err
		}
	}

	transactionQuery := `
		INSERT INTO inventory_transactions (id, clinic_address_id, product_id, quantity, transaction_type, appointment_id, created_at)
		VALUES ($1, $2, $3, $4, 'used', $5, NOW())
	`
	_, err = tx.Exec(ctx, transactionQuery, uuid.New(), appointment.clinicAddressId, material.productId, material.quantity, appointment.id)
	return err
}

func markAppointmentCompleted(ctx context.Context, tx pgx.Tx, appointmentId uuid.UUID) error {
	result, err := tx.Exec(ctx, `UPDATE appointments SET status = 'completed' WHERE id = $1 AND status = 'booked'`, appointmentId)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
