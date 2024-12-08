package utility

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// generate membershipId
func GenerateMembershipId(db *sql.DB) (string, error) {
	var lastMemId string

	// query to get lastes MemId, LIMIT 1 ensures only one row is returned
	query := "SELECT MembershipId from Membership ORDER BY MembershipId DESC LIMIT 1"

	// run query and return the most recent useruid
	err := db.QueryRow(query).Scan(&lastMemId)

	// if there is error
	if err != nil {
		// if query selects no row
		if err == sql.ErrNoRows {
			fmt.Println("No rows found. Starting with M0001")
			return "M0001", nil
		}
		return "", err
	}

	// remove prefix "M" n convert num to int
	var lastNum int
	lastNum, err = strconv.Atoi(strings.TrimPrefix(lastMemId, "M"))

	if err != nil {
		return "", fmt.Errorf("Failed to parse MembershipId: %v", err) // return empty membershipId and err msg
	}

	// increment num by 1 n format the new id
	newMemId := fmt.Sprintf("M%04d", lastNum+1)
	return newMemId, nil

}
