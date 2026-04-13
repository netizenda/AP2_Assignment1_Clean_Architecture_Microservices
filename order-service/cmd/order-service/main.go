package main

import (
	"database/sql"
	"log"
	"os"

	"order-service/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:12345@localhost:5432/orderdb?sslmode=disable")
	restAddr := getEnv("REST_ADDR", ":8080")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app.StartServers(db, restAddr)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
