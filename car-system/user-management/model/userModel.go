package model

import "time"

// struct for UserService
type UserService struct {
	UserId int `json:"user id"`
	Name string `json:"name"`
	Email string `json:"email"`
	ContactNo string `json:"contact no"`
	Dob time.Time `json:"dob"`
	Address string `json:"address"`
	HashedPassword string `json:"hashed password"`
	CreatedDateTime time.Time `json:"created date time"`
}

// struct for Membership
type Membership struct {
	MembershipId int `json:"membership id"`
	UserId int `json:"user id"`
	MembershipTier string `json:"membership tier"`
	HourlyRate float64 `json:"hourly rate"`
	MemberDiscount float64 `json:"member discount"`
	PriorityLevel int `json:"priority level"`
	TotalCostPerMonth float64 `json:"total cost per month"`
	MembershipExpiryDate time.Time `json:"membership expiry date"`	
	EligibleForUpgradeNextMonth bool `json:"eligible for upgrade next month"`
}