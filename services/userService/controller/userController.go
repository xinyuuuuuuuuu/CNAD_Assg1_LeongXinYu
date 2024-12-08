package controller

import (
	"bufio"
	"cnad_assg1_leongxinyu/services/userService/model"
	"cnad_assg1_leongxinyu/services/userService/utility"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

// hash password
func HashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// check password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// User sign up for an account - POST
func Signup(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	var user model.UserService
	var userId string
	var err error

	// userid
	userId, err = utility.GenerateUserId(db)
	if err != nil {
		fmt.Println("Error generating user id: ", err)
		return
	}
	user.UserId = userId

	// Name
	fmt.Print("Name: ")
	user.Name, _ = reader.ReadString('\n')
	user.Name = strings.TrimSpace(user.Name)

	// Email
	for {
		fmt.Print("Email: ")
		user.Email, _ = reader.ReadString('\n')
		user.Email = strings.TrimSpace(user.Email)
		// if email input contains "@"
		if strings.Contains(user.Email, "@") {
			break;
		} 

		fmt.Println("Invalid email format. Please try again.")
	}


	// Password
	fmt.Print("Password: ")
	var pw string
	pw, _ = reader.ReadString('\n')
	pw = strings.TrimSpace(pw)
	hash, err := HashedPassword(pw)
	if err != nil {
		fmt.Println("Error hashing password ", err)
	}

	// store hashed password
	user.Password = hash
	//fmt.Println("Hash: ", user.Password) // checking purpose

	// ContactNo
	fmt.Print("Contact Number: ")
	user.ContactNo, _ = reader.ReadString('\n')
	user.ContactNo = strings.TrimSpace(user.ContactNo)

	// Dob
	for {
		fmt.Print("Date of Birth (YYYY-MM-DD): ")
		dobInput, _ := reader.ReadString('\n') // Read the input as a string
		dobInput = strings.TrimSpace(dobInput)
		user.Dob, err = time.Parse("2006-01-02", dobInput)
		if err != nil {
			fmt.Println("Date format is invalid. Please use YYYY-MM-DD.")
			continue // user can input if format was prev invalid
		}
		break // break loop when user input valid date
	}


	// Address
	fmt.Print("Address: ")
	user.Address, _ = reader.ReadString('\n')
	user.Address = strings.TrimSpace(user.Address)

	// CreatedDateTime
	user.CreatedDateTime = time.Now()

	// insert data into UserService table
	query := `
	INSERT INTO UserService
	(UserId, Name, Email, Password, ContactNo, Dob, Address, CreatedDateTime)
	VALUES(?,?,?,?,?,?,?,?)
	`
	result, err := db.Exec(query, user.UserId, user.Name, user.Email, user.Password, user.ContactNo, user.Dob, user.Address, user.CreatedDateTime)

	if err != nil {
		panic(err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Number of rows affected: ", rows)

}

// User login to their account - POST
func Login(db *sql.DB) string {
	reader := bufio.NewReader(os.Stdin)

	var userLog model.UserService

	// Prompt for Email
	fmt.Print("Email: ")
	userLog.Email, _ = reader.ReadString('\n')
	userLog.Email = strings.TrimSpace(userLog.Email)

	// Prompt for Password
	fmt.Print("Password: ")
	userLog.Password, _ = reader.ReadString('\n')
	userLog.Password = strings.TrimSpace(userLog.Password)

	// query to fetch for hashed password according to given email
	query := `
	SELECT UserId, Password FROM UserService
	WHERE Email = ? 
	`

	// var that holds the hashed pw retrieved from the db
	var userId, storedHash string

	// execute the query to look for hashed pw n store it in storedHash
	err := db.QueryRow(query, userLog.Email).Scan(&userId, &storedHash)
	if err != nil {
		// when no matching row is found
		if err == sql.ErrNoRows{
			fmt.Println("Invalid email or password.")
			return "" // return empty string for failed login
		}
		fmt.Println("Error trying to query database ", err)
		return "" // return empty string for failed database query
	}
	
	// check if pw matches the one in the db
	if !CheckPasswordHash(userLog.Password, storedHash){
		fmt.Println("Invalid email or password")
		return "" // return empty string if pw doesnt match
	}

	// successful login
	fmt.Println("Login successful")
	return userId

}

// User update their account details - PUT
func UpdateUserDetails(db * sql.DB, userId string) {
	reader := bufio.NewReader(os.Stdin)

	// retrieve current details for user
	query := `
	SELECT Name, Email, ContactNo, Address
	FROM UserService
	WHERE UserId = ?
	`

	// current user details
	var currentName, currentEmail, currentContactNo, currentAddress string

	// execute the query to look for current user details
	err := db.QueryRow(query, userId).Scan(&currentName, &currentEmail, &currentContactNo, &currentAddress)
	if err != nil {
		// when no matching row is found
		if err == sql.ErrNoRows{
			fmt.Println("User not found")
			return
		}
		fmt.Println("Error fetching user details ", err)
		return
	}

	fmt.Println("Update Your Details (Press enter to skip)")

	// prompt for Name
	fmt.Printf("Name [%s]: ", currentName)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)

	// no input - keep current value
	if newName == "" {
		newName = currentName
	}

	// prompt for Email
	var newEmail string
	for {
		fmt.Printf("Email [%s]: ", currentEmail)
		newEmail, _ = reader.ReadString('\n')
		newEmail = strings.TrimSpace(newEmail)

		// no input - keep current value
		if newEmail == "" {
			newEmail = currentEmail
			break
		}

		// if new email input does not contain "@"
		if !strings.Contains(newEmail, "@") {
			fmt.Println("Invalid email format. Please try again.")
			continue
		} 

		// check if email is in use for other users
		emailCheckQuery := `
		SELECT COUNT(*)
		FROM UserService
		WHERE Email = ? AND UserId != ? 
		`

		// email count - no.of same emails
		var emailCount int

		// execute the query to look for the number of same emails n store it in emailCount
		err := db.QueryRow(emailCheckQuery, newEmail, userId).Scan(&emailCount)
		if err != nil {
			fmt.Println("Error checking email: ", err)
			return
		}

		// more than one same email
		if emailCount > 0 {
			fmt.Println("Email is already in use. Please try another email.")
			continue
		}

		// email is valid 
		break

	}

	// prompt for ContactNo
	fmt.Printf("Contact Number [%s]: ", currentContactNo)
	newContactNo, _ := reader.ReadString('\n')
	newContactNo = strings.TrimSpace(newContactNo)

	// no input - keep current value
	if newContactNo == "" {
		newContactNo = currentContactNo
	}

	// prompt for Address
	fmt.Printf("Address [%s]: ", currentAddress)
	newAddress, _ := reader.ReadString('\n')
	newAddress = strings.TrimSpace(newAddress)

	// no input - keep current value
	if newAddress == "" {
		newAddress = currentAddress
	}

	// update database 
	updateQuery := `
	UPDATE UserService 
	SET Name = ?, Email = ?, ContactNo = ?, Address = ?
	WHERE UserId = ?
	`

	_, err = db.Exec(updateQuery, newName, newEmail, newContactNo, newAddress, userId)

	// error updating user details
	if err != nil {
		fmt.Println("Error updating user details ", err)
		return
	}

	// successful update
	fmt.Println("User details updated successfully.")

}