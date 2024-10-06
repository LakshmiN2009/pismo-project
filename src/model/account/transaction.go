package account

import "time"

// Transaction model represents the transactions table
type Transaction struct {
	TransactionID   uint      `gorm:"primaryKey" json:"transaction_id"`
	AccountID       uint      `json:"account_id"`
	OperationTypeID uint      `json:"operation_type_id"`
	Amount          float64   `gorm:"not null" json:"amount"`
	EventDate       time.Time `gorm:"autoCreateTime" json:"event_date"`
}
