package main

import (
	"database/sql"
	"log"
	"order-service/internal/app"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/orderdb?sslmode=disable")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	paymentURL := getEnv("PAYMENT_SERVICE_URL", "http://localhost:8081")

	server := app.NewApp(db, paymentURL)
	port := getEnv("PORT", "8080")

	log.Printf("Order Service started on :%s", port)
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
