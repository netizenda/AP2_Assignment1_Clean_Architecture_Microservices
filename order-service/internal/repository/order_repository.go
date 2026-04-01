package repository

import (
	"context"
	"database/sql"
	"order-service/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

type postgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) Create(ctx context.Context, o *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, item_name, amount, status, created_at, idempotency_key)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, o.ID, o.CustomerID, o.ItemName, o.Amount, o.Status, o.CreatedAt, o.IdempotencyKey)
	return err
}

func (r *postgresOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, created_at, idempotency_key 
	          FROM orders WHERE id = $1`
	var o domain.Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt, &o.IdempotencyKey)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &o, err
}

func (r *postgresOrderRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Order, error) {
	if key == "" {
		return nil, nil
	}
	query := `SELECT id, customer_id, item_name, amount, status, created_at, idempotency_key 
	          FROM orders WHERE idempotency_key = $1 LIMIT 1`
	var o domain.Order
	err := r.db.QueryRowContext(ctx, query, key).Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt, &o.IdempotencyKey)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &o, err
}

func (r *postgresOrderRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	return err
}
