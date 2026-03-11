package dto

import (
	"time"

)

type CreateScheduleRequest struct {
	// Doctor_id   			string    `json:"doctor_id"`
	Clinic_address_id      	string    `json:"clinic_address_id"`
	Day_of_week    			int    `json:"day_of_week"`
	Start_time  			time.Time    `json:"start_time"`
	End_time  				time.Time   `json:"end_time"`
}

type CreateScheduleResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Schedule_id  string `json:"schedule_id"`
}