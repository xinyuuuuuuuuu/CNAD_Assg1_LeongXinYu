package controller

import (
	"database/sql"
	"fmt"
	"strings"
)

// view reservation
func ViewReservation(db *sql.DB, userId string) {
	// query to view reservation details
	query := `
	SELECT VehicleId, ReserveStatus, ReserveStartDate, ReserveEndDate, EstimatedTotalCost
	FROM Reservation
	WHERE UserId = ?
	`
	// execute the query
	results, err := db.Query(query, userId)
	if err != nil {
		fmt.Println("Error retrieving reservations ", err)
		return
	}

	// close the result when the func has ended
	defer results.Close()

	// header for displaying reservation
	fmt.Println("Reservations")
	fmt.Printf("%-10s %-18s %-22s %-22s %-10s\n", "VehicleId", "Reserve Status", "Start Date", "End Date", "Estimated Cost")
	fmt.Println(strings.Repeat("-", 100))

	for results.Next() != false {
		var vId, resStat, sD, eD string
		var estimCost float64
		//scan to get results of each row
		err := results.Scan(&vId, &resStat, &sD, &eD, &estimCost)

		// if there is error
		if err != nil {
			fmt.Println("There is an error in scanning reservation data ", err)
			return
		}

		// display the result
		fmt.Printf("%-10s %-18s %-22s %-22s %-10.2f\n", vId, resStat, sD, eD, estimCost)
	}

	// checking for any errs aft iteration is done
	if err = results.Err(); err != nil {
		fmt.Println("Error iterating over vehicle data ", err)
		return
	}
}

// update reservation
// delete reservation
// create reservation
