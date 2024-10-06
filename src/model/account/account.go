package account

// Account model represents the accounts table
type Accounts struct {
	AccountID      uint   `gorm:"primaryKey" json:"account_id"`
	DocumentNumber string `gorm:"not null" json:"document_number"`
}
