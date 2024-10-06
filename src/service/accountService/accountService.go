package accountService

import (
	"math"
	"net/http"
	"src/model/account"
	"src/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// createAccount handles the POST request to create a new account
func CreateAccount(c *gin.Context) {

	reqAccount := account.Accounts{}

	// Bind JSON request body to the account struct
	if err := c.Bind(&reqAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !util.IsNumeric(reqAccount.DocumentNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document number"})
		return
	}

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection not found"})
		return
	}

	dbOrm, ok := db.(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DB connection"})
		return
	}

	acc := account.Accounts{}
	dbOrm.First(&acc, "document_number = ?", reqAccount.DocumentNumber)
	if acc.AccountID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account already exists"})
		return
	}

	// Save the account to the database
	err := dbOrm.Create(&reqAccount).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Return the created account as a response
	c.JSON(http.StatusOK, reqAccount)
}

// getAccount handles the GET request to retrieve an account by ID
func GetAccount(c *gin.Context) {

	// Get the account ID from the URL parameters
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account := account.Accounts{}

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection not found"})
		return
	}

	dbOrm, ok := db.(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DB connection"})
		return
	}

	// Find the account in the database by the given ID
	if result := dbOrm.First(&account, accountID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
		return
	}

	// Return the account information as a response
	c.JSON(http.StatusOK, account)
}

// createTransaction handles the POST request to create a new transaction
func CreateTransaction(c *gin.Context) {

	transaction := account.Transaction{}

	// Bind JSON request body to the transaction struct
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection not found"})
		return
	}

	dbOrm, ok := db.(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DB connection"})
		return
	}

	// fetch the operation_type
	operationType := account.OperationType{}
	if result := dbOrm.First(&operationType, transaction.OperationTypeID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
		return
	}

	// calclute trasaction amount based on operation type
	amount := transaction.Amount
	transaction.Amount = determineTransactionAmount(transaction.Amount, operationType)

	// Check if the account associated with the transaction exists
	account := account.Accounts{}
	if result := dbOrm.First(&account, transaction.AccountID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
		return
	}

	// Save the transaction to the database
	tx := dbOrm.Create(&transaction)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"}) // Internal server error
		return
	}

	transaction.Amount = amount

	// Return the created transaction as a response
	c.JSON(http.StatusOK, transaction)
}

// the transactionAmount
// - negative for withdrawls and purchases
// + positive for credits
func determineTransactionAmount(transactionAmount float64, operationType account.OperationType) float64 {

	if operationType.TransactionType == 0 {
		transactionAmount = -transactionAmount
	}

	return math.Floor(transactionAmount*100) / 100
}
