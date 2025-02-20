package handlers

import (
	"bytes"
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
	"mime/multipart"
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

func TestUpdateTodo(t *testing.T) {
	// Setup Gin router and database
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Use real database for test
	db := setupTestDB()
	database.DB = db

	// Test PUT endpoint
	router.PUT("/todos/:id", func(c *gin.Context) {
		handlers.UpdateTodo(c, db)
	})

	t.Run("Update Todo successfully", func(t *testing.T) {
		// Insert a sample todo for testing
		todo := models.Todo{
			Title:       "Test Todo",
			Description: "This is a test todo",
			Attachment:  "",
		}
		db.Create(&todo)

		// Prepare the form data to update the todo
		formData := new(bytes.Buffer)
		writer := multipart.NewWriter(formData)
		writer.WriteField("title", "Updated Test Todo")
		writer.WriteField("description", "This is an updated test todo")
		// No files provided, just text fields
		writer.Close()

		// Send PUT request to update the todo
		req, _ := http.NewRequest("PUT", "/todos/"+strconv.Itoa(int(todo.ID)), formData)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Check response status code
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify the updated todo in the database
		var updatedTodo models.Todo
		db.First(&updatedTodo, todo.ID)
		assert.Equal(t, "Updated Test Todo", updatedTodo.Title)
		assert.Equal(t, "This is an updated test todo", updatedTodo.Description)

		// Cleanup: truncate the table
		truncateTable(db)
	})

	t.Run("Fail when Todo not found", func(t *testing.T) {
		// Send PUT request for a non-existing todo
		req, _ := http.NewRequest("PUT", "/todos/999999", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Expect status 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.Code)

		// Check the error message in the response
		var response map[string]string
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "Todo not found", response["error"])
	})

	t.Run("Fail when ID format is invalid", func(t *testing.T) {
		// Send PUT request with an invalid ID format
		req, _ := http.NewRequest("PUT", "/todos/invalid-id", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Expect status 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// Check the error message in the response
		var response map[string]string
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "Invalid ID format", response["error"])
	})
}
