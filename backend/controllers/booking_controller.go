package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

	if role != "customer" {
		http.Error(w, "Only customers can book vehicles", http.StatusForbidden)
		return
	}

	var input BookingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	startDate, err1 := time.Parse("2006-01-02", input.StartDate)
	endDate, err2 := time.Parse("2006-01-02", input.EndDate)
	if err1 != nil || err2 != nil || endDate.Before(startDate) {
		http.Error(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	db := config.GetDB()

	// Step 1: Check vehicle availability
	var pricePerDay float64
	var available bool
	err := db.QueryRow("SELECT price_per_day, availability FROM vehicles WHERE id = $1", input.VehicleID).
		Scan(&pricePerDay, &available)

	if err != nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	if !available {
		http.Error(w, "Vehicle is not available", http.StatusBadRequest)
		return
	}

	// Step 2: Calculate total
	duration := endDate.Sub(startDate).Hours() / 24
	if duration < 1 {
		duration = 1 // at least 1 day
	}
	total := pricePerDay * duration

	// Step 3: Insert booking
	_, err = db.Exec(`
		INSERT INTO bookings (customer_id, vehicle_id, start_date, end_date, total_price, status, payment_status)
		VALUES ($1, $2, $3, $4, $5, 'active')
 		`, userID, input.VehicleID, startDate, endDate, total)

	if err != nil {
		http.Error(w, "Could not create booking", http.StatusInternalServerError)
		return
	}

	// Step 4: Mark vehicle as unavailable
	_, err = db.Exec("UPDATE vehicles SET availability = false WHERE id = $1", input.VehicleID)
	if err != nil {
		http.Error(w, "Failed to update vehicle availability", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Booking created successfully and vehicle marked unavailable",
		"total":   strconv.FormatFloat(total, 'f', 2, 64),
	})
}

func GetMyBookings(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(middleware.UserIDKey)
	roleVal := r.Context().Value(middleware.UserRoleKey)

	fmt.Println("UserID context value:", userIDVal)
	fmt.Println("Role context value:", roleVal)

	if userIDVal == nil || roleVal == nil {
		http.Error(w, "Unauthorized: missing user context", http.StatusUnauthorized)
		return
	}

	userID := userIDVal.(int)
	role := roleVal.(string)

	if role != "customer" {
		http.Error(w, "Forbidden: Only customers can view their bookings", http.StatusForbidden)
		return
	}

	db := config.GetDB()

	rows, err := db.Query(`
		SELECT b.id, b.start_date, b.end_date, b.total_price, b.status,
				v.name, v.type, v.model, v.number_plate
		FROM bookings b
		JOIN vehicles v ON b.vehicle_id = v.id
		WHERE b.customer_id = $1
		ORDER BY b.created_at DESC
 `, userID)

	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type BookingResponse struct {
		BookingID   int     `json:"booking_id"`
		StartDate   string  `json:"start_date"`
		EndDate     string  `json:"end_date"`
		TotalPrice  float64 `json:"total_price"`
		Status      string  `json:"status"`
		VehicleName string  `json:"vehicle_name"`
		VehicleType string  `json:"vehicle_type"`
		Model       string  `json:"model"`
		NumberPlate string  `json:"number_plate"`
	}

	var bookings []BookingResponse

	for rows.Next() {
		var b BookingResponse
		err := rows.Scan(&b.BookingID, &b.StartDate, &b.EndDate, &b.TotalPrice, &b.Status,
			&b.VehicleName, &b.VehicleType, &b.Model, &b.NumberPlate)
		if err != nil {
			http.Error(w, "Failed to read bookings", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	role := r.Context().Value(middleware.UserRoleKey).(string)

	vars := mux.Vars(r)
	bookingIDStr := vars["id"]
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()

	// Step 1: Get booking info
	var customerID int
	var status string
	var vehicleID int

	err = db.QueryRow(`SELECT customer_id, status, vehicle_id FROM bookings WHERE id = $1`, bookingID).
		Scan(&customerID, &status, &vehicleID)
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Step 2: Authorization check
	if userID != customerID && role != "admin" {
		http.Error(w, "Unauthorized to cancel this booking", http.StatusUnauthorized)
		return
	}

	if status == "cancelled" {
		http.Error(w, "Booking is already cancelled", http.StatusBadRequest)
		return
	}

	// Step 3: Cancel the booking
	_, err = db.Exec(`UPDATE bookings SET status = 'cancelled' WHERE id = $1`, bookingID)
	if err != nil {
		http.Error(w, "Failed to cancel booking", http.StatusInternalServerError)
		return
	}

	// Step 4: Mark vehicle as available again
	_, err = db.Exec(`UPDATE vehicles SET availability = true WHERE id = $1`, vehicleID)
	if err != nil {
		http.Error(w, "Failed to update vehicle availability", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Booking cancelled and vehicle marked available",
	})
}

func GetAllBookings(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.UserRoleKey).(string)

	if role != "admin" {
		http.Error(w, "Forbidden: Only admin can view all bookings", http.StatusForbidden)
		return
	}

	db := config.GetDB()

	rows, err := db.Query(`
		SELECT b.id, b.start_date, b.end_date, b.total_price, b.status,
		       v.name, v.type, v.model, v.number_plate,
		       u.id, u.name, u.email
		FROM bookings b
		JOIN vehicles v ON b.vehicle_id = v.id
		JOIN users u ON b.customer_id = u.id
		ORDER BY b.created_at DESC
	`)

	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type BookingAdminView struct {
		BookingID     int     `json:"booking_id"`
		StartDate     string  `json:"start_date"`
		EndDate       string  `json:"end_date"`
		TotalPrice    float64 `json:"total_price"`
		Status        string  `json:"status"`
		VehicleName   string  `json:"vehicle_name"`
		VehicleType   string  `json:"vehicle_type"`
		VehicleModel  string  `json:"vehicle_model"`
		NumberPlate   string  `json:"number_plate"`
		CustomerID    int     `json:"customer_id"`
		CustomerName  string  `json:"customer_name"`
		CustomerEmail string  `json:"customer_email"`
	}

	var bookings []BookingAdminView

	for rows.Next() {
		var b BookingAdminView
		err := rows.Scan(&b.BookingID, &b.StartDate, &b.EndDate, &b.TotalPrice, &b.Status,
			&b.VehicleName, &b.VehicleType, &b.VehicleModel, &b.NumberPlate,
			&b.CustomerID, &b.CustomerName, &b.CustomerEmail)
		if err != nil {
			http.Error(w, "Failed to read bookings", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}
