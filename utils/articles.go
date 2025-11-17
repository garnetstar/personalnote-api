package utils

import (
	"database/sql"
	"fmt"
	"log"

	"personalnote.eu/simple-go-api/models"
)

// GetAllArticles retrieves all articles from the database (excluding deleted ones)
func GetAllArticles() ([]models.Article, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT id, title, content, updated, deleted 
		FROM article 
		WHERE deleted IS NULL 
		ORDER BY updated DESC, id DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var articles []models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Updated,
			&article.Deleted,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	log.Printf("📚 Retrieved %d articles from database", len(articles))
	return articles, nil
}

// GetArticleByID retrieves a single article by its ID
func GetArticleByID(id int) (*models.Article, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT id, title, content, updated, deleted 
		FROM article 
		WHERE id = ? AND deleted IS NULL
	`

	var article models.Article
	err := DB.QueryRow(query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Updated,
		&article.Deleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article with ID %d not found", id)
		}
		log.Printf("Error querying article by ID: %v", err)
		return nil, fmt.Errorf("failed to query article: %v", err)
	}

	log.Printf("📄 Retrieved article ID %d: %s", article.ID, article.Title)
	return &article, nil
}

func FindArticlesByTitle(title string) ([]models.Article, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT id, title, content, updated, deleted 
		FROM article 
		WHERE title LIKE ? AND deleted IS NULL
		ORDER BY updated DESC
	`

	rows, err := DB.Query(query, "%"+title+"%")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var articles []models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Updated,
			&article.Deleted,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	log.Printf("🔍 Found %d articles matching title '%s'", len(articles), title)
	return articles, nil
}

func FindArticlesByAll(keyword string) ([]models.Article, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT id, title, content, updated, deleted 
		FROM article 
		WHERE (title LIKE ? OR content LIKE ?) AND deleted IS NULL
		ORDER BY updated DESC
	`

	rows, err := DB.Query(query, "%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var articles []models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Updated,
			&article.Deleted,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	log.Printf("🔍 Found %d articles matching title '%s'", len(articles), keyword)
	return articles, nil
}

// UpdateArticle updates an existing article
func UpdateArticle(id int, title, content string) error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := `
		UPDATE article 
		SET title = ?, content = ?, updated = NOW() 
		WHERE id = ? AND deleted IS NULL
	`

	result, err := DB.Exec(query, title, content, id)
	if err != nil {
		log.Printf("Error updating article: %v", err)
		return fmt.Errorf("failed to update article: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article with ID %d not found", id)
	}

	log.Printf("✏️ Updated article ID %d: %s", id, title)
	return nil
}

// CreateArticle creates a new article in the database
func CreateArticle(title, content string) (int, error) {
	if DB == nil {
		return 0, fmt.Errorf("database connection not initialized")
	}

	query := `
		INSERT INTO article (title, content, updated) 
		VALUES (?, ?, NOW())
	`

	result, err := DB.Exec(query, title, content)
	if err != nil {
		log.Printf("Error creating article: %v", err)
		return 0, fmt.Errorf("failed to create article: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	log.Printf("✨ Created new article ID %d: %s", id, title)
	return int(id), nil
}

// DeleteArticle performs a soft delete on an article by setting the deleted timestamp
func DeleteArticle(id int) error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := `
		UPDATE article 
		SET deleted = NOW() 
		WHERE id = ? AND deleted IS NULL
	`

	result, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting article: %v", err)
		return fmt.Errorf("failed to delete article: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article with ID %d not found or already deleted", id)
	}

	log.Printf("🗑️ Soft deleted article ID %d", id)
	return nil
}
