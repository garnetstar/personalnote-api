package handlers

import (
	"encoding/json"
	"fmt"
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

// ArticlesHandler handles requests for listing all articles
func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if !utils.ValidateHTTPMethod(w, r, http.MethodGet) {
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
	if !utils.ValidateHTTPMethod(w, r, http.MethodGet) {
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

// ArticleFindHandler handles requests for finding articles with filter parameters
func ArticleFindHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if !utils.ValidateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	// Extract filter parameters from URL path
	// Expected format: /article/filter/aaa/bbb
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 || !(parts[2] == "title" || parts[2] == "all") {
		utils.SendErrorResponse(w, http.StatusBadRequest,
			"Invalid URL", "Expected format: /article/filter/(title|all)/{param2}")
		return
	}

	log.Printf("Extracted parts: %v", parts)

	param1 := parts[2]
	keyword := parts[3]

	var article interface{}
	var err error

	switch param1 {
	case "title":
		article, err = utils.FindArticlesByTitle(keyword)
	case "all":
		article, err = utils.FindArticlesByAll(keyword)
	}

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.SendErrorResponse(w, http.StatusNotFound,
				"Article not found", fmt.Sprintf("Article with title '%s' not found", keyword))
		} else {
			log.Printf("Error fetching article: %v", err)
			utils.SendErrorResponse(w, http.StatusInternalServerError,
				"Database error", "Failed to retrieve article from database")
		}
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, article)

}

// UserHandler handles user creation
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if !utils.ValidateHTTPMethod(w, r, http.MethodPost) {
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest,
			"Invalid JSON", "Could not parse request body")
		return
	}

	// Validate user data
	if user.Name == "" || user.ID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest,
			"Invalid user data", "Name and ID are required")
		return
	}

	log.Printf("👤 Creating user: %s (ID: %d)", user.Name, user.ID)
	utils.SendJSONResponse(w, http.StatusCreated, user)
}
