package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"pet_shop/database"
	"pet_shop/models"
	generate "pet_shop/tokens"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var PetCollection *mongo.Collection = database.PetData(database.Client, "Pets")
var BlogCollection *mongo.Collection = database.BlogData(database.Client, "Blogs")
var Validate = validator.New()

var (
	ErrCantFindProduct    = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't find product")
	ErrUserIDIsNotValid   = errors.New("user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add product to cart")
	ErrCantRemoveItem     = errors.New("cannot remove item from cart")
	ErrCantGetItem        = errors.New("cannot get item from cart ")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Bind the JSON data to the user model
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate user input
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		// Check if user already exists by email
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}

		// Check if phone number is already in use
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}

		// Hash the password
		password := HashPassword(*user.Password)
		user.Password = &password

		// Set created_at and updated_at timestamps
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// Generate user ID and tokens
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Order_Status = make([]models.Order, 0)

		// Insert user into the database
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Bind the JSON data to the user model
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// Find user by email
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		// Verify password
		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		// Generate tokens
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()
		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		c.JSON(http.StatusFound, founduser)

	}
}

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a new ObjectID for the product
		products.Product_ID = primitive.NewObjectID()

		// Insert the product into the database
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added our Product!!")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// Retrieve all products from the database
		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
			return
		}

		// Decode the cursor into a list of products
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productlist)

	}
}

func SearchProductByType() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get the type from the URL parameter
		productType := c.Param("type")

		// Create a filter to search for products of the specified type
		filter := bson.M{"type": productType}

		// Search for products matching the filter
		cursor, err := ProductCollection.Find(ctx, filter)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something Went Wrong. Please Try Again Later.")
			return
		}
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid Request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productlist)
	}
}

func AddPet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var pets models.Pet
		defer cancel()
		if err := c.BindJSON(&pets); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pets.Pet_ID = primitive.NewObjectID()
		_, anyerr := PetCollection.InsertOne(ctx, pets)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added our Pet!!")
	}
}

func SearchPet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var petlist []models.Pet
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := PetCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
			return
		}
		err = cursor.All(ctx, &petlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, petlist)

	}
}

// AdoptPet handles the adoption of a pet => actually it is deleting the pet from the database (it is adopted)
func AdoptPet() gin.HandlerFunc {
	return func(c *gin.Context) {
		petID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(petID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
			return
		}

		var pet models.Pet
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Find the pet by ID
		err = PetCollection.FindOne(ctx, bson.M{"_id": objectID, "available": true}).Decode(&pet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find pet"})
			return
		}

		// Mark the pet as adopted
		update := bson.M{"$set": bson.M{"available": false}}
		_, err = PetCollection.UpdateOne(ctx, bson.M{"_id": objectID, "available": true}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark pet as adopted"})
			return
		}

		// Delete the pet from the database
		_, err = PetCollection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pet from database"})
			return
		}

		c.JSON(http.StatusOK, pet)
	}
}

// creating a blog post
func CreateBlogPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var blogs models.BlogPost
		defer cancel()
		if err := c.BindJSON(&blogs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		blogs.Blog_ID = primitive.NewObjectID()
		_, anyerr := BlogCollection.InsertOne(ctx, blogs)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Blog Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully created your Blog!!")
	}
}

func SearchBlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var blogslist []models.BlogPost
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := BlogCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
			return
		}
		err = cursor.All(ctx, &blogslist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, blogslist)

	}
}

