package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB
// connection to database
func ConnectDB()(*sql.DB, error) {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// retrieve database credentials from .env file
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// open connection
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}

	// test connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Database is not reachable: %v", err)
	}

	fmt.Println("Connected to the database succesfully.")
	return db, nil
}

