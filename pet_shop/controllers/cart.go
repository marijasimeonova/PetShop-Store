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
