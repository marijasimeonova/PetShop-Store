package middleware

import (
	"net/http"
	token "pet_shop/tokens"

	"github.com/gin-gonic/gin"
)

// Authentication is a middleware function that checks for a valid token in the request header.
// If no token is provided or the token is invalid, it returns a 500 Internal Server Error.
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the request header
		ClientToken := c.Request.Header.Get("token")
		// No token provided, return error response
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			// Token validation failed, return error response
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		// Set email and uid claims in the context for further processing
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
