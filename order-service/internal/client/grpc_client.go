package client

import (
	"context"

	"google.golang.org/grpc"
	pb "order-service/proto/v1"
)

type PaymentGRPCClient struct {
	client pb.PaymentServiceClient
}

func NewPaymentGRPCClient() (*PaymentGRPCClient, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &PaymentGRPCClient{client: pb.NewPaymentServiceClient(conn)}, nil
}

func (c *PaymentGRPCClient) Authorize(ctx context.Context, orderID string, amount int64) (string, string, error) {
	resp, err := c.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
	if err != nil {
		return "", "", err
	}
	return resp.Status, resp.TransactionId, nil
}
