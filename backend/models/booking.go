package models

import "time"

type Booking struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	VehicleID  int       `json:"vehicle_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
