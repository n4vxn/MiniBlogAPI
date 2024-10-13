package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
		DB.Close()
	}

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = DB.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	if err := CreateTables(); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

}

func CreateTables() error {
	queries := []struct {
		name  string
		query string
	}{
		{
			name: "users",
			query: `CREATE TABLE IF NOT EXISTS users(
		user_id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		},
		{
			name: "blogs",
			query: `CREATE TABLE IF NOT EXISTS blogs(
		blog_id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		category TEXT NOT NULL,
		published_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		username TEXT NOT NULL,
    	FOREIGN KEY (username) REFERENCES users(username));`,
		},
	}
	for _, table := range queries {
		if _, err := DB.Exec(table.query); err != nil {
			return fmt.Errorf("error creating %s table: %w", table.name, err)
		}
	}
	return nil
}
