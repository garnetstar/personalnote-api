package utils

import (
	"database/sql"
	"fmt"
	"log"

	"personalnote.eu/simple-go-api/models"
)

// CreateOrUpdateUser creates a new user or updates an existing one
func CreateOrUpdateUser(googleID, email, name, picture string) (*models.User, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	// Check if user exists
	var user models.User
	query := `SELECT id, google_id, email, name, picture, created_at, updated_at FROM users WHERE google_id = ?`
	err := DB.QueryRow(query, googleID).Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// User doesn't exist, create new one
		insertQuery := `INSERT INTO users (google_id, email, name, picture, created_at, updated_at) 
			VALUES (?, ?, ?, ?, NOW(), NOW())`
		result, err := DB.Exec(insertQuery, googleID, email, name, picture)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return nil, fmt.Errorf("failed to create user: %v", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get last insert ID: %v", err)
		}

		user.ID = int(id)
		user.GoogleID = googleID
		user.Email = email
		user.Name = name
		user.Picture = picture

		log.Printf("👤 Created new user: %s (%s)", name, email)
		return &user, nil
	} else if err != nil {
		log.Printf("Error querying user: %v", err)
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	// User exists, update info
	updateQuery := `UPDATE users SET email = ?, name = ?, picture = ?, updated_at = NOW() WHERE google_id = ?`
	_, err = DB.Exec(updateQuery, email, name, picture, googleID)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	user.Email = email
	user.Name = name
	user.Picture = picture

	log.Printf("👤 Updated user: %s (%s)", name, email)
	return &user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id int) (*models.User, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	var user models.User
	query := `SELECT id, google_id, email, name, picture, created_at, updated_at FROM users WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with ID %d not found", id)
	} else if err != nil {
		log.Printf("Error querying user: %v", err)
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	return &user, nil
}
