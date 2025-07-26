package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikhilsshet/wheelz-on-rent/backend/controllers"
	"github.com/nikhilsshet/wheelz-on-rent/backend/middleware"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/api/register", controllers.Register).Methods("POST")
	r.HandleFunc("/api/login", controllers.Login).Methods("POST")
	
	// Protected route example
	r.Handle("/api/profile", middleware.JWTMiddleware(http.HandlerFunc(controllers.Profile))).Methods("GET")
}

func RegisterVehicleRoutes(router *mux.Router) {
	// Protected route - only authenticated users can add vehicles
	router.HandleFunc("/api/vehicles", middleware.AuthMiddleware(controllers.AddVehicle)).Methods("POST")

	// Public route - anyone can view vehicles
	router.HandleFunc("/api/vehicles", controllers.GetAllVehicles).Methods("GET")
}


