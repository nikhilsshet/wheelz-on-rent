package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikhilsshet/wheelz-on-rent/backend/config"
)

func main() {
    config.ConnectDB() // Connect to DB

    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Wheelz On Rent Backend")
    })

    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", r)
}
