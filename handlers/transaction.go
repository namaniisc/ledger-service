package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ledger/models"
)

// CreateTransaction handles POST /transactions.
func CreateTransaction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CustomerID          uuid.UUID `json:"customer_id"`
			Type                string    `json:"type"`
			Amount              float64   `json:"amount"`
			ClientTransactionID string    `json:"transaction_id,omitempty"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Parsed transaction request: %+v\n", req)

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start transaction"})
			return
		}
		defer tx.Rollback()

		if req.ClientTransactionID != "" {
			var exists bool
			err = tx.QueryRow(
				`SELECT exists (SELECT 1 FROM transactions WHERE client_transaction_id = $1)`,
				req.ClientTransactionID,
			).Scan(&exists)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "idempotency check failed"})
				return
			}
			if exists {
				c.JSON(http.StatusOK, gin.H{"message": "transaction already processed"})
				return
			}
		}

		var balance float64
		err = tx.QueryRow(
			`SELECT balance FROM customers WHERE customer_id = $1 FOR UPDATE`,
			req.CustomerID,
		).Scan(&balance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "customer not found"})
			return
		}

		if req.Type == "debit" && balance < req.Amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}

		var newBalance float64
		if req.Type == "credit" {
			newBalance = balance + req.Amount
		} else {
			newBalance = balance - req.Amount
		}

		_, err = tx.Exec(
			`UPDATE customers SET balance = $1 WHERE customer_id = $2`,
			newBalance, req.CustomerID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update balance"})
			return
		}

		transactionID := uuid.New()
		_, err = tx.Exec(
			`INSERT INTO transactions (transaction_id, customer_id, type, amount, client_transaction_id, created_at)
             VALUES ($1, $2, $3, $4, $5, $6)`,
			transactionID, req.CustomerID, req.Type, req.Amount, req.ClientTransactionID, time.Now().UTC(),
		)
		if err != nil {
			log.Println("Insert Transaction Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not insert transaction"})
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"transaction_id": transactionID,
			"status":         "success",
			"balance":        newBalance,
		})
	}
}

// GetTransactions handles GET /customers/:customer_id/transactions.
func GetTransactions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID := c.Param("customer_id")
		rows, err := db.Query(
			`SELECT transaction_id, type, amount, created_at FROM transactions WHERE customer_id = $1 ORDER BY created_at DESC`,
			customerID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve transactions"})
			return
		}
		defer rows.Close()

		var transactions []models.Transaction

		for rows.Next() {
			var t models.Transaction
			if err := rows.Scan(&t.TransactionID, &t.Type, &t.Amount, &t.Timestamp); err != nil {
				continue
			}
			transactions = append(transactions, t)
		}
		c.JSON(http.StatusOK, transactions)
	}
}
