package model

type Vehicle struct {
	VehicleId          string
	VehicleMake        string `json:"Vehicle Make"`
	VehicleModel       string `json:"Vehicle Model"`
	VehicleType        string `json:"Vehicle Type"`
	LicensePlate       string `json:"License Plate"`
	VehicleStatus      string `json:"Vehicle Status"`
	VehicleLocation    string `json:"Vehicle Location"`
	VehicleChargeLevel int64  `json:"Vehicle Charge Level"`
	VehicleCleanliness string `json:"Vehhicle Cleanliness"`
}
