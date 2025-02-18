package model

import "time"

// struct for Billing
type Billing struct {
	BillingId int `json:"billing id"`
	UserId int `json:"user id"`
	ReservationId int `json:"reservation id"`
	BillingTotal float64 `json:"billing total"`
	MembershipDiscount float64 `json:"membership discount"`
	PromoDiscount float64 `json:"promo discount"`
	BillingFinal float64 `json:"billing final"`
	PaymentStatus string `json:"payment status"`
	BillingDate time.Time `json:"billing date"`
}

// struct for Promotion
type Promotion struct {
	PromoId int `json:"promo id"`
	PromoType string `json:"promo type"`
	PromoDiscountPercentage float64 `json:"promo discount percentage"`
	PromoStartDate time.Time `json:"promo start date"`
	PromoEndDate time.Time `json:"promo end date"`
}

// struct for PaymentTransactions
type PaymentTransactions struct {
	TransactionId int `json:"transaction id"`
	BillingId int `json:"billing id"`
	UserId int `json:"user id"`
	PaymentMethod string `json:"payment method"`
	PaymentStatus string `json:"payment status"`
	TransactionAmount float64 `json:"transaction amount"`
	TransactionDate time.Time `json:"transaction date"`
}


