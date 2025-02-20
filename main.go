package main

import (
	"fmt"
	"log"
	"todo-app/internal/database"
	"todo-app/internal/routes"
	"todo-app/internal/s3helper"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	database.InitDatabase()
	s3helper.InitS3()

	// Get the database instance
	db := database.GetDB()

	// Initialize Gin router
	r := gin.Default()

	// Register the routes
	routes.RegisterRoutes(r, db)

	// Start the server
	port := ":8080"
	fmt.Println("Server is running on port", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
