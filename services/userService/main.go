package main

import (
	"cnad_assg1_leongxinyu/services/userService/controller"
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
		fmt.Println("1. Sign up for an account")
		fmt.Println("2. Login to account")
		fmt.Println("3. Quit")
		fmt.Println("Enter an option: ")
		fmt.Scanln(&opt)
		fmt.Println("\n")

		switch opt {
		case 1: 
			controller.Signup(db)

		case 2:
			userId := controller.Login(db)
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
		fmt.Println("1. Update user details")
		fmt.Println("2. View Membership Details")
		fmt.Println("3. View Rental history") //dk if wan or not - might need to create rentalHistory table
		fmt.Println("4. View Billing History") // must do
		fmt.Println("5. View Available Vehicles") //must do
		fmt.Println("6. View Reservations") //must do
		fmt.Println("7. Update Reservation") //must do, should update n delete reservation be tgt
		fmt.Println("8. Quit")
		fmt.Println("Enter an option: ")
		fmt.Scanln(&opt)
		fmt.Println("\n")

		switch opt {
			case 1: 
				controller.UpdateUserDetails(db, userId)
	
			case 2:


			case 3:
				

			case 4:

			case 5:


			case 6:


			case 7:
			
	
			case 8:
				println("Exiting out of application...")
				return
	
			default:
				fmt.Println("Invalid option, please try again.")
			}
	}

}