func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get productID and userID from the route parameters
		productID := c.Param("productID")
		userID := c.Param("userID")

		// Check if productID or userID is empty
		if productID == "" || userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "productID or userID is empty"})
			return
		}

		// Convert productID to ObjectID
		productObjID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid productID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = AddProductToCart(ctx, ProductCollection, UserCollection, productObjID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully added product to cart"})
	}
}

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	// Retrieve the product details from the products collection
	searchResult := prodCollection.FindOne(ctx, bson.M{"_id": productID})
	if searchResult.Err() != nil {
		return ErrCantFindProduct
	}

	// Decode the product details into a ProductUser object
	var product models.ProductUser
	err := searchResult.Decode(&product)
	if err != nil {
		return ErrCantDecodeProducts
	}

	// Convert userID to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrUserIDIsNotValid
	}

	// Update the user's cart with the product
	filter := bson.M{"_id": userObjID}
	update := bson.M{"$push": bson.M{"usercart": product}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("productID")
		userID := c.Param("userID")
		if productID == "" || userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "productID or userID is empty"})
			return
		}

		productObjID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid productID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = RemoveCartItem(ctx, UserCollection, productObjID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully removed item from cart"})
	}
}

func RemoveCartItem(ctx context.Context, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}
	return nil
}

func ShowItemsFromUserCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get userID from the route parameter
		userID := c.Param("userID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Retrieve the user's cart items
		userCart, err := GetUserCart(ctx, UserCollection, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, userCart)
	}
}

func GetUserCart(ctx context.Context, userCollection *mongo.Collection, userID string) ([]models.ProductUser, error) {
	// Convert userID to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Retrieve the user's document from the collection
	result := userCollection.FindOne(ctx, bson.M{"_id": userObjID})
	if result.Err() != nil {
		return nil, result.Err()
	}

	// Decode the user document into a User object
	var user models.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	// Return the user's cart items
	return user.UserCart, nil
}

func BuyProductFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract userID from the request parameters
		userID := c.Param("userID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call the BuyItemFromCart function to process the purchase
		err := BuyItemFromCart(ctx, UserCollection, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully placed order from cart"})
	}
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	// Convert userID to ObjectID
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	// Initialize variables for storing cart items and order details
	var getCartItems models.User
	var orderCart models.Order
	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Orderered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	// Define aggregation pipeline stages to calculate the total price of items in the cart
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	if err != nil {
		log.Println(err)
		return ErrCantGetItem
	}

	// Decode the aggregated results into a slice of BSON documents
	var getUserCart []bson.M
	if err := currentResults.All(ctx, &getUserCart); err != nil {
		log.Println(err)
		return ErrCantGetItem
	}

	// Calculate the total price from the aggregated results
	var totalPrice int32
	for _, userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}

	// Set the total price in the orderCart struct
	orderCart.Price = int(totalPrice)

	// Define the filter to update the user's document
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	// Update the user document with the order details
	if _, err := userCollection.UpdateMany(ctx, filter, update); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	// Retrieve the user's cart items after the purchase
	if err := userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	// Define the filter to update the user's document again
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}
	if _, err := userCollection.UpdateOne(ctx, filter2, update2); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	// Set the user's cart to empty after the purchase
	userCartEmpty := make([]models.ProductUser, 0)
	filtered := bson.D{primitive.E{Key: "_id", Value: id}}
	updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: userCartEmpty}}}}
	if _, err := userCollection.UpdateOne(ctx, filtered, updated); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get userID and productID from the route parameters
		userID := c.Param("userID")
		productID := c.Param("productID")

		// Check if userID or productID is empty
		if userID == "" || productID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID or productID is empty"})
			return
		}

		// Convert productID to ObjectID
		productObjID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid productID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call the InstantBuyer function to process the instant purchase
		err = InstantBuyer(ctx, ProductCollection, UserCollection, productObjID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully placed instant order"})
	}
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	// Convert userID to ObjectID
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	// Retrieve product details from the products collection
	var productDetails models.ProductUser
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	// Create a new order
	var order models.Order
	order.Order_ID = primitive.NewObjectID()
	order.Orderered_At = time.Now()
	order.Order_Cart = []models.ProductUser{productDetails}
	order.Payment_Method.COD = true
	order.Price = productDetails.Price

	// Push the new order into the user's orders list
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	// Push the product details into the order_list of the new order
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}
