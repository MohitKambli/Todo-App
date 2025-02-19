package handlers

import (
    "net/http"
    "strings"
    "strconv"

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

func GetTodoByID(c *gin.Context) {
    var todo models.Todo
    // Extract ID from URL parameter
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32) // Convert ID to uint
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    // Fetch the todo item from the database by ID
    if err := database.DB.First(&todo, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }
    // Respond with the todo item
    c.JSON(http.StatusOK, todo)
}


func CreateTodo(c *gin.Context) {
    var todo models.Todo
    // Extract JSON fields manually (since ShouldBindJSON doesn't work with multipart)
    idStr := c.PostForm("id")
    id, err := strconv.ParseUint(idStr, 10, 32) // Parsing the string into a uint
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    todo.ID = uint(id)
    todo.Title = c.PostForm("title")
    todo.Description = c.PostForm("description")

    // Handle multiple file uploads (multipart/form-data)
    var attachmentURLs []string
    form, err := c.MultipartForm() // Get the form, handle error
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
        return
    }

    files := form.File["files"] // Access the array of files using the field name "files"
    if len(files) > 0 {
        for _, file := range files {
            openedFile, openErr := file.Open()
            if openErr != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
                return
            }
            defer openedFile.Close()
            // Upload each file to S3
            fileURL, uploadErr := s3helper.UploadFile(openedFile, file.Filename)
            if uploadErr != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
                return
            }
            attachmentURLs = append(attachmentURLs, fileURL) // Store file URL
        }
        todo.Attachment = strings.Join(attachmentURLs, ",") // Save all URLs as comma-separated string
    }

    // Save the todo item in the database
    if err := database.DB.Create(&todo).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
        return
    }
    // Respond with the created todo
    c.JSON(http.StatusCreated, todo)
}


func UpdateTodo(c *gin.Context) {
    var todo models.Todo
    // Extract ID from URL parameter and parse it to uint
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32) // Parse the string into a uint
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    todo.ID = uint(id)

    // Check if the todo item exists in the database
    if err := database.DB.First(&todo, todo.ID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }

    // Update fields from the form (multipart/form-data)
    todo.Title = c.PostForm("title")
    todo.Description = c.PostForm("description")

    // Handle multiple file uploads (multipart/form-data)
    var attachmentURLs []string
    form, err := c.MultipartForm() // Get the form, handle error
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
        return
    }

    files := form.File["files"] // Access the array of files using the field name "files"
    if len(files) > 0 {
        // Delete old files (if any)
        if todo.Attachment != "" {
            oldFiles := strings.Split(todo.Attachment, ",")
            for _, oldFile := range oldFiles {
                fileName := strings.Split(oldFile, "/")[len(strings.Split(oldFile, "/"))-1]
                s3helper.DeleteFile(fileName)
            }
        }

        // Upload each new file to S3
        for _, file := range files {
            openedFile, _ := file.Open()
            defer openedFile.Close()

            // Upload each file to S3
            fileURL, uploadErr := s3helper.UploadFile(openedFile, file.Filename)
            if uploadErr != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
                return
            }
            attachmentURLs = append(attachmentURLs, fileURL) // Store file URL
        }
        todo.Attachment = strings.Join(attachmentURLs, ",") // Save all URLs as comma-separated string
    }

    // Save the updated todo item in the database
    if err := database.DB.Save(&todo).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
        return
    }

    // Respond with the updated todo item
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
