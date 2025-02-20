package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

func Signup(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")
	log.Println("[DEBUG] Signup handler invoked")

	var user model.UserService
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("[ERROR] Error parsing JSON:", err)
		return
	}

	// Debug received struct
	log.Printf("[DEBUG] Received struct: %+v\n", user)

	// Validate required fields
	missingFields := []string{}
	if user.Name == "" {
		missingFields = append(missingFields, "Name")
	}
	if user.Email == "" {
		missingFields = append(missingFields, "Email")
	}
	if user.ContactNo == "" {
		missingFields = append(missingFields, "ContactNo")
	}
	if user.Dob == "" {
		missingFields = append(missingFields, "Dob")
	}
	if user.Address == "" {
		missingFields = append(missingFields, "Address")
	}
	if user.HashedPassword == "" {
		missingFields = append(missingFields, "HashedPassword")
	}

	// If any fields are missing, return an error
	if len(missingFields) > 0 {
		http.Error(w, "Missing required fields: "+strings.Join(missingFields, ", "), http.StatusBadRequest)
		log.Println("[ERROR] Missing required fields:", strings.Join(missingFields, ", "))
		return
	}

	// Validate DOB
	parsedDob, err := time.Parse("2006-01-02", user.Dob)
	if err != nil {
		http.Error(w, "Invalid Date of Birth format. Use YYYY-MM-DD", http.StatusBadRequest)
		log.Println("[ERROR] Error parsing DOB:", err)
		return
	}

	// Assign parsedDob to user obj
	userTime := parsedDob.UTC()
	log.Printf("[DEBUG] Parsed DOB (UTC): %s\n", userTime)

	// Hash password
	hashPassword, err := HashedPassword(strings.TrimSpace(user.HashedPassword))
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		log.Println("[ERROR] Error hashing password:", err)
		return
	}
	user.HashedPassword = hashPassword
	log.Println("[DEBUG] Password hashed successfully")

	// Set created timestamp
	user.CreatedDateTime = time.Now().UTC()
	log.Printf("[DEBUG] CreatedDateTime (UTC): %s\n", user.CreatedDateTime)

	// Insert data into UserService table
	query := `
	INSERT INTO UserService
	(Name, Email, ContactNo, Dob, Address, HashedPassword, CreatedDateTime)
	VALUES(?,?,?,?,?,?,?)
	`

	log.Println("[DEBUG] Executing INSERT query")
	result, err := db.Exec(query, user.Name, user.Email, user.ContactNo, userTime, user.Address, user.HashedPassword, user.CreatedDateTime)
	if err != nil {
		log.Println("[ERROR] Error inserting into database:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	log.Println("[DEBUG] Insert query executed successfully")

	// Retrieve the last inserted user ID
	userId, err := result.LastInsertId()
	if err != nil {
		log.Println("[ERROR] Error getting last inserted ID:", err)
		http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		return
	}

	// Convert userId to int
	user.UserId = int(userId)
	log.Printf("[DEBUG] User created with ID: %d\n", user.UserId)

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User account created successfully!",
		"userId":  user.UserId,
	})
	log.Println("[INFO] User signup successful")
}

