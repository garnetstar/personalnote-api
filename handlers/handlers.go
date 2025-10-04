package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"personalnote.eu/simple-go-api/models"
	"personalnote.eu/simple-go-api/utils"
)

var counter int

// HelloHandler handles the root endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	counter++
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Counter", strconv.Itoa(counter))

	resp := models.Response{Message: "Hallo, from Go!"}
	log.Printf("Received request #%d", counter)
	log.Printf("Query params: %v", r.URL.Query())
	log.Printf("URI: %v", r.RequestURI)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UserHandler handles user-related requests
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed,
			"Method not allowed", "Only POST requests are accepted")
		return
	}

	// Read and parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest,
			"Invalid request body", "Could not read request body")
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest,
			"Invalid JSON", "Could not parse JSON request body")
		return
	}

	// Validate the request
	validationErrors := utils.ValidateUser(user)
	if len(validationErrors) > 0 {
		log.Printf("Validation failed: %v", validationErrors)
		utils.SendErrorResponse(w, http.StatusBadRequest,
			"Validation failed",
			fmt.Sprintf("Validation errors: %s", strings.Join(validationErrors, ", ")))
		return
	}

	// Log the valid request to console
	log.Printf("✅ Received valid user data: Name=%s, ID=%d", user.Name, user.ID)
	log.Printf("📝 Full request body: %s", string(body))

	// Send success response
	message := fmt.Sprintf("User %s with ID %d has been processed successfully", user.Name, user.ID)
	utils.SendSuccessResponse(w, message)
}

// ArticlesHandler handles requests for listing all articles
func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed,
			"Method not allowed", "Only GET requests are accepted")
		return
	}

	log.Printf("📚 Fetching all articles from database")

	// Get articles from database
	articles, err := utils.GetAllArticles()
	if err != nil {
		log.Printf("Error fetching articles: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError,
			"Database error", "Failed to retrieve articles from database")
		return
	}

	// Create response
	response := models.ArticleListResponse{
		Articles: articles,
		Count:    len(articles),
		Message:  fmt.Sprintf("Successfully retrieved %d articles", len(articles)),
	}

	// Log success
	log.Printf("✅ Successfully fetched %d articles", len(articles))

	// Send JSON response
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// ArticleByIDHandler handles requests for getting a specific article by ID
func ArticleByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed,
			"Method not allowed", "Only GET requests are accepted")
		return
	}

	// Extract ID from URL path
	// Expected format: /article/123
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "article" {
		utils.SendErrorResponse(w, http.StatusBadRequest,
			"Invalid URL", "Expected format: /article/{id}")
		return
	}

	// Parse article ID
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest,
			"Invalid ID", "Article ID must be a valid integer")
		return
	}

	log.Printf("📄 Fetching article with ID: %d", id)

	// Get article from database
	article, err := utils.GetArticleByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.SendErrorResponse(w, http.StatusNotFound,
				"Article not found", fmt.Sprintf("Article with ID %d not found", id))
		} else {
			log.Printf("Error fetching article: %v", err)
			utils.SendErrorResponse(w, http.StatusInternalServerError,
				"Database error", "Failed to retrieve article from database")
		}
		return
	}

	// Log success
	log.Printf("✅ Successfully fetched article: %s (ID: %d)", article.Title, article.ID)

	// Send JSON response
	utils.SendJSONResponse(w, http.StatusOK, article)
}
