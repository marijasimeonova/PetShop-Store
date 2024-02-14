## **“Online PetShop Store”**

<img src="pets.jpg"  alt="drawing" width="500" height="300"/>

# project structure

- **pet_shop** 📁
  - controllers 📁
    - cart.go📝
    - controllers.go📝
  - database📁
    - databasetup.go📝
  - middleware📁
    - middleware.go📝
  - models📁
    - models.go📝
  - routes📁
    - routes.go📝
  - tokens📁
    - tokengen.go📝
  - docker-compose.yaml
  - go.sum 📝
  - go.mod 📝
  - main.go📝


```bash
# Start the mongodb container for local development
Ensure the following software is installed on your machine:
The Go Programming Language
Docker Desktop
Postman
MongoDB
Studio 3T

```

```bash
- Start the MongoDB container for local development
docker-compose up -d
- Run the application 
go run ..path/pet_shop/main.go
- Once the application is running, you can access it through your web browser using the following URL:
http://localhost:8000/
- Send requests in Postman to test the app
```

## API FUNCTIONALITIES

- Register
- Log in
- Add a product to DB
- View products that are added
- Add product to cart
- Remove product from cart
- View products in cart
- Buy products from cart
- Instant buy a product
- Add pet for adoption
- View pets available for adoption
- Adopt a pet
- Create and publish a blog post
- View and read all blogs


## models.go

This part defines us how should our database looks like.

- We have to register a user first, having an unique id and token, and after the register user can log in with email and password.

```go
type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       *string            `json:"last_name"  validate:"required,min=2,max=30"`
	Password        *string            `json:"password"   validate:"required,min=6"`
	Email           *string            `json:"email"      validate:"email,required"`
	Phone           *string            `json:"phone"      validate:"required"`
	Token           *string            `json:"token"`
	Refresh_Token   *string            `josn:"refresh_token"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updtaed_at"`
	User_ID         string             `json:"user_id"`
	UserCart        []ProductUser      `json:"usercart" bson:"usercart"`
	Address_Details []Address          `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
}

```

- We have to define a product, having a unique id, name and price. The product can be toy, food or accessories.

```go
type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Rating       *uint8             `json:"rating"`
	Image        *string            `json:"image"`
	Type         string             `json:"type"`
}
```

- We have to define an slice of array products where a user can store individual products

```go
type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price"  bson:"price"`
	Rating       *uint              `json:"rating" bson:"rating"`
	Image        *string            `json:"image"  bson:"image"`
	Type         string             `json:"type" bson:"type"`
}
```

- Order struct:

```go
type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_list"  bson:"order_list"`
	Orderered_At   time.Time          `json:"ordered_on"  bson:"ordered_on"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       *int               `json:"discount"    bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}
