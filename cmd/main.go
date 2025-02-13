package main

import (
	"fr_lab_1/pkg/handlers"
	"fr_lab_1/pkg/middleware"
	"log"
	"net/http"
)

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle routes
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/api/login", handlers.Login)
	http.HandleFunc("/api/register", handlers.Register)
	http.HandleFunc("/api/logout", handlers.Logout)
	http.HandleFunc("/api/users", handlers.Users)
	http.HandleFunc("/api/authority", handlers.Authority)

	// Apply middleware
	handler := middleware.LoggingMiddleware(middleware.AuthMiddleware(http.DefaultServeMux))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v", err)
	}
}
