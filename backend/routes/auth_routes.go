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
