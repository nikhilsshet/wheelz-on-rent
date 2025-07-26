package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/nikhilsshet/wheelz-on-rent/backend/config"
	"github.com/nikhilsshet/wheelz-on-rent/backend/models"
)

func AddVehicle(w http.ResponseWriter, r *http.Request) {
	var vehicle models.Vehicle

	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = config.DB.Exec(`INSERT INTO vehicles 
		(name, type, model, number_plate, color, availability, price_per_day) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		vehicle.Name, vehicle.Type, vehicle.Model, vehicle.NumberPlate,
		vehicle.Color, vehicle.Availability, vehicle.PricePerDay)

	if err != nil {
		http.Error(w, "Failed to add vehicle", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle added successfully"})
}

func GetAllVehicles(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, name, type, model, number_plate, color, availability, price_per_day FROM vehicles")
	if err != nil {
		http.Error(w, "Failed to fetch vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		err := rows.Scan(&v.ID, &v.Name, &v.Type, &v.Model, &v.NumberPlate, &v.Color, &v.Availability, &v.PricePerDay)
		if err != nil {
			http.Error(w, "Error parsing vehicle", http.StatusInternalServerError)
			return
		}
		vehicles = append(vehicles, v)
	}

	json.NewEncoder(w).Encode(vehicles)
}

