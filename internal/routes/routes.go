package routes

import (
	"github.com/gin-gonic/gin"
	"todo-app/internal/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/todos", handlers.GetTodos)
	r.GET("/todos/:id", handlers.GetTodoByID)
	r.POST("/todos", handlers.CreateTodo)
	r.PUT("/todos/:id", handlers.UpdateTodo)
	r.DELETE("/todos/:id", handlers.DeleteTodo)
}
