package models

import "github.com/google/uuid"

type User struct {
	Id uuid.UUID
	Role string
	Email string
	Name string
	Gender string
	Age int
	Push_consent bool
}