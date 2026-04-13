package grpc

import (
	"order-service/internal/usecase"
)

type OrderStreamServer struct {
	usecase *usecase.OrderUsecase
}

func NewOrderStreamServer(uc *usecase.OrderUsecase) *OrderStreamServer {
	return &OrderStreamServer{usecase: uc}
}

func (s *OrderStreamServer) SubscribeToOrderUpdates(req interface{}, stream interface{}) error {
	return nil
}
