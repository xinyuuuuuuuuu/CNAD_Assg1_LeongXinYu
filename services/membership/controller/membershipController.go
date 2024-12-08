package controller

import (
	"cnad_assg1_leongxinyu/services/membership/model"
	"cnad_assg1_leongxinyu/services/membership/utility"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// calculate membershipExpiryDate

// create membership details for new user
func CreateMembership(db *sql.DB, userId string) {

	var membership model.Membership

	// membershipId
	membershipId, err := utility.GenerateMembershipId(db)

	if err != nil {
		fmt.Println("Error generating membership id: ", err)
		return
	}
	membership.MembershipId = membershipId

	// assign userId
	membership.UserId = userId

	// assign MembershipTier, new member alws start from basic tier
	membership.MembershipTier = "Basic"

	// assign Hourly Rate , new member would be 15.00
	membership.HourlyRate = 15.00

	// assign Member Discount, new member start from 0.00
	membership.MemberDiscount = 0.00

	// assign Priority Level, new member start from 0
	membership.PiorityLevel = 0

	// assign total cost per month, new account starts from 0
	membership.TotalCostPerMonth = 0.00

	// assign membership expiry date
	membership.MembershipExpiryDate = time.Now().AddDate(0, 2, 0)

	// assign eligible for upgrade next month, new account starts from 'F': false
	membership.EligibleForUpgradeNextMonth = "F"

	// insert data into Membership table
	query := `
	INSERT INTO Membership(MembershipId, UserId, MembershipTier, HourlyRate, MemberDiscount, PriorityLevel, TotalCostPerMonth, MembershipExpiryDate, EligibleForUpgradeNextMonth)
	VALUES(?,?,?,?,?,?,?,?,?)
	`
	// Execute query and insert attributes
	_, err = db.Exec(query, membership.MembershipId, membership.UserId, membership.MembershipTier, membership.HourlyRate, membership.MemberDiscount, membership.PiorityLevel, membership.TotalCostPerMonth, membership.MembershipExpiryDate, membership.EligibleForUpgradeNextMonth)

	// if there is an error
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Membership successfully created.")
}

// view membership details
func ViewMembership(db *sql.DB, userId string) {
	//query to get Membership details
	query := `
	SELECT MembershipTier, HourlyRate, MemberDiscount, PriorityLevel, TotalCostPerMonth, MembershipExpiryDate, EligibleForUpgradeNextMonth
	FROM Membership
 	WHERE UserId = ?
 	`

	// execute the query to look for membership details
	results, err := db.Query(query, userId)

	// if there is error retrieving for data
	if err != nil {
		fmt.Println("Error retrieving for data ", err)
		return
	}

	// close result when the func has ended
	defer results.Close()

	// var to check if results exists
	hasResult := false

	// membership details headers
	fmt.Println("Membership Details")
	fmt.Printf("%-17s %-13s %-20s %-22s %-24s %-26s %-32s\n",
		"Membership Tier",
		"Hourly Rate",
		"Member Discount (%)",
		"Priority Level (0-2)",
		"Total Spent Per Month",
		"Membership Expiry Date",
		"Eligible to upgrade to next tier")
	fmt.Println(strings.Repeat("-", 163))

	// when results exist
	for results.Next() != false {
		hasResult = true // record exist

		var memTier, memExpiryDate, eliForUpgrade string
		var hrlyRate, memDisc, TSPM float64 // TSPM - total spent per month
		var priorLvl int

		// scan to get result of each row
		err := results.Scan(&memTier, &hrlyRate, &memDisc, &priorLvl, &TSPM, &memExpiryDate, &eliForUpgrade)

		// if there is error
		if err != nil {
			fmt.Println("There is an error in scanning data for membership ", err)
			return
		}

		// display the results
		fmt.Printf("%-17s %-13.2f %-20.2f %-22d %-24.2f %-26s %-32s\n", memTier, hrlyRate, memDisc, priorLvl, TSPM, memExpiryDate, eliForUpgrade)
	}

	// checking for any errors aft each iteration is done
	if err = results.Err(); err != nil {
		fmt.Println("Error iterating over membership details ", err)
		return
	}

	// if result doesn't exist
	if !hasResult {
		fmt.Println("No membership record for user")
		return
	}
}
