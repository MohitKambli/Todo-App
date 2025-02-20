package models_test

import (
	"todo-app/internal/database"
	"os"
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

func SetupTestDB() error {
	if err := godotenv.Load("../../.env"); err != nil {
        log.Println("Warning: No .env file found")
    }
	connStr := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return err
	}
	database.DB = db
	// Run migrations
	database.DBMigrate()

	return nil
}

func TeardownTestDB() {
	var tableNames []string

	// Query to get all table names in the current database
	database.DB.Raw("SHOW TABLES").Scan(&tableNames)

	// Iterate over each table name and drop it
	for _, tableName := range tableNames {
		database.DB.Exec("DROP TABLE IF EXISTS " + tableName)
	}
}