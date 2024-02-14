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
	//incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}

func PetRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/petProfile", controllers.AddPet())
	incomingRoutes.GET("/viewpets", controllers.SearchPet())
	//incomingRoutes.GET("/searchpet", controllers.SearchPetByQuery())
	incomingRoutes.DELETE("/pets/:id/adopt", controllers.AdoptPet())
}

func BlogRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/blogPost", controllers.CreateBlogPost())
	incomingRoutes.GET("/viewblogs", controllers.SearchBlog())
	//incomingRoutes.GET("/searchblogs", controllers.SearchBlogByQuery())
}
