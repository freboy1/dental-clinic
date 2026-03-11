package models

import (
	// "time"

	"github.com/google/uuid"
)

type Schedule struct {
	Id          		uuid.UUID 
	Doctor_id       	uuid.UUID    
	Clinic_address_id 	uuid.UUID 
	Day_of_week       	int     
	Start_time   		string 
	End_time   			string 
}
