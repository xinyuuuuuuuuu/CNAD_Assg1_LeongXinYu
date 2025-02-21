package vehicle

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Vehicle struct {
	VehicleId          int    `json:"vehicle id"`
	VehicleMake        string `json:"vehicle make"`
	VehicleModel       string `json:"vehicle model"`
	VehicleType        string `json:"vehicle type"`
	LicensePlate       string `json:"license plate"`
	VehicleStatus      string `json:"vehicle status"`
	VehicleLocation    string `json:"vehicle location"`
	VehicleChargeLevel uint8  `json:"vehicle charge level"`
	VehicleCleanliness string `json:"vehicle cleanliness"`
}

// ViewAvailableVehicles returns a list of available vehicles
func ViewAvailableVehicles(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] ViewAvailableVehicles: Received request")
	w.Header().Set("Content-Type", "application/json")

	log.Println("[DEBUG] Executing query to fetch available vehicles")
	query := `SELECT VehicleID, VehicleMake, VehicleModel, VehicleType, LicensePlate, VehicleStatus, VehicleLocation, VehicleChargeLevel, VehicleCleanliness FROM Vehicle WHERE VehicleStatus = 'Available'`
	rows, err := db.Query(query)
	if err != nil {
		log.Println("[ERROR] Database query error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []Vehicle
	log.Println("[DEBUG] Iterating over query results")
	for rows.Next() {
		var vehicle Vehicle
		if err := rows.Scan(&vehicle.VehicleId, &vehicle.VehicleMake, &vehicle.VehicleModel, &vehicle.VehicleType, &vehicle.LicensePlate, &vehicle.VehicleStatus, &vehicle.VehicleLocation, &vehicle.VehicleChargeLevel, &vehicle.VehicleCleanliness); err != nil {
			log.Println("[ERROR] Error scanning vehicle row:", err)
			http.Error(w, "Error retrieving vehicles", http.StatusInternalServerError)
			return
		}
		log.Printf("[DEBUG] Retrieved vehicle: %+v\n", vehicle)
		vehicles = append(vehicles, vehicle)
	}

	if len(vehicles) == 0 {
		log.Println("[DEBUG] No available vehicles found")
	} else {
		log.Printf("[DEBUG] Successfully retrieved %d available vehicles\n", len(vehicles))
	}

	log.Println("[DEBUG] Sending response to client")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

// RegisterRoutes sets up the reservation and vehicle routes
func RegisterRoutes(router *mux.Router, db *sql.DB) {
	log.Println("[DEBUG] Registering vehicle routes")
	router.HandleFunc("/vehicles", func(w http.ResponseWriter, r *http.Request) { 
		log.Println("[DEBUG] Handling /vehicles request")
		ViewAvailableVehicles(db, w, r) 
	}).Methods("GET")
}
