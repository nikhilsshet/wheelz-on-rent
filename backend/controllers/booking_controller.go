package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nikhilsshet/wheelz-on-rent/backend/config"
	"github.com/nikhilsshet/wheelz-on-rent/backend/middleware"
)

type BookingInput struct {
	VehicleID int    `json:"vehicle_id"`
	StartDate string `json:"start_date"` // format: "YYYY-MM-DD"
	EndDate   string `json:"end_date"`
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middleware.UserIDKey).(int)
    role := r.Context().Value(middleware.UserRoleKey).(string)

    fmt.Printf("User ID: %d, Role: %s\n", userID, role)
	
	var input BookingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Optional: Get customer ID from token context
	// userID, ok := r.Context().Value("userID").(float64)
	// if !ok {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	startDate, err1 := time.Parse("2006-01-02", input.StartDate)
	endDate, err2 := time.Parse("2006-01-02", input.EndDate)
	if err1 != nil || err2 != nil || endDate.Before(startDate) {
		http.Error(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	db := config.GetDB()

	// Get price of vehicle
	var pricePerDay float64
	err := db.QueryRow("SELECT price_per_day FROM vehicles WHERE id = $1", input.VehicleID).Scan(&pricePerDay)
	if err != nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	duration := endDate.Sub(startDate).Hours() / 24
	total := pricePerDay * duration

	// Insert booking
	_, err = db.Exec(`
		INSERT INTO bookings (customer_id, vehicle_id, start_date, end_date, total_price)
		VALUES ($1, $2, $3, $4, $5)
	`, int(userID), input.VehicleID, startDate, endDate, total)

	if err != nil {
		http.Error(w, "Could not create booking", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Booking created successfully",
		"total":   strconv.FormatFloat(total, 'f', 2, 64),
	})
}
