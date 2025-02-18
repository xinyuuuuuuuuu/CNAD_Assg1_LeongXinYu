package main

import (
	"fmt"
	"log"
	"net/http"
	"user-management/db"
	"user-management/routes"

	"github.com/gorilla/mux"
)

func main() {
	// connect to database
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer database.Close()

	// initialize router
	router := mux.NewRouter()

	// set up user routes
	routes.UserRoutes(router, database)

	// start server
	port := ":8080"
	fmt.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
