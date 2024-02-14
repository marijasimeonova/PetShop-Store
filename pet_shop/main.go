package main

import (
	"log"
	"os"
	"pet_shop/middleware"
	"pet_shop/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.CartRoutes(router)
	routes.PetRoutes(router)
	routes.BlogRoutes(router)
	router.Use(middleware.Authentication())

	log.Fatal(router.Run(":" + port))
}
