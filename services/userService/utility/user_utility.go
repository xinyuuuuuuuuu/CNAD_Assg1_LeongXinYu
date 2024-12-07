package utility

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// generate userid
func GenerateUserId(db*sql.DB) (string, error) {
	var lastUserId string

	// query to get the lastest userid, LIMIT 1 ensures only one row is returned
	query := "SELECT UserId from UserService ORDER BY UserId DESC LIMIT 1"

	err := db.QueryRow(query).Scan(&lastUserId) // if there is userid found, lastUserId = userId
	
	if err != nil{
		// if query selects no row
		if err == sql.ErrNoRows {
			fmt.Println("No rows found. Starting with U0001")
			lastUserId = "U0001"
		}

		return "", err
	}

	// remove 'U' and convert num to int
	var lastNum int
	lastNum, err = strconv.Atoi(strings.TrimPrefix(lastUserId, "U"))
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserId: %v", err) // returns empty userId and err msg
	}

	// increment num n format the new id
	newUserId := fmt.Sprintf("U%04d", lastNum + 1)
	return newUserId, nil

}