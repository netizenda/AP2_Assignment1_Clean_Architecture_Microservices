package main

import (
	"database/sql"
	"log"
	"os"
	"payment-service/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:12345@localhost:5432/paymentdb?sslmode=disable")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := app.NewApp(db)
	port := getEnv("PORT", "8081")

	log.Printf("🚀 Payment Service started on :%s", port)
	if err := server.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
