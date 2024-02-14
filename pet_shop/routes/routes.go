package routes

import (
	"pet_shop/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/addproduct", controllers.AddProduct())
	incomingRoutes.GET("/users/viewproducts", controllers.SearchProduct())
	incomingRoutes.GET("/users/viewproducts/:type", controllers.SearchProductByType())
}

func CartRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/addtocart/:productID/:userID", controllers.AddToCart())
	incomingRoutes.DELETE("/removefromcart/:productID/:userID", controllers.RemoveFromCart())
	incomingRoutes.GET("/usercart/:userID", controllers.ShowItemsFromUserCart())
	incomingRoutes.POST("/buyfromcart/:userID", controllers.BuyProductFromCart())
	incomingRoutes.POST("/instantbuy/:userID/:productID", controllers.InstantBuy())
}

func PetRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/petProfile", controllers.AddPet())
	incomingRoutes.GET("/viewpets", controllers.SearchPet())
	incomingRoutes.DELETE("/pets/:id/adopt", controllers.AdoptPet())
}

func BlogRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/blogPost", controllers.CreateBlogPost())
	incomingRoutes.GET("/viewblogs", controllers.SearchBlog())
}
