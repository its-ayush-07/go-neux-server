package middleware

import (
	"net/http"

	token "github.com/its-ayush-07/go-neux-server/token"

	"github.com/gin-gonic/gin"
)

// Middleware function to authorize protected routes
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")

		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		if tokenString == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		claims, msg := token.ValidateToken(tokenString)
		if msg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("username", claims.UserName)
		c.Next()
	}
}
