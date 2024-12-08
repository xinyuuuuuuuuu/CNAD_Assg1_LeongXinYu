package controller

import (
	"bufio"
	"cnad_assg1_leongxinyu/services/membership/utility"
	"cnad_assg1_leongxinyu/services/reservationService/model"
	vehicleController "cnad_assg1_leongxinyu/services/vehicleService/controller"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
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

	// var to check if results exist
	hasResult := false

	// header for displaying reservation
	fmt.Println("Reservations")
	fmt.Printf("%-10s %-18s %-22s %-22s %-10s\n", "VehicleId", "Reserve Status", "Start Date", "End Date", "Estimated Cost")
	fmt.Println(strings.Repeat("-", 100))

	for results.Next() != false {
		hasResult = true
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

	// when result doesn't exist
	if !hasResult {
		fmt.Println("No reservation")
		return
	}
}

// create reservation
func ReserveVehicle(db *sql.DB, userId string) {
	if !vehicleController.DisplayAvailableVehicles(db) {
		// if no vehicles are available, func will be stopped
		fmt.Println("No available vehicles to reserve")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var reserve model.Reservation

	// prompt for vehicleId
	fmt.Print("\nEnter the vehicle id of the vehicle you wish to reserve: ")
	reserve.VehicleId, _ = reader.ReadString('\n')
	reserve.VehicleId = strings.TrimSpace(reserve.VehicleId)

	// validate the vehicle id
	query := `
	SELECT COUNT(*)
	FROM Vehicle
	WHERE VehicleStatus = "A" AND VehicleId = ?
	`

	var count int
	err := db.QueryRow(query, reserve.VehicleId).Scan(&count)
	// if there is err
	if err != nil {
		fmt.Println("Error validating vehicle Id: ", err)
		return
	}

	// if vehicle id is not available or invalid
	if count == 0 {
		fmt.Println("Invalid vehicle Id or vehicle is not available. Please try again.")
		return
	}

	// prompt for reserve start n end date
	for {
		fmt.Print("Enter reservation start date (YYYY-MM-DD HH:MM:SS): ")
		startDate, _ := reader.ReadString('\n')
		startDate = strings.TrimSpace(startDate)
		reserve.ReserveStartDate, err = time.Parse("2006-01-02 21:19:09", startDate)
		if err != nil {
			fmt.Println("Date format is invalid. Please use YYYY-MM-DD HH:MM:SS.")
			continue // user can input if format was prev invalid
		}
		break // loop breaks when user input valid date
	}

	for {
		fmt.Print("Enter reservation end date (YYYY-MM-DD HH:MM:SS): ")
		endDate, _ := reader.ReadString('\n')
		endDate = strings.TrimSpace(endDate)
		reserve.ReserveEndDate, err = time.Parse("2006-01-02 21:19:09", endDate)
		if err != nil {
			fmt.Println("Date format is invalid. Please use YYYY-MM-DD HH:MM:SS.")
			continue // user can input if format was prev invalid
		}

		// check if end date is after start date
		if reserve.ReserveEndDate.Before(reserve.ReserveStartDate) {
			fmt.Println("End date must be after the start date. Please try again.")
			continue // input can be re-enter if end date is invalid
		}
		break // loop breaks when user input valid date
	}

	// calculate estimated cost

	// fetch member discount form membership
	var hrlyRate float64
	rateQuery := `
	SELECT HourlyRate
	FROM Membership
	WHERE UserId = ?
	`

	// execute query n store it
	err = db.QueryRow(rateQuery, userId).Scan(&hrlyRate)
	if err != nil {
		fmt.Println("Error fetching hourly rate for user ", err)
		return
	}

	// calculate duration of reservation
	duration := reserve.ReserveEndDate.Sub(reserve.ReserveStartDate).Hours()
	reserve.EstimatedTotalCost = hrlyRate * duration
	fmt.Println("Estimated Cost: ", reserve.EstimatedTotalCost)

	// generate reservation id
	reserveId, err := utility.GenerateMembershipId(db)
	if err != nil {
		fmt.Println("Error generating reservation id: ", err)
		return
	}

	reserve.ReservationId = reserveId

	// assign reserve status
	reserve.ReserveStatus = "Pend"

	// assign created date
	reserve.CreatedDate = time.Now()

	// query to insert data into reservation table
	insertQuery := `
	INSERT INTO Reservation
	(ReservationId, UserId, VehicleId, ReserveStatus, ReserveStartDate, ReserveEndDate, EstimatedTotalCost, CreatedDate)
	VALUES(?,?,?,?,?,?,?,?)
	`
	// Execute query and insert attributes
	_, err = db.Exec(insertQuery, reserve.ReservationId, reserve.UserId, reserve.VehicleId, reserve.ReserveStatus, reserve.ReserveStartDate, reserve.ReserveEndDate, reserve.EstimatedTotalCost, reserve.CreatedDate)
	if err != nil {
		fmt.Println("Error confirming Reservation ", err)
		return
	}

	fmt.Println("Successful Reservation")

	// change the vehicle status to "R" after the vehicle has been reserved
	updateVehAvailQuery := `
	UPDATE Vehicle
	SET VehicleStatus = "R"
	WHERE VehicleId = ?
	`
	_, err = db.Exec(updateVehAvailQuery, reserve.VehicleId)

	if err != nil {
		fmt.Println("Error updating vehicle status ", err)
		return
	}

	fmt.Println("Vehicle has been reserved successfully. Vehicle status has been updated.")
}

// update reservation
func UpdateReservation(db * sql.DB, userId string) {
	
	ViewReservation(db, userId)
	reader := bufio.NewReader(os.Stdin)
	var reserve model.Reservation

	// confirm reservation w user
	fmt.Print("Confirm this reservation (y/n): ")
	confirmInput, _ := reader.ReadString('\n')
	confirmInput = strings.TrimSpace(strings.ToLower(confirmInput))

	if confirmInput == "yes" {
		// update reservation status to "Conf"
		updateRSQuery := `
		UPDATE Reservation
		SET ReserveStatus = "Conf"
		WHERE ReservationId = ?
		`
		_, err := db.Exec(updateRSQuery, reserve.ReservationId)
		if err != nil {
			fmt.Println("Error confirming reservation: ", err)
			return
		}
		fmt.Println("Reservation confirmed successfully")
	} else {
		// user does not confirm reservation

		// update vehicle status to "A"
		changeVehStatusQuery := `
		UPDATE Vehicle
		SET VehicleStatus = "A"
		WHERE VehicleId = ?
		`

		_, err := db.Exec(changeVehStatusQuery, reserve.VehicleId)
		if err != nil {
			fmt.Println("Error changing vehicle status: ", err)
			return
		}

		// cancel reservation
		cancelQuery := `
		UPDATE Reservation
		SET ReserveStatus = "Canc"
		WHERE ReservationId = ?
		`
		_, err = db.Exec(cancelQuery, reserve.VehicleId)
		if err != nil {
			fmt.Println("Error changing vehicle status: ", err)
			return
		}

		fmt.Println("Reservation not confirmed. Reservation has been cancelled and vehicle status has been changed.")
	}
}
