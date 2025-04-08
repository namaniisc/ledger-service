package models

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer account.
type Customer struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Name       string    `json:"name"`
	Balance    float64   `json:"balance"`
}

// Transaction represents a transaction entry.
type Transaction struct {
	TransactionID       uuid.UUID `json:"transaction_id"`
	CustomerID          uuid.UUID `json:"customer_id"`
	Type                string    `json:"type"` // "credit" or "debit"
	Amount              float64   `json:"amount"`
	ClientTransactionID string    `json:"client_transaction_id,omitempty"`
	Timestamp           time.Time `json:"timestamp"`
}
