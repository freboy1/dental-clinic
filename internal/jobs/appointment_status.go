package jobs

import (
	"context"
	"log"
	"time"

	"dental_clinic/internal/modules/appointment/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartAppointmentStatusCron(ctx context.Context, db *pgxpool.Pool, interval time.Duration) {
	if db == nil {
		return
	}
	if interval <= 0 {
		interval = time.Minute
	}

	repo := repository.NewAppointmentRepository(db)

	go func() {
		runAppointmentStatusJob(ctx, repo)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runAppointmentStatusJob(ctx, repo)
			}
		}
	}()
}

func runAppointmentStatusJob(ctx context.Context, repo repository.AppointmentRepository) {
	updated, err := repo.MarkExpiredBookedCompleted(ctx)
	if err != nil {
		log.Printf("appointment status cron failed: %v", err)
		return
	}
	if updated > 0 {
		log.Printf("appointment status cron completed %d appointment(s)", updated)
	}
}
