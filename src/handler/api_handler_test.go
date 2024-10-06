package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"src/model/account"
	"src/service/accountService"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestCase struct {
	Name             string
	Method           string
	URL              string
	RequestBody      interface{}
	ExpectedStatus   int
	ExpectedResponse interface{}
	SetupDatabase    func(db *gorm.DB)
	PreQueries       []func(db *gorm.DB)
}

// Unit Test for createAccount Handler with Mock DB
func TestAccountsApi(t *testing.T) {

	// Setup Gin
	gin.SetMode(gin.TestMode)

	testCases := []TestCase{
		{
			Name:             "CreateAccount Success",
			Method:           "POST",
			URL:              "/accounts",
			RequestBody:      `{"document_number": "123456765432"}`,
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: account.Accounts{AccountID: 1, DocumentNumber: "123456765432"},
			SetupDatabase: func(db *gorm.DB) {
				db.AutoMigrate(&account.Accounts{})
			},
		},
		{
			Name:             "CreateAccount Failure - Invalid document number",
			Method:           "POST",
			URL:              "/accounts",
			RequestBody:      `{"document_number": "abcd"}`,
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: gin.H{"error": "Invalid document number"},
			SetupDatabase: func(db *gorm.DB) {
				db.AutoMigrate(&account.Accounts{})
			},
		},
		{
			Name:             "CreateAccount Failed - Account already exists",
			Method:           "POST",
			URL:              "/accounts",
			RequestBody:      `{"document_number": "12345678911"}`,
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: gin.H{"error": "Account already exists"},
			SetupDatabase: func(db *gorm.DB) {
				db.AutoMigrate(&account.Accounts{})
				db.Create(&account.Accounts{AccountID: 1, DocumentNumber: "12345678911"})
			},
		},
		{
			Name:             "GetAccountById Success",
			Method:           "GET",
			URL:              "/accounts/1",
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: account.Accounts{AccountID: 1, DocumentNumber: "12345678900"},
			SetupDatabase: func(db *gorm.DB) {
				db.AutoMigrate(&account.Accounts{})
				db.Create(&account.Accounts{AccountID: 1, DocumentNumber: "12345678900"})
			},
		},
		{
			Name:             "GetAccountById Failure - Account not found",
			Method:           "GET",
			URL:              "/accounts/17",
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: gin.H{"error": "Account not found"},
			SetupDatabase: func(db *gorm.DB) {
				db.AutoMigrate(&account.Accounts{})
				db.Create(&account.Accounts{AccountID: 1, DocumentNumber: "12345678900"})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Initialize in-memory database
			db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

			// Ensure the table exists by calling AutoMigrate
			dbErr := db.AutoMigrate(&account.Accounts{})
			if dbErr != nil {
				t.Fatal("Failed to migrate accounts table:", dbErr)
			}

			if tc.SetupDatabase != nil {
				tc.SetupDatabase(db)
			}

			// Create a Gin router and register handlers
			router := gin.Default()
			router.Use(func(c *gin.Context) {
				db.AutoMigrate(&account.Accounts{})
				c.Set("db", db)
				c.Next()
			})

			router.POST("/accounts", accountService.CreateAccount)
			router.GET("/accounts/:account_id", accountService.GetAccount)

			// Create a request to test the API
			var req *http.Request
			var err error

			if tc.Method == "POST" {
				// Handle POST request
				fmt.Println("tc.ReqBody", tc.RequestBody)
				// body, err := json.Marshal(tc.RequestBody)
				// fmt.Println("err", err)
				// fmt.Println("Body", string(body))
				req, err = http.NewRequest(tc.Method, tc.URL, bytes.NewBufferString(tc.RequestBody.(string)))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tc.Method, tc.URL, nil)
			}
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			router.ServeHTTP(rr, req)

			fmt.Println("rr.Body.String()::::> ", rr.Body.String())

			// Assert the status code
			assert.Equal(t, tc.ExpectedStatus, rr.Code)

			// Assert the response body if the status is OK or Created
			if tc.ExpectedStatus == http.StatusOK || tc.ExpectedStatus == http.StatusCreated {
				var actualResponse account.Accounts
				err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.ExpectedResponse.(account.Accounts).AccountID, actualResponse.AccountID)
				assert.Equal(t, tc.ExpectedResponse.(account.Accounts).DocumentNumber, actualResponse.DocumentNumber)
			} else {
				var actualResponse gin.H
				err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.ExpectedResponse, actualResponse)
			}
		})
	}
}

