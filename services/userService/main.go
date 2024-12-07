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
		//fmt.Println("\n")

		switch opt {
		case 1: 
			controller.Signup(db)

		case 2:
			controller.Login(db)

		case 3:
			println("Exiting out of application...")
			return

		default:
			fmt.Println("Invalid option, please try again.")
		}
	}

	// LoggedIn menu for users who had login to their account
	//func LoggedInMenu ()






}