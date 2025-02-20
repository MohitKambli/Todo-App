package handlers

import (
	// "bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"os"
	"fmt"
	"testing"
	"todo-app/internal/database"
	"todo-app/internal/models"
	"todo-app/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
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

// Test GetTodoByID Handler
func TestGetTodoByID(t *testing.T) {
	// Initialize test DB
	db := setupTestDB()
	database.DB = db

	// Insert test data
	todo := models.Todo{Title: "Test Todo", Description: "This is a test todo"}
	db.Create(&todo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/todos/:id", func(c *gin.Context) {
		handlers.GetTodoByID(c, db)
	})

	// Test Case: Valid ID
	t.Run("Valid Todo ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/todos/"+strconv.Itoa(int(todo.ID)), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var responseTodo models.Todo
		json.Unmarshal(resp.Body.Bytes(), &responseTodo)
		assert.Equal(t, todo.Title, responseTodo.Title)
		assert.Equal(t, todo.Description, responseTodo.Description)
		truncateTable(db)
	})

	// Test Case: Invalid ID format
	t.Run("Invalid ID Format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/todos/abc", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid ID format")
		truncateTable(db)
	})

	// Test Case: Non-existing ID
	t.Run("Non-existing ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/todos/9999", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), "Todo not found")
		truncateTable(db)
	})
}
