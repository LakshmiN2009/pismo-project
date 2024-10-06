package drivers

import (
	"fmt"
	"os"
	"src/model/account"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database initialization
func InitDB() *gorm.DB {
	var err error
	fmt.Println("Connecting to the database", os.Getenv("DB_HOST"))
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// AutoMigrate will create the tables and update the schema if necessary
	db.AutoMigrate(&account.Accounts{}, &account.OperationType{}, &account.Transaction{})

	return db
}
