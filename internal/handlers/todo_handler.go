package handlers

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "todo-app/internal/database"
    "todo-app/internal/models"
    "todo-app/internal/s3helper"
)

func GetTodos(c *gin.Context) {
	var todos []models.Todo
	database.DB.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

func CreateTodo(c *gin.Context) {
    var todo models.Todo

    // Bind JSON fields
    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Handle file upload
    file, err := c.FormFile("attachment")
    if err == nil {
        openedFile, _ := file.Open()
        defer openedFile.Close()

        fileURL, uploadErr := s3helper.UploadFile(openedFile, file.Filename)
        if uploadErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
            return
        }
        todo.Attachment = fileURL
    }

    database.DB.Create(&todo)
    c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c *gin.Context) {
    var todo models.Todo
    id := c.Param("id")

    if err := database.DB.First(&todo, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }

    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Handle file upload
    file, err := c.FormFile("attachment")
    if err == nil {
        openedFile, _ := file.Open()
        defer openedFile.Close()

        // Delete the old file if it exists
        if todo.Attachment != "" {
            oldFileName := strings.Split(todo.Attachment, "/")[len(strings.Split(todo.Attachment, "/"))-1]
            s3helper.DeleteFile(oldFileName)
        }

        fileURL, uploadErr := s3helper.UploadFile(openedFile, file.Filename)
        if uploadErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
            return
        }
        todo.Attachment = fileURL
    }

    database.DB.Save(&todo)
    c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
    var todo models.Todo
    id := c.Param("id")

    if err := database.DB.First(&todo, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }

    // Delete attachment if it exists
    if todo.Attachment != "" {
        fileName := strings.Split(todo.Attachment, "/")[len(strings.Split(todo.Attachment, "/"))-1]
        s3helper.DeleteFile(fileName)
    }

    database.DB.Delete(&todo)
    c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
