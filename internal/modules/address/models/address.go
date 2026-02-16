package models

import (
	"github.com/google/uuid"
)

type Address struct {
	ID        uuid.UUID 
	Country   string    
	City      string    
	Street    string    
	Building  string    
	Latitude  float64   
	Longitude float64   
}
