package model

import (
	"time"
)

type Reservation struct {
	ReservationId string
	UserId string
	VehicleId string
	ReserveStatus string `json:"Reservation Status"`
	ReserveStartDate time.Time `json:"Reservation Start Date"`
	ReserveEndDate time.Time `json:"Reservation End Date"`
	EstimatedTotalCost float64 `json:"Estimated Total Cost"`
	CreatedDate time.Time 
	ModifiedDate time.Time
}

