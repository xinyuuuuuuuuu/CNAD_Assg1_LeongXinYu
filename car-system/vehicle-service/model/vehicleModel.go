package model

import "time"

// struct for Vehicle
type Vehicle struct {
	VehicleId int `json:"vehicle id"`
	VehicleMake string `json:"vehicle make"`
	VehicleModel string `json:"vehicle model"`
	VehicleType string `json:"vehicle type"`
	LicensePlate string `json:"license plate"`
	VehicleStatus string `json:"vehicle status"`
	VehicleLocation string `json:"vehicle location"`
	VehicleChargeLevel uint8 `json:"vehicle charge level"`
	VehicleCleanliness string `json:"vehicle cleanliness"`
}

// struct for Reservation
type Reservation struct {
	ReservationId int `json:"reservation id"`
	UserId int `json:"user id"`
	VehicleId int `json:"vehicle id"`
	ReserveStatus string `json:"reserve status"`
	ReserveStartDate time.Time `json:"reserve start date"`
	ReserveEndDate time.Time `json:"reserve end date"`
	EstimatedTotalCost float64 `json:"estimated total cost"`
	CreatedDate time.Time `json:"created date"`
	ModifiedDate *time.Time `json:"modified date"`
}