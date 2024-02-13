package main

import (
	"log"
	"os"
	"pet_shop/controllers"
	"pet_shop/database"
	"pet_shop/middleware"
	"pet_shop/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"), database.PetData(database.Client, "Pets"), database.BlogData(database.Client, "Blogs"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.PetRoutes(router)
	routes.BlogRoutes(router)

	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	// Log the error that the router can possibly return.
	log.Fatal(router.Run(":" + port))
}
