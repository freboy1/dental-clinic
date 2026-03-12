package dto

import (
	"time"

)

type CreateScheduleRequest struct {
	// Doctor_id   			string    `json:"doctor_id"`
	Clinic_address_id      	string    `json:"clinic_address_id"`
	Day_of_week    			int    `json:"day_of_week"`
	Start_time  			string    `json:"start_time"`
	End_time  				string   `json:"end_time"`
}

type CreateScheduleResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Schedule_id  string `json:"schedule_id"`
}

type GenerateSlotsRequest struct {
	From_date      	string    `json:"from_date"`
	To_date    		string    `json:"to_date"`
}

type ScheduleResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

type SlotResponse struct {
	Id        string `json:"id"` 
	Slot_start   time.Time `json:"slot_start"` 
	Slot_end      time.Time `json:"slot_end"` 
	Status    string `json:"status"` 
}

type ScheduleDoctorResponse struct {
	Id        				string  
	Clinic_address_id      	string    
	Day_of_week    			int    
	Start_time  			string    
	End_time  				string   
}