package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-app/internal/handlers"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/todos", func(c *gin.Context) { handlers.GetTodos(c, db) })
	r.GET("/todos/:id", func(c *gin.Context) { handlers.GetTodoByID(c, db) })
	r.POST("/todos", func(c *gin.Context) { handlers.CreateTodo(c, db) })
	r.PUT("/todos/:id", func(c *gin.Context) { handlers.UpdateTodo(c, db) })
	r.DELETE("/todos/:id", func(c *gin.Context) { handlers.DeleteTodo(c, db) })
}
