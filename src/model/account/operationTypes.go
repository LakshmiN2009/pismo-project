package account

// OperationType model represents the operation_types table
type OperationType struct {
	OperationTypeID uint   `gorm:"primaryKey" json:"operation_type_id"`
	Description     string `json:"description"`
	TransactionType int    `json:"transaction_type"`
}
