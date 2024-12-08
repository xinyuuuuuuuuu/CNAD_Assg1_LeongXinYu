package main

import (
	billingController "cnad_assg1_leongxinyu/services/billingService/controller"
	membershipController "cnad_assg1_leongxinyu/services/membership/controller"
	userController "cnad_assg1_leongxinyu/services/userService/controller"

	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// connect to database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/cnad_assg1")

	// handle error
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	for {
		var opt int
		fmt.Println("============")
		fmt.Println("User Console")
		fmt.Println("1. Sign up for an account") // done
		fmt.Println("2. Login to account")       // done
		fmt.Println("3. Quit")
		fmt.Println("Enter an option: ")
		fmt.Scanln(&opt)
		fmt.Println("\n")

		switch opt {
		case 1:
			userController.Signup(db)

		case 2:
			userId := userController.Login(db)
			if userId != "" {
				LoggedInMenu(db, userId)
			} else {
				fmt.Println("Login has failed. Please try again...")
			}

		case 3:
			println("Exiting out of application...")
			return

		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}

// menu that members have access to after logging into their account
func LoggedInMenu(db *sql.DB, userId string) {
	for {
		var opt int
		fmt.Println("============")
		fmt.Println("Welcome!")
		fmt.Println("Member Console")
		fmt.Println("1. Update user details") //(done)
		fmt.Println("2. View Membership Details")
		fmt.Println("3. View Rental history")     //dk if wan or not - might need to create rentalHistory table
		fmt.Println("4. View Billing History")    // must do - (done)
		fmt.Println("5. View Available Vehicles") //must do
		fmt.Println("6. View Reservations")       //must do
		fmt.Println("7. Update Reservation")      //must do, should update n delete reservation be tgt
		fmt.Println("8. Log Out")
		fmt.Println("Enter an option: ")
		fmt.Scanln(&opt)
		fmt.Println("\n")

		switch opt {
		case 1:
			userController.UpdateUserDetails(db, userId)

		case 2:
			membershipController.ViewMembership(db, userId)

		case 3:

		case 4:
			billingController.GetPastBilling(db, userId)

		case 5:

		case 6:

		case 7:

		case 8:
			println("Logging out of account...")
			return

		default:
			fmt.Println("Invalid option, please try again.")
		}
	}

}
