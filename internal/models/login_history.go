package models

import "time"

type LoginHistory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	IPAddress   string    `json:"ip_address"`
	Success     bool      `json:"success"`
	UserAgent   string    `json:"user_agent"`
	AttemptTime time.Time `json:"attempt_time"`
}
