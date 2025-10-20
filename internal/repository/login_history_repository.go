package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginHistoryRepository interface {
	LogLogin(ctx context.Context, userID string, ip string, success bool, userAgent string) error
}

type loginHistoryRepository struct {
	db *pgxpool.Pool
}

func NewLoginHistoryRepository(db *pgxpool.Pool) LoginHistoryRepository {
	return &loginHistoryRepository{db: db}
}

func (r *loginHistoryRepository) LogLogin(ctx context.Context, userID string, ip string, success bool, userAgent string) error {
	query := `
		INSERT INTO login_history (user_id, ip_address, success, user_agent, attempt_time)
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err := r.db.Exec(ctx, query, userID, ip, success, userAgent, time.Now())
	return err
}
