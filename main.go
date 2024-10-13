package main

import (
	"blogapi-naveen/db"
	"blogapi-naveen/handlers"
	"blogapi-naveen/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()

	// Public routes
	r.POST("/user/signup", handlers.CreateUserHandler)
	r.POST("/user/login", handlers.LoginUserHandler)
	r.POST("/users/logout", handlers.LogoutUserHandler)
	r.GET("/blogs", handlers.GetAllBlogsHandler)

	// Routes requiring authentication
	authorized := r.Group("/")
	authorized.Use(utils.AuthMiddleware())
	{
		authorized.GET("/blogs/:id", handlers.GetBlogByIdHandler)
		authorized.POST("/blogs/create", handlers.CreatePostHandler)
		authorized.DELETE("/blogs/:id", handlers.DeleteBlogsByIdHandler)
		authorized.GET("/blogs/user/:username", handlers.GetBlogByUsername)
	}

	// Start the server
	r.Run()
}