func TestCreateTransction(t *testing.T) {

	// Setup Gin
	gin.SetMode(gin.TestMode)

	testCases := []TestCase{
		{
			Name:             "CreateTransaction Success",
			Method:           "POST",
			URL:              "/transactions",
			RequestBody:      account.Transaction{AccountID: 1, OperationTypeID: 4, Amount: 100},
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: account.Transaction{TransactionID: 1, AccountID: 1, OperationTypeID: 4, Amount: 100},
			SetupDatabase: func(db *gorm.DB) {
				db.AutoMigrate(&account.Accounts{}, &account.Transaction{}, &account.OperationType{})
				db.Create(&account.Accounts{AccountID: 1, DocumentNumber: "12345678900"})
				db.Create(&account.OperationType{OperationTypeID: 4, Description: "Credit Voucher"})
			},
		},
		{
			Name:             "CreateTransaction Failed - Invalid operation type",
			Method:           "POST",
			URL:              "/transactions",
			RequestBody:      account.Transaction{AccountID: 1, OperationTypeID: 10, Amount: 100},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: gin.H{"error": "Invalid operation type"},
		},
		{
			Name:             "CreateTransaction Failed - Account not found",
			Method:           "POST",
			URL:              "/transactions",
			RequestBody:      account.Transaction{AccountID: 17, OperationTypeID: 4, Amount: 100},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: gin.H{"error": "Account not found"},
			SetupDatabase: func(db *gorm.DB) {
				db.AutoMigrate(&account.Accounts{}, &account.Transaction{}, &account.OperationType{})
				db.Create(&account.OperationType{OperationTypeID: 4, Description: "Credit Voucher"})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Initialize in-memory database
			db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if tc.SetupDatabase != nil {
				tc.SetupDatabase(db)
			}

			// userService := UserService{DB: db}

			// Create a Gin router and register handlers
			router := gin.Default()
			router.Use(func(c *gin.Context) {
				db.AutoMigrate(&account.Accounts{}, &account.Transaction{}, &account.OperationType{})
				c.Set("db", db)
				c.Next()
			})

			router.POST("/transactions", accountService.CreateTransaction)

			// Create a request to test the API
			var req *http.Request
			var err error

			if tc.Method == "POST" {
				// Handle POST request
				body, _ := json.Marshal(tc.RequestBody)
				req, err = http.NewRequest(tc.Method, tc.URL, bytes.NewBuffer(body))
			} else {
				req, err = http.NewRequest(tc.Method, tc.URL, nil)
			}
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			router.ServeHTTP(rr, req)

			// Assert the status code
			assert.Equal(t, tc.ExpectedStatus, rr.Code)

			// Assert the response body if the status is OK or Created
			if tc.ExpectedStatus == http.StatusOK || tc.ExpectedStatus == http.StatusCreated {
				var actualResponse account.Transaction
				err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tc.ExpectedResponse.(account.Transaction).AccountID, actualResponse.AccountID)
				assert.Equal(t, tc.ExpectedResponse.(account.Transaction).TransactionID, actualResponse.TransactionID)
				assert.Equal(t, tc.ExpectedResponse.(account.Transaction).OperationTypeID, actualResponse.OperationTypeID)
				assert.Equal(t, tc.ExpectedResponse.(account.Transaction).Amount, actualResponse.Amount)
			} else {
				var actualResponse gin.H
				err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.ExpectedResponse, actualResponse)
			}
		})
	}
}
