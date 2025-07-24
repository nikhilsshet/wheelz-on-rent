package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/nikhilsshet/wheelz-on-rent/backend/config"
	"github.com/nikhilsshet/wheelz-on-rent/backend/middleware"
	"github.com/nikhilsshet/wheelz-on-rent/backend/models"
	"github.com/nikhilsshet/wheelz-on-rent/backend/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" || user.PasswordHash == "" || user.Name == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password encryption failed", http.StatusInternalServerError)
		return
	}

	err = config.DB.QueryRow(`INSERT INTO users (name, email, password_hash, role) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Name, user.Email, string(hashedPassword), "customer").Scan(&user.ID)

	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusBadRequest)
		return
	}

	user.PasswordHash = ""
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	err := config.DB.QueryRow(`SELECT id, name, email, password_hash, role 
		FROM users WHERE email = $1`, input.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.PasswordHash))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func Profile(w http.ResponseWriter, r *http.Request) {
	// Extract claims from context
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*utils.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
	})
}
