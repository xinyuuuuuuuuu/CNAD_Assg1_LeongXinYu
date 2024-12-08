package controller

import (
	"database/sql"
	"fmt"
	"strings"
)

// view all available vehicles
func ViewAvailableVehicles(db *sql.DB, userId string) {
	// retrieve all available vehicles to rent
	query := `
	SELECT VehicleId, VehicleMake, VehicleModel, VehicleType, LicensePlate, VehicleLocation, VehicleChargeLevel, VehicleCleanliness
	FROM Vehicle
	WHERE VehicleStatus = "A"	
	`
	// execute the query
	results, err := db.Query(query)
	if err != nil {
		fmt.Println("Error retrieving vehicles ", err)
		return
	}

	// close the result when the func has ended
	defer results.Close()

	// header for displaying vehicles
	fmt.Println("Available Vehicles")
	fmt.Printf("%-10s %-10s %-10s %-10s %-15s %-28s %-15s %-10s\n", "Make", "Model", "Type", "License Plate", "Location", "Charge Level", "Cleanliness")
	fmt.Println(strings.Repeat("-", 120))

	for results.Next() != false {
		var vId, vMake, vModel, vType, licensePlate, loc, cleanliness string
		var chargeLvl int 

		// scan to get result of each row
		err := results.Scan(&vId, &vMake, &vModel, &vType, &licensePlate, &loc, &chargeLvl, &cleanliness)

		// if there is error
		if err != nil {
			fmt.Println("There is an error in scanning vehicle data ", err)
			return
		}

		// display the result
		fmt.Printf("%-10s %-10s %-10s %-10s %-15s %-28s %-15d %-10s\n", vId, vMake, vModel, vType, licensePlate, loc, chargeLvl, cleanliness)
	}
	
	// checking for any errs aft iteration is done
	if err = results.Err(); err != nil {
		fmt.Println("Error iterating over vehicle data ", err)
		return
	}
}


