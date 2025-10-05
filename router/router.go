package router

import (
	"net/http"

	"personalnote.eu/simple-go-api/handlers"
	"personalnote.eu/simple-go-api/middleware"
)

// SetupRoutes configures all the application routes
func SetupRoutes() {
	register := func(pattern string, handler http.HandlerFunc) {
		http.Handle(pattern, middleware.WithCORS(handler))
	}

	register("/", handlers.HelloHandler)
	register("/articles", handlers.ArticlesHandler)
	register("/article/filter/", handlers.ArticleFindHandler)
	register("/article/", handlers.ArticleByIDHandler)
}
