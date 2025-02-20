package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"gorm.io/gorm"
	"os"
	"fmt"
	"todo-app/internal/database"
	"todo-app/internal/models"
	"todo-app/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
)

// Truncate the table
func truncateTable(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE todos RESTART IDENTITY CASCADE;")
}

// Setup Test Database
func setupTestDB() *gorm.DB {
	godotenv.Load("../../.env")
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	// Auto Migrate
	db.AutoMigrate(&models.Todo{})
	return db
}

// Test CreateTodo API
func TestCreateTodo(t *testing.T) {
	db := setupTestDB()
	database.DB = db

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/todos", func(c *gin.Context) {
		handlers.CreateTodo(c, db)
	})

	t.Run("Fail when no form data is provided", func(t *testing.T) {
		// Send an empty request with no form data
		req, _ := http.NewRequest("POST", "/todos", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Perform request
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Expect a 400 Bad Request due to missing form data
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// Truncate table after test
		truncateTable(db)
	})
}

