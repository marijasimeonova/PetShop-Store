package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBSet initializes and returns a MongoDB client.
func DBSet() *mongo.Client {
	// Set up MongoDB connection options
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://development:testpassword@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Ping MongoDB to check connectivity
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil
	}

	fmt.Println("Successfully Connected to the mongodb")
	return client
}

// Client holds the MongoDB client instance.
var Client *mongo.Client = DBSet()

// UserData returns a MongoDB collection for user-related data.
func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("PetStore").Collection(CollectionName)
	return collection
}

// ProductData returns a MongoDB collection for product-related data.
func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("PetStore").Collection(CollectionName)
	return productcollection
}

// PetData returns a MongoDB collection for pet-related data.
func PetData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var petcollection *mongo.Collection = client.Database("PetStore").Collection(CollectionName)
	return petcollection
}

// BlogData returns a MongoDB collection for blog-related data.
func BlogData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var blogcollection *mongo.Collection = client.Database("PetStore").Collection(CollectionName)
	return blogcollection
}
