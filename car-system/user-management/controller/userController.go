package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"user-management/config"
	"user-management/model"

	"github.com/golang-jwt/jwt"

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

// struct to return token response
type TokenResponse struct {
	Token   string `json:"Token"`
	Message string `json:"Message"`
}

// function to create jwt token
func GenerateJWT(userId int) (string, error) {
	// define expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// create jwt claims
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    expirationTime.Unix(),
	}

	// create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the token with the secret key
	secretKey := config.GetJWTSecret()
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// POST - user signs up for an account
func Signup(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// response header
	w.Header().Set("Content-Type", "application/json")

	var user model.UserService
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// hashed password
	hashpassword, err := HashedPassword(user.HashedPassword)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.HashedPassword = hashpassword

	// set CreatedDateTime
	user.CreatedDateTime = time.Now()

	// insert data into UserService table
	query := `
	INSERT INTO UserService
	(Name, Email, ContactNo, Dob, Address, HashedPassword, CreatedDateTime)
	VALUES(?,?,?,?,?,?,?)
	`
	_, err = db.Exec(query, user.Name, user.Email, user.ContactNo, user.Dob, user.Address, user.HashedPassword, user.CreatedDateTime)
	if err != nil {
		log.Println("Error inserting into database:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// getting the most recent user id
	// userId, err := result.LastInsertId()
	// if err != nil {
	// 	log.Println("Error getting last inserted id:", err)
	// 	http.Error(w, "Failed to retrieve user Id", http.StatusInternalServerError)
	// 	return
	// }

	// return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User account created successfully!",
	})
}

// POST - user login into account
func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// response header
	w.Header().Set("Content-Type", "application/json")

	var user model.UserService
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"message": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// query to fetch for hashed password according to given email
	query := `
	SELECT UserId, HashedPassword FROM UserService
	WHERE Email = ? 
	`

	// var that holds the hashed pw retrieved from the db
	var hashedpw string
	var userId int

	// execute the query
	err = db.QueryRow(query, user.Email).Scan(&userId, &hashedpw)
	if err != nil {
		// when no matching row is found
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// check if pw matches the hashedpw in the db
	if !CheckPasswordHash(user.HashedPassword, hashedpw) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// generate jwt token
	token, err := GenerateJWT(userId)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// successful login
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Token":   token,
		"message": "Login sucessful",
	})
}

// GET - user view account details
func ViewAccountDetails(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// response header
	w.Header().Set("Content-Type", "application/json")

	// retrieve userId from context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// query to view account details
	query := `
	SELECT Name, Email, ContactNo, Dob, Address
	FROM UserService
	WHERE UserId = ?
	`
	var user model.UserService

	// execute the query
	err := db.QueryRow(query, userId).Scan(&user.Name, &user.Email, &user.ContactNo, &user.Dob, &user.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// close the result when the func has ended
	// defer results.Close()

	// return user details in json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// PUT - User update their account details
func UpdateAccountDetails(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// response header
	w.Header().Set("Content-Type", "application/json")

	// retrieve userId from context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// parse incoming json req body
	var updatedUser model.UserService
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// retrieve current details for user
	query := `
	SELECT Name, Email, ContactNo, Address
	FROM UserService
	WHERE UserId = ?
	`
	var currentUser model.UserService

	// execute the query to look for current user details
	err = db.QueryRow(query, userId).Scan(&currentUser.Name, &currentUser.Email, &currentUser.ContactNo, &currentUser.Address)
	if err != nil {
		// when no matching row is found
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// keep original value if the field is empty
	// Name
	if updatedUser.Name == "" {
		updatedUser.Name = currentUser.Name
	}

	// Email
	if updatedUser.Email == "" {
		updatedUser.Email = currentUser.Email
	}

	// validate that updated email is not in used by another user
	if updatedUser.Email != currentUser.Email {
		emailCheckQuery := `
		SELECT COUNT(*) 
		FROM UserService
		WHERE Email = ? AND UserId != ?
		`
		var emailCount int
		err := db.QueryRow(emailCheckQuery, updatedUser.Email, userId).Scan(&emailCount)
		if err != nil {
			http.Error(w, "Error checking email", http.StatusInternalServerError)
			return
		}
		if emailCount > 0 {
			http.Error(w, "Email is already in use", http.StatusConflict)
			return
		}
	}

	// update database
	updateQuery := `
	UPDATE UserService 
	SET Name = ?, Email = ?, ContactNo = ?, Address = ?
	WHERE UserId = ?
	`

	_, err = db.Exec(updateQuery,updatedUser.Name, updatedUser.Email, updatedUser.ContactNo, updatedUser.Address, userId)

	// error updating user details
	if err != nil {
		http.Error(w, "Error updating user details ", http.StatusInternalServerError)
		return
	}

	// successful update
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string){
		"message": "User details updated successfully!"
	}
}
