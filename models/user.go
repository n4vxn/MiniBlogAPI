package models

import (
	"blogapi-naveen/db"
	"blogapi-naveen/utils"
	"fmt"
	"strings"
	"time"
)

type User struct {
	UserID       int       `json:"user_id"`
	FirstName    string    `json:"first_name" binding:"required"`
	LastName     string    `json:"last_name" binding:"required"`
	Username     string    `json:"username" binding:"required"`
	PasswordHash string    `json:"password" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

func (u *User) UserSave() error {
	u.Username = strings.ToLower(u.Username)
	hashedPass, err := utils.HashPassword(u.PasswordHash)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	u.PasswordHash = hashedPass
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	

	query := `INSERT INTO users(first_name, last_name, username, password_hash, created_at)
	VALUES($1, $2, $3, $4, $5)
	RETURNING user_id`

	err = db.DB.QueryRow(query, u.FirstName, u.LastName, u.Username, u.PasswordHash, u.CreatedAt).Scan(&u.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (l *LoginRequest) ValidateCredentials() error {
	l.Username = strings.ToLower(l.Username)

	query := `SELECT password_hash FROM users WHERE username = $1`
	row := db.DB.QueryRow(query, l.Username)

	var retrivedPass string
	err := row.Scan(&retrivedPass)
	if err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}
	err = utils.CheckPassword(l.Password, retrivedPass)
	if err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}
	return nil
}

