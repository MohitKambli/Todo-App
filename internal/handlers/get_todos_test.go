package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
	"fmt"
	"todo-app/internal/database"
	"todo-app/internal/models"
	"todo-app/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

// Test GetTodos Handler
func TestGetTodos(t *testing.T) {
	// Initialize test DB
	db := setupTestDB()
	database.DB = db

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/todos", func(c *gin.Context) {
		handlers.GetTodos(c, db)
	})

	// Seed test data
	testTodos := []models.Todo{
		{Title: "Todo 1", Description: "First test todo"},
		{Title: "Todo 2", Description: "Second test todo"},
	}
	db.Create(&testTodos)

	// Test Case: Get all todos (non-empty DB)
	t.Run("Get All Todos (Non-Empty)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/todos", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var responseTodos []models.Todo
		json.Unmarshal(resp.Body.Bytes(), &responseTodos)

		assert.Len(t, responseTodos, 2)
		assert.Equal(t, testTodos[0].Title, responseTodos[0].Title)
		assert.Equal(t, testTodos[1].Title, responseTodos[1].Title)
		truncateTable(db)
	})

	// Test Case: Get todos when DB is empty
	t.Run("Get Todos (Empty DB)", func(t *testing.T) {
		db.Exec("DELETE FROM todos") // Clear table

		req, _ := http.NewRequest("GET", "/todos", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var responseTodos []models.Todo
		json.Unmarshal(resp.Body.Bytes(), &responseTodos)

		assert.Len(t, responseTodos, 0) // Expect empty array
		truncateTable(db)
	})
}
