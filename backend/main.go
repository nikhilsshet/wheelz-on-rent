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
	routes.AuthRoutes(r) // Mount auth routes
	routes.RegisterVehicleRoutes(r) // Mount vehicle routes
	routes.RegisterBookingRoutes(r)



	fmt.Println("🚀 Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
