package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
    // Load environment variables from .env
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("Failed to open DB: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }

    fmt.Println("✅ Connected to PostgreSQL database")
}
