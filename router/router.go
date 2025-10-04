package router

import (
	"net/http"

	"personalnote.eu/simple-go-api/handlers"
)

// SetupRoutes configures all the application routes
func SetupRoutes() {
	http.HandleFunc("/", handlers.HelloHandler)
	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/articles", handlers.ArticlesHandler)
	http.HandleFunc("/article/", handlers.ArticleByIDHandler)
}
