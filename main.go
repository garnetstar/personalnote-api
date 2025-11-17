package main

import (
	"log"
	"net/http"

	"personalnote.eu/simple-go-api/handlers"
	"personalnote.eu/simple-go-api/router"
	"personalnote.eu/simple-go-api/utils"
)

func main() {
	// Initialize database connection
	if err := utils.InitDB(); err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Printf("🔄 Continuing without database - some endpoints may not work")
	} else {
		defer utils.CloseDB()
	}

	// Initialize OAuth
	handlers.InitOAuth()

	// Setup all routes
	router.SetupRoutes()

	// Start the server
	addr := ":8080"
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("📡 Endpoints available:")
	log.Printf("   GET  / - Health check")
	log.Printf("   POST /user - User management")
	log.Printf("   GET  /articles - List all articles")
	log.Printf("   GET  /article/{id} - Get article by ID")
	log.Printf("   GET  /auth/google/login - Google OAuth login")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("❌ Could not start server: %s", err)
	}
}
