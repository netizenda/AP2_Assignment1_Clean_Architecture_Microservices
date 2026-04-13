package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"payment-service/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:12345@localhost:5432/paymentdb?sslmode=disable")
	grpcAddr := getEnv("GRPC_ADDR", ":50051")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	grpcServer := app.NewGRPCServer(db)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Payment gRPC Server started on %s", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
