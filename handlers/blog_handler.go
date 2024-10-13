package handlers

import (
	"blogapi-naveen/db"
	"blogapi-naveen/models"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	username := c.GetString("username")
	fmt.Printf("%v\n", username)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Username not found in context"})
		return
	}

	var blog models.Blog

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request"})
		return
	}

	blog.Username = username

	if err := blog.BlogSave(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not post blog, try again later"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Blog Created!"})
}

func GetAllBlogsHandler(c *gin.Context) {
	var blogs []models.Blog

	query := `SELECT * FROM blogs`
	rows, err := db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error retrieving blogs", "error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var blog models.Blog
		if err := rows.Scan(&blog.BlogID, &blog.Title, &blog.Content, &blog.Category, &blog.PublishedDate, &blog.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error processing blog data"})
			return
		}
		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error retrieving blog data"})
	}
	c.JSON(http.StatusOK, gin.H{"blogs": blogs})
}

func GetBlogByIdHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "authorization token is required"})
		return
	}

	id := c.Param("id")
	query := `SELECT * FROM blogs WHERE blog_id = $1`

	row := db.DB.QueryRow(query, id)

	var blog models.Blog
	err := row.Scan(&blog.BlogID, &blog.Title, &blog.Content, &blog.Category, &blog.PublishedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "blog not found"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving blog"})
	}
	c.JSON(http.StatusOK, gin.H{"blog": blog})
}

func DeleteBlogsByIdHandler(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM blogs WHERE blog_id = $1`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting blog", "error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error checking result", "error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Blog not found"})
		return
	}
}

func GetBlogByUsername(c *gin.Context) {
	var blogs []models.Blog
	username := c.Param("username")
	fmt.Printf("username %v\n", username)
	query := `SELECT * FROM blogs WHERE username = $1`

	rows, err := db.DB.Query(query, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error retrieving blogs"})
		return
	}
	for rows.Next() {
		var blog models.Blog
		err := rows.Scan(&blog.BlogID, &blog.Title, &blog.Content, &blog.Category, &blog.PublishedDate, &blog.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error processing blog data"})
			return
		}
		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error retrieving blog data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": blogs})
}
