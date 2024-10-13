package models

import (
	"blogapi-naveen/db"
	"time"
)

type Blog struct {
	BlogID        int       `json:"id"`
	Username      string    `json:"username"`
	Title         string    `json:"title"    binding:"required"`
	Content       string    `json:"content" binding:"required"`
	Category      string    `json:"category" binding:"required"`
	PublishedDate time.Time `json:"published_date"`
}

func (b *Blog) BlogSave() error {
	query := `
		INSERT INTO blogs(title, content, category, published_date, username)
		VALUES($1, $2, $3, $4, $5)
		RETURNING blog_id`

	err := db.DB.QueryRow(query, b.Title, b.Content, b.Category, b.PublishedDate, b.Username).Scan(&b.BlogID)
	if err != nil {
		return err
	}

	return nil
}
