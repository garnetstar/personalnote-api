package models

import "time"

// User represents a user entity from the database
type User struct {
	ID        int        `json:"id" db:"id"`
	GoogleID  string     `json:"google_id" db:"google_id"`
	Email     string     `json:"email" db:"email"`
	Name      string     `json:"name" db:"name"`
	Picture   string     `json:"picture" db:"picture"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// GoogleUserInfo represents the user info from Google
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
