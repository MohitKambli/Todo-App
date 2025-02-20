package database

import (
	"fmt"
	"log"
	"os"
	"todo-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func InitDatabase() {
	godotenv.Load()
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }

    fmt.Println("Connected to PostgreSQL successfully!")

    // Migrate the models
    if err := DB.AutoMigrate(&models.Todo{}); err != nil {
        log.Fatal("AutoMigrate error:", err)
    }
}

func GetDB() *gorm.DB {
	return DB
}

func DBMigrate() {
	DB.AutoMigrate(&models.Todo{})
}