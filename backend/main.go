package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikhilsshet/wheelz-on-rent/backend/config"
	"github.com/nikhilsshet/wheelz-on-rent/backend/routes"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()

	r.Use(enableCORS)

	routes.AuthRoutes(r)            // Mount auth routes
	routes.RegisterVehicleRoutes(r) // Mount vehicle routes
	routes.RegisterBookingRoutes(r)

	fmt.Println("🚀 Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
