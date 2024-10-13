package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Middleware triggered")
		tokenString := c.GetHeader("Authorization")
		fmt.Println("Authorization header:", tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims, err := VerifyToken(tokenString)
		if err != nil {
			fmt.Println("Token verification error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			fmt.Println("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("username", username)
		fmt.Println("Username set in context:", username)
		c.Next()
	}
}
