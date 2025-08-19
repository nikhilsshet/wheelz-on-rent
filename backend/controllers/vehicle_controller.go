package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nikhilsshet/wheelz-on-rent/backend/config"
	"github.com/nikhilsshet/wheelz-on-rent/backend/models"
)

// Helper function to extract role from JWT token
func getRoleFromToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if role, ok := claims["role"].(string); ok {
			return role, nil
		}
	}
	return "", http.ErrNoCookie
}

func AddVehicle(w http.ResponseWriter, r *http.Request) {

	role, err := getRoleFromToken(r)
	if err != nil || role != "admin" {
		http.Error(w, "Unauthorized: Admins can only add vehicles", http.StatusNotAcceptable)
		return
	}

	var vehicle models.Vehicle

	err = json.NewDecoder(r.Body).Decode(&vehicle)
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
