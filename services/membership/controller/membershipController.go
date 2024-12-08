package controller

import (
	"database/sql"
	"fmt"
	"strings"
)

// calculate membershipExpiryDate

// create membership details for new user
//func CreateMembership(db *sqlDB, userId string) {
// memberid

//}

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

// User sign up for an account - POST
// func Signup(db *sql.DB) {
// 	reader := bufio.NewReader(os.Stdin)

// 	var user model.UserService
// 	var userId string
// 	var err error

// 	// userid
// 	userId, err = utility.GenerateUserId(db)
// 	if err != nil {
// 		fmt.Println("Error generating user id: ", err)
// 		return
// 	}
// 	user.UserId = userId

// 	// Name
// 	fmt.Print("Name: ")
// 	user.Name, _ = reader.ReadString('\n')
// 	user.Name = strings.TrimSpace(user.Name)

// 	// Email
// 	for {
// 		fmt.Print("Email: ")
// 		user.Email, _ = reader.ReadString('\n')
// 		user.Email = strings.TrimSpace(user.Email)
// 		// if email input contains "@"
// 		if strings.Contains(user.Email, "@") {
// 			break;
// 		}

// 		fmt.Println("Invalid email format. Please try again.")
// 	}

// 	// Password
// 	fmt.Print("Password: ")
// 	var pw string
// 	pw, _ = reader.ReadString('\n')
// 	pw = strings.TrimSpace(pw)
// 	hash, err := HashedPassword(pw)
// 	if err != nil {
// 		fmt.Println("Error hashing password ", err)
// 	}

// 	// store hashed password
// 	user.Password = hash
// 	//fmt.Println("Hash: ", user.Password) // checking purpose

// 	// ContactNo
// 	fmt.Print("Contact Number: ")
// 	user.ContactNo, _ = reader.ReadString('\n')
// 	user.ContactNo = strings.TrimSpace(user.ContactNo)

// 	// Dob
// 	for {
// 		fmt.Print("Date of Birth (YYYY-MM-DD): ")
// 		dobInput, _ := reader.ReadString('\n') // Read the input as a string
// 		dobInput = strings.TrimSpace(dobInput)
// 		user.Dob, err = time.Parse("2006-01-02", dobInput)
// 		if err != nil {
// 			fmt.Println("Date format is invalid. Please use YYYY-MM-DD.")
// 			continue // user can input if format was prev invalid
// 		}
// 		break // break loop when user input valid date
// 	}

// 	// Address
// 	fmt.Print("Address: ")
// 	user.Address, _ = reader.ReadString('\n')
// 	user.Address = strings.TrimSpace(user.Address)

// 	// CreatedDateTime
// 	user.CreatedDateTime = time.Now()

// 	// insert data into UserService table
// 	query := `
// 	INSERT INTO UserService
// 	(UserId, Name, Email, Password, ContactNo, Dob, Address, CreatedDateTime)
// 	VALUES(?,?,?,?,?,?,?,?)
// 	`
// 	result, err := db.Exec(query, user.UserId, user.Name, user.Email, user.Password, user.ContactNo, user.Dob, user.Address, user.CreatedDateTime)

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println("Number of rows affected: ", rows)

// }
