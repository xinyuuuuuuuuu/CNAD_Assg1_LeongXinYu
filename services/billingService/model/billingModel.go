package model

import (
	"time"
)

type UserService struct {
	BillingId        string
	UserId           string
	BillingDate      time.Time `json:"Billing Date"`
	BillingTotal     float64   `json:"Billing Total"`
	PaymentMethod    string    `json:"Payment Method"`
	PaymentStatus    string    `json:"Payment Status"`
	LastActivityDate time.Time
}
