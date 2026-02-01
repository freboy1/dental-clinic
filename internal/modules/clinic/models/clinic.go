package models

import (
	"time"

	"github.com/google/uuid"
)

type Clinic struct {
	Id          uuid.UUID 
	Name        string    
	Description string    
	Phone       string    
	Email       string    
	Website     string    
	IsActive    bool      
	CreatedAt   time.Time 
}
