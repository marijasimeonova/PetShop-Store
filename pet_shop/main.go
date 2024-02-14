package main

import (
	"log"
	"os"
	"pet_shop/middleware"
	"pet_shop/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get the port from the environment variable, default to 8000 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Initialize a new Gin router
	router := gin.New()

	// Use Gin's default logger middleware
	router.Use(gin.Logger())

	// Define routes
	routes.UserRoutes(router)
	routes.CartRoutes(router)
	routes.PetRoutes(router)
	routes.BlogRoutes(router)

	// Apply authentication middleware to all routes
	router.Use(middleware.Authentication())

	// Start the server
	log.Fatal(router.Run(":" + port))
}
