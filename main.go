package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/kewlspy/web-backend/models"
	"github.com/kewlspy/web-backend/routes"
)

func main() {
	// Use a single router instance
	router := gin.Default()

	// Apply CORS to this router
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Connect to DB
	models.ConnectDB()

	// Register routes on the same router
	routes.RegisterRoutes(router)

	// Start the server
	router.Run(":8080")
}
