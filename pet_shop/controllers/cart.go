package controllers

import (
	//"context"
	//"errors"
	//"log"
	//"net/http"
	//"pet_shop/database"

	//"pet_shop/models"
	//"time"

	//"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
	petCollection  *mongo.Collection
	blogCollection *mongo.Collection
}

func NewApplication(prodCollection *mongo.Collection, userCollection *mongo.Collection, petCollection *mongo.Collection, blogCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
		petCollection:  petCollection,
		blogCollection: blogCollection,
	}
}

/*
func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filledCart, err := database.GetUserCart(ctx, app.userCollection, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, filledCart)
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := database.BuyItemsFromCart(ctx, app.userCollection, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully placed order from cart"})
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		productID := c.Param("productID")
		if userID == "" || productID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID or productID is empty"})
			return
		}

		productObjID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid productID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuyProduct(ctx, app.prodCollection, app.userCollection, productObjID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully placed instant order"})
	}
}
*/
