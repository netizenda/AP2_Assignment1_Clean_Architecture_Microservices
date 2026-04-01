package usecase

import (
	"context"
	"errors"
	"order-service/internal/domain"
	"order-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidAmount      = errors.New("amount must be > 0")
	ErrOrderNotFound      = errors.New("order not found")
	ErrCannotCancel       = errors.New("only pending orders can be cancelled")
	ErrPaymentUnavailable = errors.New("payment service unavailable")
)

type ExternalPaymentClient interface {
	Authorize(ctx context.Context, orderID string, amount int64) (status string, transactionID string, err error)
}

type OrderUsecase struct {
	repo          repository.OrderRepository
	paymentClient ExternalPaymentClient
}

func NewOrderUsecase(repo repository.OrderRepository, pc ExternalPaymentClient) *OrderUsecase {
	return &OrderUsecase{repo: repo, paymentClient: pc}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, customerID, itemName string, amount int64, idempotencyKey string) (*domain.Order, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Idempotency (bonus)
	if idempotencyKey != "" {
		if existing, _ := u.repo.GetByIdempotencyKey(ctx, idempotencyKey); existing != nil {
			return existing, nil
		}
	}

	order := &domain.Order{
		ID:             uuid.New().String(),
		CustomerID:     customerID,
		ItemName:       itemName,
		Amount:         amount,
		Status:         "Pending",
		CreatedAt:      time.Now(),
		IdempotencyKey: idempotencyKey,
	}

	if err := u.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Call Payment Service
	status, _, err := u.paymentClient.Authorize(ctx, order.ID, amount)
	if err != nil {
		_ = u.repo.UpdateStatus(ctx, order.ID, "Failed")
		if ctx.Err() != nil {
			return nil, ErrPaymentUnavailable
		}
		return nil, err
	}

	if status == "Authorized" {
		_ = u.repo.UpdateStatus(ctx, order.ID, "Paid")
		order.Status = "Paid"
	} else {
		_ = u.repo.UpdateStatus(ctx, order.ID, "Failed")
		order.Status = "Failed"
	}
	return order, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *OrderUsecase) CancelOrder(ctx context.Context, id string) error {
	order, err := u.repo.GetByID(ctx, id)
	if err != nil || order == nil {
		return ErrOrderNotFound
	}
	if order.Status != "Pending" {
		return ErrCannotCancel
	}
	return u.repo.UpdateStatus(ctx, id, "Cancelled")
}
