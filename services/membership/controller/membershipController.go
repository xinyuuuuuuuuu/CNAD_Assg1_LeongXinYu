package controller 

import (
	"database/sql"
)

// calculate membershipExpiryDate

// create membership details for new user
//func CreateMembership(db *sqlDB, userId string) {
// memberid

//}

// view membership details
func ViewMembership(db *sql.DB, userId string) {
	// 
	query := `
	SELECT * FROM Membership
	`

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
