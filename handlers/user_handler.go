package handlers

import (
	"blogapi-naveen/models"
	"blogapi-naveen/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request"})
		return
	}

	err = user.UserSave()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user, try again"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created succesfully"})
}

func LoginUserHandler(c *gin.Context) {
	var user models.LoginRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request"})
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		return
	}
	token, err := utils.CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "error creating token"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})

}

func LogoutUserHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token is required"})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// err := deleteses
}

func CurrentUserHandler(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The username cookie could not be read"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username})

}
