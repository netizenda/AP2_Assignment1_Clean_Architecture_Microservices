package repository

import (
	"context"
	"database/sql"
	"payment-service/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *domain.Payment) error
	GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
}

type postgresPaymentRepository struct {
	db *sql.DB
}

func NewPostgresPaymentRepository(db *sql.DB) PaymentRepository {
	return &postgresPaymentRepository{db: db}
}

func (r *postgresPaymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	query := `INSERT INTO payments (id, order_id, transaction_id, amount, status, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status, p.CreatedAt)
	return err
}

func (r *postgresPaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	query := `SELECT id, order_id, transaction_id, amount, status, created_at 
	          FROM payments WHERE order_id = $1 LIMIT 1`
	var p domain.Payment
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}
