package handlers
import (
	"encoding/json"
	"todo-app/internal/handlers"
	"todo-app/internal/models"
	"todo-app/internal/database"
	"todo-app/tests/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := models_test.SetupTestDB()
	if err != nil {
		log.Fatalf("Error setting up test database: %v", err)
	}

	code := m.Run()

	// Cleanup
	models_test.TeardownTestDB()

	os.Exit(code)
}

func TestGetTodos(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	// Seed the database with some test data
	todo1 := models.Todo{Title: "Test Todo 1", Description: "Description 1"}
	todo2 := models.Todo{Title: "Test Todo 2", Description: "Description 2"}
	database.DB.Create(&todo1)
	database.DB.Create(&todo2)

	handlers.GetTodos(w, database.DB)
	res := w.Result()

	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, res.StatusCode, http.StatusOK, "API should return 200 status code")

	// Read data from the body and parse the JSON
	var todos []models.Todo
	err := json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(t, err)

	// Check the length of the todos array
	assert.Len(t, todos, 2)

	// Check that the first todo matches the expected data
	assert.Equal(t, todos[0].Title, todo2.Title)
	assert.Equal(t, todos[0].Description, todo2.Description)

	// Check that the second todo matches the expected data
	assert.Equal(t, todos[1].Title, todo1.Title)
	assert.Equal(t, todos[1].Description, todo1.Description)
}