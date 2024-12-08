package model

import (
	"time"
)

type Membership struct {
	MembershipId                string
	UserId                      string
	MembershipTier              string    `json:"Membership Tier"`
	HourlyRate                  float64   `json:"Hourly Rate"`
	MemberDiscount              float64   `json:"Member Discount(%)"`
	PiorityLevel                int       `json:"Priority Level (0-2)"`
	TotalCostPerMonth           float64   `json:"Total Cost Per Month"`
	MembershipExpiryDate        time.Time `json:"Membership Expiry Date"`
	EligibleForUpgradeNextMonth string    `json:"Eligible to upgrade membership tier"`
}
