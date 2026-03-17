package dto

import (
	// "time"

)

type CreateAppointmentRequest struct {
	Doctor_id   			string    	`json:"doctor_id"`
	Clinic_address_id      	string    	`json:"clinic_address_id"`
	Service_id      		string    	`json:"service_id"`
	Slot_id      			string    	`json:"slot_id"`
	Date	      			string    	`json:"date"`
	Name	      			string    	`json:"name"`
	Email	      			string    	`json:"email"`
}

type CreateAppointmentResponse struct {
	Success 		string `json:"success"`
	Message 		string `json:"message"`
	Appointment_id  string `json:"appointment_id"`
}