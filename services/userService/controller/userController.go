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

// User sign up for an account
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
	fmt.Print("Email: ")
	user.Email, _ = reader.ReadString('\n')
	user.Email = strings.TrimSpace(user.Email)

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

// User login to their account
func Login(db *sql.DB) {
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
	SELECT Password FROM UserService
	WHERE Email = ? 
	`

	// var that holds the hashed pw retrieved from the db
	var storedHash string

	// execute the query to look for hashed pw n store it in storedHash
	err := db.QueryRow(query, userLog.Email).Scan(&storedHash)
	if err != nil {
		// when no matching row is found
		if err == sql.ErrNoRows{
			fmt.Println("Invalid email or password.")
			return
		}
		fmt.Println("Error trying to query database ", err)
		return
	}
	
	// check if pw matches the one in the db
	if !CheckPasswordHash(userLog.Password, storedHash){
		fmt.Println("Invalid email or password")
	}

	// successful login
	fmt.Println("Login successful")

}
