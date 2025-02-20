package model

import "time"

// struct for UserService
type UserService struct {
	UserId          int       `json:"user_id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	ContactNo       string    `json:"contact_no"`
	Dob             string    `json:"dob"`
	Address         string    `json:"address"`
	HashedPassword  string    `json:"hashed_password"`
	MembershipLevel string    `json:"membership_level"` // Added MembershipLevel
	CreatedDateTime time.Time `json:"created_datetime"`
}
