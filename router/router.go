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

	// Static file serving
	http.Handle("/garnetstar.ico", middleware.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/garnetstar.ico")
	})))
	http.Handle("/garnetstar.jpeg", middleware.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/garnetstar.jpeg")
	})))

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
