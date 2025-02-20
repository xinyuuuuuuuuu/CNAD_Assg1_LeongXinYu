package main

import (
	"fmt"
	"log"
	"net/http"
	"user-management/db"
	"user-management/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Connect to the database
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer database.Close()

	// Initialize router
	router := mux.NewRouter()

	// Set up user routes
	routes.UserRoutes(router, database)

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Change this to specific origins if needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start server with CORS middleware
	port := ":8080"
	fmt.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(port, corsHandler.Handler(router)))
}
