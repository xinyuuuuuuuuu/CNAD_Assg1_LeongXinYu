package main

import (
	"fmt"
	"log"
	"net/http"
	"vehicle-service/db"
	"vehicle-service/reservation"
	"vehicle-service/vehicle"

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

	// Set up reservation and vehicle routes
	reservation.RegisterRoutes(router, database)
	vehicle.RegisterRoutes(router, database) // <-- Add this line

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start server with CORS middleware
	port := ":8081"
	fmt.Println("Vehicle Service is running on port", port)
	log.Fatal(http.ListenAndServe(port, corsHandler.Handler(router)))
}
