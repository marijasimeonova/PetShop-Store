package routes

import (
	"pet_shop/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	//incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}

func PetRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/petProfile", controllers.PetViewerAdmin())
	incomingRoutes.GET("/viewpet", controllers.SearchPet())
	//incomingRoutes.GET("/searchpet", controllers.SearchPetByQuery())
	incomingRoutes.DELETE("/pets/:id/adopt", controllers.AdoptPet())
}

func BlogRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/blogPost", controllers.BlogViewerAdmin())
	incomingRoutes.GET("/viewblogs", controllers.SearchBlog())
	//incomingRoutes.GET("/searchblogs", controllers.SearchBlogByQuery())
}
