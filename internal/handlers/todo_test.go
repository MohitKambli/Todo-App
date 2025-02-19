package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"todo-app/internal/database"
	"todo-app/internal/models"
	"todo-app/internal/routes"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	RegisterRoutes(r)
	return r
}

func TestCreateTodo(t *testing.T) {
	router := setupRouter()

	// Prepare form data for file upload
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add title and description
	_ = writer.WriteField("title", "Test Todo")
	_ = writer.WriteField("description", "This is a test todo.")

	// Add file (mock file upload)
	filePath := "test_file.txt"
	file, _ := os.Create(filePath)
	defer os.Remove(filePath) // Cleanup test file
	file.WriteString("This is a test file.") // Mock file content
	file.Close()

	part, _ := writer.CreateFormFile("file", filepath.Base(filePath))
	mockFile, _ := os.Open(filePath)
	io.Copy(part, mockFile)
	mockFile.Close()
	writer.Close()

	// Perform the request
	req, _ := http.NewRequest("POST", "/todos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate response
	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.Todo
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Test Todo", response.Title)
}

func TestGetTodos(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSingleTodo(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/todos/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateTodo(t *testing.T) {
	router := setupRouter()

	// JSON request body
	todo := map[string]string{"title": "Updated Title", "description": "Updated description"}
	jsonBody, _ := json.Marshal(todo)
	req, _ := http.NewRequest("PUT", "/todos/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTodo(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("DELETE", "/todos/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
