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

	// Public routes
	register("/", handlers.HelloHandler)
	register("/articles", handlers.ArticlesHandler)
	register("/article/filter/", handlers.ArticleFindHandler)
	register("/article/", handlers.ArticleHandler)

	// Auth routes
	register("/auth/google/login", handlers.GoogleLoginHandler)
	register("/auth/google/callback", handlers.GoogleCallbackHandler)
	register("/auth/user", handlers.UserInfoHandler)
}
