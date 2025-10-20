package models

import "time"

type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Phone           string    `json:"phone,omitempty"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Role            string    `json:"role"`
	PushConsent     bool      `json:"push_consent"`
	Activated       bool      `json:"activated"`
	RegisteredAt    time.Time `json:"registered_at"`
	ActivationToken string    `json:"-"`
}
