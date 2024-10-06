package api

import (
	"src/middleware"
	"src/service/accountService"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes() {
	// Initialize Gin router
	r := gin.Default()

	r.Use(middleware.DbMiddleware())

	// Define the API routes

	// POST /accounts - create a new account
	r.POST("/accounts", accountService.CreateAccount)

	// GET /accounts/:accountId - get account information by account ID
	r.GET("/accounts/:account_id", accountService.GetAccount)

	// POST /transactions - create a new transaction
	r.POST("/transactions", accountService.CreateTransaction)

	// Start the server on port 8080
	r.Run(":8080")
}
