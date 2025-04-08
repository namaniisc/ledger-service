package main

import (
	"log"

	"ledger/database"
	"ledger/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()

	router.POST("/customers", handlers.CreateCustomer(db))
	router.POST("/transactions", handlers.CreateTransaction(db))
	router.GET("/customers/:customer_id/balance", handlers.GetBalance(db))
	router.GET("/customers/:customer_id/transactions", handlers.GetTransactions(db))

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
