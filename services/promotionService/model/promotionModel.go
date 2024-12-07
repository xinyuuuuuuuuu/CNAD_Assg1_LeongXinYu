package model

import (
	"time"
)

type Promotion struct {
	PromoId string
	PromoType string `json:"Promotion Type"`
	PromoDiscount float64 `json:"Promotion Discount"`
	PromoStartDate time.Time `json:"Promotion Start Date"`
	PromoEndDate time.Time `json:"Promotion End Date"`
}