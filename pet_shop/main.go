package main

import (
	"log"
	"os"
	"pet_shop/controllers"

	//"pet_shop/database"
	"pet_shop/middleware"
	"pet_shop/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	//app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"), database.PetData(database.Client, "Pets"), database.BlogData(database.Client, "Blogs"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.CartRoutes(router)
	routes.PetRoutes(router)
	routes.BlogRoutes(router)

	router.Use(middleware.Authentication())

	router.POST("/addaddress/:userID", controllers.AddAddress())
	router.PUT("/edithomeaddress/:userID", controllers.EditHomeAddress())
	router.PUT("/editworkaddress/:userID", controllers.EditWorkAddress())
	router.DELETE("/deleteaddresses/:userID", controllers.DeleteAddress())
	//router.POST("/cartcheckout/:userID", app.BuyFromCart())
	//router.POST("/instantbuy/:userID/:productID", app.InstantBuy())

	// Log the error that the router can possibly return.
	log.Fatal(router.Run(":" + port))
}
