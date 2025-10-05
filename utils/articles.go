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
		ORDER BY id DESC
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
		SELECT id, title, content 
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
		SELECT id, title, content 
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
