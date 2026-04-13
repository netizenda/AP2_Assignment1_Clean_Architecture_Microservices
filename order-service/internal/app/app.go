package app

import (
	"database/sql"
	"log"

	"order-service/internal/client"
	"order-service/internal/repository"
	httpTransport "order-service/internal/transport/http"
	"order-service/internal/usecase"
)

func StartServers(db *sql.DB, restAddr string) {
	repo := repository.NewPostgresOrderRepository(db)

	paymentClient, err := client.NewPaymentGRPCClient()
	if err != nil {
		log.Fatalf("failed to create payment client: %v", err)
	}

	uc := usecase.NewOrderUsecase(repo, paymentClient)

	httpServer := httpTransport.NewServer(uc)
	log.Printf("sOrder REST Server started on %s", restAddr)
	httpServer.Run(restAddr)
}
