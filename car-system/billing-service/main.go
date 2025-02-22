package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"billing-service/billing"
	"github.com/gorilla/mux"
	"github.com/rs/cors" // Import the CORS package
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load database connection
	db, err := connectDB()
	if err != nil {
		log.Fatal("[ERROR] Failed to connect to database:", err)
	}
	defer db.Close()

	// Create router
	router := mux.NewRouter()

	// Apply CORS middleware to specific routes
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:5500", // Allow your frontend domain
			"http://localhost:5500",  // Localhost frontend URL
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, // Allow specific methods
		AllowedHeaders: []string{"Authorization", "Content-Type"}, // Allow specific headers
	})

	// Register billing routes
	billing.RegisterBillingRoutes(router, db)

	// Apply CORS handler to the router
	http.Handle("/", corsHandler.Handler(router))

	// Start the server
	port := getPort()
	log.Printf("[INFO] Billing Service running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil)) // Use nil to pass the default mux (http.Handle)
}

// connectDB initializes a MySQL database connection
func connectDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Fallback to default values if not set
	if dbUser == "" {
		dbUser = "electric-car-sharing-admin"
	}
	if dbPassword == "" {
		dbPassword = "Password123"
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "billing_service"
	}

	// Create DSN string for MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	log.Println("[INFO] Connecting to MySQL database...")

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("[ERROR] Failed to open database connection:", err)
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Println("[ERROR] Database connection failed:", err)
		return nil, err
	}

	log.Println("[INFO] Successfully connected to MySQL database!")
	return db, nil
}

// getPort retrieves the server port from environment variables or defaults to 8080
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	return port
}
