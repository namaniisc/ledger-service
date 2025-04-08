package handlers

import (
	"database/sql"
	"net/http"

	"ledger/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCustomer handles POST /customers.
func CreateCustomer(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name           string  `json:"name"`
			InitialBalance float64 `json:"initial_balance"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		customerID := uuid.New()
		_, err := db.Exec(
			`INSERT INTO customers (customer_id, name, balance) VALUES ($1, $2, $3)`,
			customerID, req.Name, req.InitialBalance,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.Customer{
			CustomerID: customerID,
			Name:       req.Name,
			Balance:    req.InitialBalance,
		})
	}
}

// GetBalance handles GET /customers/:customer_id/balance.
func GetBalance(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID := c.Param("customer_id")
		var balance float64
		err := db.QueryRow(`SELECT balance FROM customers WHERE customer_id = $1`, customerID).Scan(&balance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"customer_id": customerID,
			"balance":     balance,
		})
	}
}
