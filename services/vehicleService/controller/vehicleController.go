package controller

import "database/sql"

// view all available vehicles
func ViewAvailableVehicles(db *sql.DB, userId string) {
	// retrieve all available vehicles 

}

// // get past billing
// func GetPastBilling(db *sql.DB, userId string) {
// 	// retrieve past billing for user
// 	query := `
// 		SELECT BillingDate, BIllingTotal, PaymentMethod, PaymentStatus
// 		FROM Billing
// 		WHERE UserId = ?
// 		`
// 	// execute the query to look for current user details
// 	results, err := db.Query(query, userId)

// 	// if there is error in retrieving for data
// 	if err != nil {
// 		fmt.Println("Error retrieving for data ", err)
// 		return
// 	}

// 	// close the result when the func has ended
// 	defer results.Close()

// 	// var to check if results exist
// 	hasResult := false

// 	// billing headers
// 	fmt.Println("Past Billings")
// 	fmt.Printf(("%-21s %-10s %-16s %-16s\n"), "Date", "Total", "Payment Method", "Payment Status")
// 	fmt.Println(strings.Repeat("-", 60))

// 	// when results exist
// 	for results.Next() != false { // .Next returns true/ false
// 		hasResult = true // hv 1 or more record

// 		var billDate, paymMethod, payStatus string
// 		var billTotal float64

// 		// scan to get result of each row
// 		err := results.Scan(&billDate, &billTotal, &paymMethod, &payStatus)

// 		// if there is error
// 		if err != nil {
// 			fmt.Println("There is an error in scanning data for billing ", err)
// 			return
// 		}

// 		// display the results, %.2f - return 2 dp for float
// 		fmt.Printf("%-21s %-10.2f %-15s %-15s\n", billDate, billTotal, paymMethod, payStatus)
// 	}

// 	// checking for any errs after iteration is done
// 	if err = results.Err(); err != nil {
// 		fmt.Println("Error iterating over billing data ", err)
// 		return
// 	}

// 	// when result doesn't exist
// 	if !hasResult {
// 		fmt.Println("No past billing record for user")
// 		return
// 	}
// }

