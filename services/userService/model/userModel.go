package model

import (
	"time"
)

type UserService struct {
	UserId string 
	Name   string `json:"Name"`
	Email   string `json:"Email"`
	ContactNo   string `json:"Contact Number"`
	Dob  time.Time  `json:"Date Of Birth"`	
	Address   string `json:"Address"`
	Password string    
	CreatedDateTime time.Time    
}

