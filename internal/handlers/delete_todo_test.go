package handlers

import (
	// "bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"gorm.io/gorm"
	"os"
	"fmt"
	"strconv"
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

func TestDeleteTodo(t *testing.T) {
	// Setup Gin router and database
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Use real database for test
	db := setupTestDB()
	database.DB = db

	// Test DELETE endpoint
	router.DELETE("/todos/:id", func(c *gin.Context) {
		handlers.DeleteTodo(c, db)
	})

	t.Run("Delete a Todo successfully", func(t *testing.T) {
		// Insert a sample todo for testing
		todo := models.Todo{
			Title:       "Test Todo",
			Description: "This is a test todo",
			Attachment:  "some/path/to/file.jpg",
		}
		db.Create(&todo)

		// Send DELETE request
		req, _ := http.NewRequest("DELETE", "/todos/"+strconv.Itoa(int(todo.ID)), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Check response status code
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify if the todo is deleted from the database
		var deletedTodo models.Todo
		err := db.First(&deletedTodo, todo.ID).Error
		assert.Error(t, err) // Expected error as the todo should be deleted

		// Cleanup: truncate the table
		truncateTable(db)
	})

	t.Run("Fail when Todo not found", func(t *testing.T) {
		// Send DELETE request for a non-existing todo
		req, _ := http.NewRequest("DELETE", "/todos/non-existing-id", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Expect status 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.Code)

		// Check the error message in the response
		var response map[string]string
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "Todo not found", response["error"])
	})
}