```

The Payment struct:

```go
type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod"     bson:"cod"`
}
```

We can define a pet for adoption, having unique pet id.

```go
type Pet struct {
	Pet_ID      primitive.ObjectID `bson:"_id"`
	Pet_Name    *string            `json:"pet_name"`
	Breed       string             `json:"breed"`
	Age         int                `json:"age"`
	Temperament string             `json:"temperament"`
	Description string             `json:"description"`
	Available   bool               `json:"available"`
}
```

- Blog struct:

```go
type BlogPost struct {
	Blog_ID     primitive.ObjectID `bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Content     string             `json:"content" bson:"content"`
	Author      string             `json:"author" bson:"author"`
	PublishDate time.Time          `json:"publish_date" bson:"publish_date"`
	Tags        []string           `json:"tags" bson:"tags"`
}
```

# Requests made in Postman

- **SIGNUP FUNCTION API CALL (POST REQUEST)**

http://localhost:8000/users/signup

```json
Body:
{
  "first_name": "Lina",
  "last_name": "Wood",
  "email": "lina@example.com",
  "password": "lina12",
  "phone": "+1666426655"
}


```

Answ:	"Successfully Signed Up!!"

- **LOGIN FUNCTION API CALL (POST REQUEST)**

  http://localhost:8000/users/login


```json
body:
{
  "email": "lina@example.com",
  "password": "lina12",
}

```

Answ:
```json
{
    "_id": "65cbf3f9a31529feb8412f07",
    "first_name": "Lina",
    "last_name": "Wood",
    "password": "$2a$14$XFfkTALhJg1HKPrBTlkum.Gb1.Fm2.x1IzC5GNvSpaENgyvx3Rq7S",
    "email": "lina@example.com",
    "phone": "+1666426655",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6Im1hcmVAZXhhbXBsZS5jb20iLCJGaXJzdF9OYW1lIjoiTWFyaWphIiwiTGFzdF9OYW1lIjoiU2ltZW9ub3ZhIiwiVWlkIjoiNjVjYmYzZjlhMzE1MjlmZWI4NDEyZjA3IiwiZXhwIjoxNzA3OTUxNDgxfQ.NN0CS85idQZrCZBojBAWwseIzfujzm1wY-YtVEWMHa4",
    "Refresh_Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0X05hbWUiOiIiLCJMYXN0X05hbWUiOiIiLCJVaWQiOiIiLCJleHAiOjE3MDg0Njk4ODF9.Guh0Gu9TVytUhnSG0VUtEr4uTtfB-nhOyWSpImF67XI",
    "created_at": "2024-02-13T22:58:01Z",
    "updtaed_at": "2024-02-13T22:58:01Z",
    "user_id": "65cbf3f9a31529feb8412f07",
    "usercart": [],
    "orders": []
}
```


- **ADD PRODUCT FUNCTION (POST REQUEST)**

  http://localhost:8000/addproduct


```json
Authorization:
Type: Bearer Token
Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0X05hbWUiOiIiLCJMYXN0X05hbWUiOiIiLCJVaWQiOiIiLCJleHAiOjE3MDg0NzA1OTZ9.mryH_dsRvpF0MZAB4J_iRR_EZGFvN3jf_VdNfUA4_2c

body:
{
  "product_name": "bone",
  "price": 279,
  "rating": 10,
  "image": "bone.jpg"
  "type": "food"
}
```

Answ:
```json
{
    "Successfully added our Product!!"

}
```

  
- **VIEW PRODUCTS FUNCTION (GET REQUEST)**

  http://localhost:8000/users/viewproducts


```json
Auth:
	Bearer token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0X05hbWUiOiIiLCJMYXN0X05hbWUiOiIiLCJVaWQiOiIiLCJleHAiOjE3MDg0Njk4ODF9.Guh0Gu9TVytUhnSG0VUtEr4uTtfB-nhOyWSpImF67XI
```

Answ:
```json
{
[
    {
        "Product_ID": "65cbfcdda31529feb8412f0a",
        "product_name": "dog toy",
        "price": 90,
        "rating": 8,
        "image": "toy.jpg"
	"type": "toy"
    },
    {
        "Product_ID": "65cbfe12a31529feb8412f0b",
        "product_name": "daily rewards - food sticks",
        "price": 150,
        "rating": 10,
        "image": "sticks.jpg"
"type": "food"
    },
    {
        "Product_ID": "65cbfec79a05d86baff6c4d9",
        "product_name": " sticks",
        "price": 150,
        "rating": 10,
        "image": "ss.jpg"
	"type": "food"
    }
]

}
```

- **SEARCH PRODUCT BY TYPE FUNCTION (GET REQUEST)**

  http://localhost:8000/users/viewproducts/:type


```json
Auth:
	Bearer token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0X05hbWUiOiIiLCJMYXN0X05hbWUiOiIiLCJVaWQiOiIiLCJleHAiOjE3MDg0Njk4ODF9.Guh0Gu9TVytUhnSG0VUtEr4uTtfB-nhOyWSpImF67XI


```

Answ:
```json
[
    {
        "Product_ID": "65cbfcdda31529feb8412f0a",
        "product_name": "dog toy",
        "price": 90,
        "rating": 8,
        "image": "toy.jpg"
	  "type": "toy"
    }
]
```

- **ADD PRODUCT TO CART FUNCTION (POST REQUEST)**

  http://localhost:8000/addtocart/:productID/:userID
  Exp : http://localhost:8000/addtocart/65cbfcdda31529feb8412f0a/65cbf3f9a31529feb8412f07


```json
Auth:
	Bearer token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0X05hbWUiOiIiLCJMYXN0X05hbWUiOiIiLCJVaWQiOiIiLCJleHAiOjE3MDg0Njk4ODF9.Guh0Gu9TVytUhnSG0VUtEr4uTtfB-nhOyWSpImF67XI
```

Answ:
```json
{
    "message": "Successfully added product to cart"
}
```


- **REMOVE PRODUCT FROM CART FUNCTION (DELETE REQUEST)**

  DELETE: http://localhost:8000/removefromcart/:productID/:userID

  Exp: http://localhost:8000/removefromcart/65cceb2c7b24863b8f0c3281/65cbf3f9a31529feb8412f07

```json
Auth:
	Bearer token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0X05hbWUiOiIiLCJMYXN0X05hbWUiOiIiLCJVaWQiOiIiLCJleHAiOjE3MDg0Njk4ODF9.Guh0Gu9TVytUhnSG0VUtEr4uTtfB-nhOyWSpImF67XI
CREATING PET PROFILE
```

Answ:
```json
{
    "message": "Successfully removed item from cart"
}
```