// POST - user login into account
func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// response header
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("[DEBUG] Login attempt received") // Debug statement

	var user model.UserService
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("[ERROR] Error decoding request body:", err) // Debug statement
		http.Error(w, `{"message": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	fmt.Printf("[DEBUG] Decoded request body: %+v\n", user) // Debug statement

	// Ensure email is not empty
	if user.Email == "" {
		fmt.Println("[ERROR] Email is empty in request body")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Ensure password is not empty
	if user.HashedPassword == "" {
		fmt.Println("[ERROR] Password is empty in request body")
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Query to fetch hashed password based on email
	query := `
	SELECT UserId, HashedPassword FROM UserService
	WHERE Email = ? 
	`

	var hashedpw string
	var userId int

	// Execute the query
	err = db.QueryRow(query, user.Email).Scan(&userId, &hashedpw)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("[ERROR] No matching user found for email:", user.Email) // Debug statement
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		fmt.Println("[ERROR] Database query error:", err) // Debug statement
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	fmt.Println("[DEBUG] User found with ID:", userId)                                  // Debug statement
	fmt.Println("[DEBUG] Fetched hashed password from DB:", hashedpw)                   // Debug statement
	fmt.Println("[DEBUG] Received raw password (before hashing):", user.HashedPassword) // Debug statement

	// Validate if the fetched hashed password matches the entered password
	match := CheckPasswordHash(user.HashedPassword, hashedpw)
	if !match {
		fmt.Println("[ERROR] Password mismatch for user ID:", userId)         // Debug statement
		fmt.Println("[DEBUG] Hashed password does not match the stored hash") // Debug statement
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	fmt.Println("[DEBUG] Password matched successfully for user ID:", userId) // Debug statement

	// Generate JWT token
	token, err := GenerateJWT(userId)
	if err != nil {
		fmt.Println("[ERROR] Error generating JWT token:", err) // Debug statement
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	fmt.Println("[INFO] Login successful for user ID:", userId) // Debug statement

	// Successful login
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Token":   token,
		"message": "Login successful",
	})
}

// GET - user view account details
func ViewAccountDetails(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// response header
	w.Header().Set("Content-Type", "application/json")

	log.Println("ViewAccountDetails: Received request")

	// retrieve userId from context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		log.Println("ViewAccountDetails: Unauthorized access - missing userId in context")
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	log.Printf("ViewAccountDetails: Retrieved userId: %d\n", userId)

	// query to view account details, including MembershipLevel
	query := `
	SELECT UserId, Name, Email, ContactNo, Dob, Address, HashedPassword, MembershipLevel
	FROM UserService
	WHERE UserId = ?`

	var user model.UserService

	// execute the query
	err := db.QueryRow(query, userId).Scan(
		&user.UserId,
		&user.Name,
		&user.Email,
		&user.ContactNo,
		&user.Dob,
		&user.Address,
		&user.HashedPassword,
		&user.MembershipLevel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ViewAccountDetails: User not found for userId: %d\n", userId)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("ViewAccountDetails: Database error: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("ViewAccountDetails: Successfully retrieved user details: %+v\n", user)

	// return user details in json
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("ViewAccountDetails: Error encoding JSON response: %v\n", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Println("ViewAccountDetails: Response sent successfully")
}

// PUT - update user account details
func UpdateAccountDetails(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// response header
	w.Header().Set("Content-Type", "application/json")

	log.Println("UpdateAccountDetails: Received request")

	// retrieve userId from context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		log.Println("UpdateAccountDetails: Unauthorized access - missing userId in context")
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	log.Printf("UpdateAccountDetails: Retrieved userId: %d\n", userId)

	// parse request body
	var updatedUser model.UserService
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Printf("UpdateAccountDetails: Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// update query
	query := `
	UPDATE UserService 
	SET Name = ?, Email = ?, ContactNo = ?, Address = ?
	WHERE UserId = ?
	`

	// execute update query
	_, err = db.Exec(query, updatedUser.Name, updatedUser.Email, updatedUser.ContactNo, updatedUser.Address, userId)
	if err != nil {
		log.Printf("UpdateAccountDetails: Database error: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Println("UpdateAccountDetails: User details updated successfully")

	// send success response
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Account details updated successfully"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("UpdateAccountDetails: Error encoding JSON response: %v\n", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Println("UpdateAccountDetails: Response sent successfully")
}

// UpdateMembershipLevel updates the membership level of a user
func UpdateMembershipLevel(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request struct {
		UserID          int    `json:"user_id"`
		MembershipLevel string `json:"membership_level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate membership level
	validLevels := map[string]bool{"Basic": true, "Premium": true, "VIP": true}
	if !validLevels[request.MembershipLevel] {
		http.Error(w, "Invalid membership level", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement
	query := "UPDATE UserService SET MembershipLevel = ? WHERE UserId = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error preparing query: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the update query
	_, err = stmt.Exec(request.MembershipLevel, request.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating membership level: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Successfully updated user %d to membership level: %s", request.UserID, request.MembershipLevel)})
}
