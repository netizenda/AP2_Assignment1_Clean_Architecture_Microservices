package usecase

import (
	"context"
	"payment-service/internal/domain"
	"payment-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type PaymentUsecase struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: repo}
}

func (u *PaymentUsecase) CreatePayment(ctx context.Context, orderID string, amount int64) (*domain.Payment, error) {
	status := "Authorized"
	if amount > 100000 {
		status = "Declined"
	}

	payment := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: "txn_" + uuid.New().String()[:8],
		Amount:        amount,
		Status:        status,
		CreatedAt:     time.Now(),
	}

	if err := u.repo.Create(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (u *PaymentUsecase) GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return u.repo.GetByOrderID(ctx, orderID)
}
