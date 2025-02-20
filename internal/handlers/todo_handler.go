package handlers

import (
	"net/http"
	"strings"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-app/internal/database"
	"todo-app/internal/models"
	"todo-app/internal/s3helper"
)

func GetTodos(c *gin.Context, db *gorm.DB) {
	var todos []models.Todo
	database.DB.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

func GetTodoByID(c *gin.Context, db *gorm.DB) {
	var todo models.Todo
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := database.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func CreateTodo(c *gin.Context, db *gorm.DB) {
	var todo models.Todo
	todo.Title = c.PostForm("title")
	todo.Description = c.PostForm("description")
	var attachmentURLs []string
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}
	files := form.File["files"]
	if len(files) > 0 {
		for _, file := range files {
			openedFile, openErr := file.Open()
			if openErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
				return
			}
			defer openedFile.Close()
			fileURL, uploadErr := s3helper.UploadFile(openedFile, file.Filename)
			if uploadErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
				return
			}
			attachmentURLs = append(attachmentURLs, fileURL)
		}
		todo.Attachment = strings.Join(attachmentURLs, ",")
	}
	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c *gin.Context, db *gorm.DB) {
	var todo models.Todo
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	todo.ID = uint(id)
	if err := database.DB.First(&todo, todo.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	todo.Title = c.PostForm("title")
	todo.Description = c.PostForm("description")
	var attachmentURLs []string
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}
	files := form.File["files"]
	if len(files) > 0 {
		if todo.Attachment != "" {
			oldFiles := strings.Split(todo.Attachment, ",")
			for _, oldFile := range oldFiles {
				fileName := strings.Split(oldFile, "/")[len(strings.Split(oldFile, "/"))-1]
				s3helper.DeleteFile(fileName)
			}
		}
		for _, file := range files {
			openedFile, _ := file.Open()
			defer openedFile.Close()
			fileURL, uploadErr := s3helper.UploadFile(openedFile, file.Filename)
			if uploadErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
				return
			}
			attachmentURLs = append(attachmentURLs, fileURL)
		}
		todo.Attachment = strings.Join(attachmentURLs, ",")
	}
	if err := database.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context, db *gorm.DB) {
	var todo models.Todo
	id := c.Param("id")
	if err := database.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	if todo.Attachment != "" {
		fileName := strings.Split(todo.Attachment, "/")[len(strings.Split(todo.Attachment, "/"))-1]
		s3helper.DeleteFile(fileName)
	}
	database.DB.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
