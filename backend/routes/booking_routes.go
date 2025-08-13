package routes

import (
	"github.com/gorilla/mux"
	"github.com/nikhilsshet/wheelz-on-rent/backend/controllers"
	"github.com/nikhilsshet/wheelz-on-rent/backend/middleware"
)

func RegisterBookingRoutes(router *mux.Router) {
	router.HandleFunc("/api/bookings", middleware.AuthMiddleware(controllers.CreateBooking)).Methods("POST")
	router.HandleFunc("/api/bookings", middleware.AuthMiddleware(controllers.GetMyBookings)).Methods("GET")
	router.HandleFunc("/api/bookings/{id}/cancel", middleware.AuthMiddleware(controllers.CancelBooking)).Methods("PATCH")
	router.HandleFunc("/api/bookings/all", middleware.AuthMiddleware(controllers.GetAllBookings)).Methods("GET")
	router.HandleFunc("/api/mybookings", middleware.AuthMiddleware(controllers.GetMyBookings)).Methods("GET", "OPTIONS")
}
