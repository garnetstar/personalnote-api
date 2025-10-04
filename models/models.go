package models

import "time"

// Response represents a standard API response
type Response struct {
	Message string `json:"message"`
}

// User represents a user entity
type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// Article represents an article entity from the database
type Article struct {
	ID      int        `json:"id" db:"id"`
	Title   string     `json:"title" db:"title"`
	Content string     `json:"content" db:"content"`
	Updated *time.Time `json:"updated" db:"updated"`
	Deleted *time.Time `json:"deleted" db:"deleted"`
}

// ArticleListResponse represents a response containing multiple articles
type ArticleListResponse struct {
	Articles []Article `json:"articles"`
	Count    int       `json:"count"`
	Message  string    `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
