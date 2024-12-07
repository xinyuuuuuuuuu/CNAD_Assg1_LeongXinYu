package controller

import (
	"fmt"
	"time"
	"cnad_assg1_leongxinyu/services/userService/model"
	"cnad_assg1_leongxinyu/services/user-service/utility"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

// hash password
func HashedPassword(password string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// check password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) 
	return err == nil
}

func Signup (db *sql.DB) {
	var user model.UserService
	var userId string

	// userid 
	userId, err = utility.GenerateUserId(db)
	if err != nil {
		fmt.Println("Error generating user id: ", err)
		return
	}
	user.UserId = userId

	// Name
	fmt.Print("Name: ")
	fmt.Scanln(&user.Name)

	// Email
	fmt.Print("Email: ")
	fmt.Scanln(&user.Email)

	// Password
	fmt.Print("Password: ")
	var pw string
	fmt.Scanln(&pw)
	hash, err := HashedPassword(pw)
	if err != nil {
		fmt.Println("Error hashing password ", err)
	}

	// store hashed password
	user.Password = hash 
	fmt.Println("Hash: ", user.Password)

	// ContactNo
	fmt.Print("Contact Number: ")
	fmt.Scanln(&user.ContactNo)

	// Dob
	fmt.Print("Date of Birth (YYYY-MM-DD): ")
	fmt.Scanln(&user.Dob)

	// Address
	fmt.Print("Address: ")
	fmt.Scanln(&user.Address)

	// CreatedDateTime
	user.CreatedDateTime = time.Now()


}
