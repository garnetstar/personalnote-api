package router

import (
	"net/http"

	"personalnote.eu/simple-go-api/handlers"
	"personalnote.eu/simple-go-api/utils"
)

// SetupRoutes configures all the application routes
func SetupRoutes() {
	http.HandleFunc("/", utils.CORSHandlerFunc(handlers.HelloHandler))
	http.HandleFunc("/user", utils.CORSHandlerFunc(handlers.UserHandler))
	http.HandleFunc("/articles", utils.CORSHandlerFunc(handlers.ArticlesHandler))
	http.HandleFunc("/article/filter/", utils.CORSHandlerFunc(handlers.ArticleFindHandler))
	http.HandleFunc("/article/", utils.CORSHandlerFunc(handlers.ArticleByIDHandler))
}
