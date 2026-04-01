package domain

import "time"

type Order struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customer_id"`
	ItemName       string    `json:"item_name"`
	Amount         int64     `json:"amount"`
	Status         string    `json:"status"` // Pending, Paid, Failed, Cancelled
	CreatedAt      time.Time `json:"created_at"`
	IdempotencyKey string    `json:"-"` // для bonus
}
