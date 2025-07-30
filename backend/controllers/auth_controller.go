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

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.PasswordHash == "" || user.Name == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// Save raw password to separate field before hashing
	user.PlainPassword = user.PasswordHash

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Insert into DB
	// db := config.GetDB()
	_, err = config.DB.Exec("INSERT INTO users (name, email, password_hash, plain_password, role) VALUES ($1, $2, $3, $4, $5)",
		user.Name, user.Email, user.PasswordHash, user.PlainPassword, "customer")
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
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
