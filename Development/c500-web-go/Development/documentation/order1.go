package models

import (
	"database/sql"
	"time"
)

// Order represents a completed transaction.
type Order struct {
	ID              int       `json:"id"`
	StripeSessionID string    `json:"stripe_session_id"`
	ProductID       int       `json:"product_id"`
	AmountTotal     int       `json:"amount_total"` // in cents
	CustomerEmail   string    `json:"customer_email"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// OrderModel holds the database connection pool.
type OrderModel struct {
	DB *sql.DB
}

// Insert adds a new order record to the database.
func (m *OrderModel) Insert(order *Order) error {
	stmt := `
		INSERT INTO orders (stripe_session_id, product_id, amount_total, customer_email, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	// Use QueryRow because we expect one row back (the new ID and timestamp)
	err := m.DB.QueryRow(
		stmt,
		order.StripeSessionID,
		order.ProductID,
		order.AmountTotal,
		order.CustomerEmail,
		order.Status,
	).Scan(&order.ID, &order.CreatedAt) // Update the struct with the new DB ID and time

	if err != nil {
		return err
	}
	return nil
}
