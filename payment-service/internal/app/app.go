package app

import (
	"database/sql"
	"log"

	"payment-service/internal/repository"
	grpcTransport "payment-service/internal/transport/grpc"
	"payment-service/internal/usecase"

	"google.golang.org/grpc"
	pb "payment-service/proto/v1"
)

func NewGRPCServer(db *sql.DB) *grpc.Server {
	repo := repository.NewPostgresPaymentRepository(db)
	uc := usecase.NewPaymentUsecase(repo)

	s := grpc.NewServer()

	pb.RegisterPaymentServiceServer(s, grpcTransport.NewPaymentServer(uc))

	log.Println("gRPC Payment Server initialized successfully")
	return s
}
