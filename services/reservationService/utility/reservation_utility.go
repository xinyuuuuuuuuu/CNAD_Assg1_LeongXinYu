package utility

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// generate userid
func GenerateReservationId(db*sql.DB) (string, error) {
	var lastReserveId string

	// query to get the lastest userid, LIMIT 1 ensures only one row is returned
	query := "SELECT ReservationId from Reservation ORDER BY UserId DESC LIMIT 1"

	err := db.QueryRow(query).Scan(&lastReserveId) // if there is reserve id found, lastReserveId = reserveId
	
	if err != nil{
		// if query selects no row
		if err == sql.ErrNoRows {
			fmt.Println("No rows found. Starting with R0001")
			return "R0001", nil
		}

		return "", err
	}

	// remove 'R' and convert num to int
	var lastNum int
	lastNum, err = strconv.Atoi(strings.TrimPrefix(lastReserveId, "R"))
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserId: %v", err) // returns empty userId and err msg
	}

	// increment num n format the new id
	newReserveId := fmt.Sprintf("R%04d", lastNum + 1)
	return newReserveId, nil
}