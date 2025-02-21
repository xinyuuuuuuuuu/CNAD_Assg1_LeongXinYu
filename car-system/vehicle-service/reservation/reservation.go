package reservation

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"user-management/middleware"

	"github.com/gorilla/mux"
)

// struct for Reservation
type Reservation struct {
	ReservationId      int        `json:"reservation_id"`
	UserId             int        `json:"user_id"`
	VehicleId          int        `json:"vehicle_id"`
	ReserveStatus      string     `json:"reserve_status"`
	ReserveStartDate   time.Time  `json:"reserve_start_date"`
	ReserveEndDate     time.Time  `json:"reserve_end_date"`
	EstimatedTotalCost float64    `json:"estimated_total_cost"`
	CreatedDate        time.Time  `json:"created_date"`
	ModifiedDate       *time.Time `json:"modified_date"`
}

// GetReservations retrieves reservations for logged-in user
func GetReservations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] GetReservations: Received request")

	// Extract userId from request context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		log.Println("[ERROR] Unauthorized - Missing userId in request context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("[DEBUG] Fetching reservations for userId: %d\n", userId)

	// Updated query to fetch reservations
	query := `SELECT ReservationId, UserId, VehicleId, ReserveStatus, ReserveStartDate, ReserveEndDate, EstimatedTotalCost, CreatedDate, ModifiedDate FROM Reservation WHERE UserId = ?`
	rows, err := db.Query(query, userId)
	if err != nil {
		log.Println("[ERROR] Database query error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var reservation Reservation
		var reserveStartDate, reserveEndDate, createdDate, modifiedDate sql.NullString

		// Use sql.NullString to handle possible NULL values and avoid conversion errors
		if err := rows.Scan(
			&reservation.ReservationId,
			&reservation.UserId,
			&reservation.VehicleId,
			&reservation.ReserveStatus,
			&reserveStartDate,
			&reserveEndDate,
			&reservation.EstimatedTotalCost,
			&createdDate,
			&modifiedDate,
		); err != nil {
			log.Println("[ERROR] Error scanning reservation row:", err)
			http.Error(w, "Error retrieving reservations", http.StatusInternalServerError)
			return
		}

		// Convert date strings to time.Time
		const timeLayout = "2006-01-02 15:04:05" // MySQL DATETIME format

		if reserveStartDate.Valid {
			reservation.ReserveStartDate, err = time.Parse(timeLayout, reserveStartDate.String)
			if err != nil {
				log.Println("[ERROR] Error parsing ReserveStartDate:", err)
			}
		}

		if reserveEndDate.Valid {
			reservation.ReserveEndDate, err = time.Parse(timeLayout, reserveEndDate.String)
			if err != nil {
				log.Println("[ERROR] Error parsing ReserveEndDate:", err)
			}
		}

		if createdDate.Valid {
			reservation.CreatedDate, err = time.Parse(timeLayout, createdDate.String)
			if err != nil {
				log.Println("[ERROR] Error parsing CreatedDate:", err)
			}
		}

		if modifiedDate.Valid {
			parsedTime, err := time.Parse(timeLayout, modifiedDate.String)
			if err == nil {
				reservation.ModifiedDate = &parsedTime
			} else {
				log.Println("[ERROR] Error parsing ModifiedDate:", err)
			}
		}

		log.Printf("[DEBUG] Retrieved reservation: %+v\n", reservation)
		reservations = append(reservations, reservation)
	}

	log.Printf("[DEBUG] Total reservations found: %d\n", len(reservations))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
}

// MakeReservation allows users to reserve a vehicle
func MakeReservation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] MakeReservation: Received request")

	// Extract userId from JWT middleware
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		log.Println("[ERROR] Unauthorized - Missing userId in request context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("[DEBUG] Making reservation for userId: %d\n", userId)

	w.Header().Set("Content-Type", "application/json")

	var reservation struct {
		VehicleId        int       `json:"vehicle_id"`
		ReserveStartDate time.Time `json:"reserve_start_date"`
		ReserveEndDate   time.Time `json:"reserve_end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		log.Println("[ERROR] Invalid request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure start and end dates are valid
	if reservation.ReserveStartDate.After(reservation.ReserveEndDate) {
		log.Println("[ERROR] Invalid reservation dates")
		http.Error(w, "Reserve start date must be before end date", http.StatusBadRequest)
		return
	}

	// Fetch user's membership level from `user_management` database
	var membershipLevel string
	err := db.QueryRow("SELECT MembershipLevel FROM user_management.UserService WHERE UserId = ?", userId).Scan(&membershipLevel)
	if err != nil {
		log.Println("[ERROR] Failed to fetch user membership level:", err)
		http.Error(w, "Failed to fetch user membership level", http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG] User membership level: %s\n", membershipLevel)

	// Assign fixed cost based on membership level
	var estimatedTotalCost float64
	switch membershipLevel {
	case "Premium":
		estimatedTotalCost = 40.00
	case "VIP":
		estimatedTotalCost = 30.00
	default: // Basic membership
		estimatedTotalCost = 50.00
	}
	log.Printf("[DEBUG] Assigned Estimated Total Cost: %.2f\n", estimatedTotalCost)

	// Determine vehicle availability to set `reserve_status`
	var vehicleStatus string
	err = db.QueryRow("SELECT VehicleStatus FROM Vehicle WHERE VehicleId = ?", reservation.VehicleId).Scan(&vehicleStatus)
	if err != nil {
		log.Println("[ERROR] Failed to fetch vehicle status:", err)
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}
	log.Printf("[DEBUG] Vehicle Status: %s\n", vehicleStatus)

	// Assign `reserve_status` based on vehicle availability
	reserveStatus := "Pending"
	if vehicleStatus == "Available" {
		reserveStatus = "Active"
	}
	log.Printf("[DEBUG] Assigned Reserve Status: %s\n", reserveStatus)

	// Insert reservation into the database
	query := `
		INSERT INTO Reservation (UserId, VehicleId, ReserveStatus, ReserveStartDate, ReserveEndDate, EstimatedTotalCost) 
		VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, userId, reservation.VehicleId, reserveStatus, reservation.ReserveStartDate, reservation.ReserveEndDate, estimatedTotalCost)
	if err != nil {
		log.Println("[ERROR] Error inserting reservation into database:", err)
		http.Error(w, "Failed to make reservation", http.StatusInternalServerError)
		return
	}

	// Retrieve last inserted ID
	reservationID, _ := result.LastInsertId()
	log.Printf("[DEBUG] Reservation created successfully with ID: %d\n", reservationID)

	// Send response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":              "Reservation created successfully",
		"reservation_id":       reservationID,
		"reserve_status":       reserveStatus,
		"estimated_total_cost": estimatedTotalCost, // Return computed cost in response
	})
}

// UpdateReservation allows users to update their reservation details
func UpdateReservation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] UpdateReservation: Received request")

	w.Header().Set("Content-Type", "application/json")

	var reservation Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		log.Println("[ERROR] Invalid request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] Updating reservation ID: %d\n", reservation.ReservationId)
	query := `UPDATE Reservation SET ReserveStatus = ?, ReserveStartDate = ?, ReserveEndDate = ?, EstimatedTotalCost = ?, ModifiedDate = CURRENT_TIMESTAMP WHERE ReservationId = ?`
	_, err := db.Exec(query, reservation.ReserveStatus, reservation.ReserveStartDate, reservation.ReserveEndDate, reservation.EstimatedTotalCost, reservation.ReservationId)
	if err != nil {
		log.Println("[ERROR] Error updating reservation:", err)
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}

	log.Println("[DEBUG] Reservation updated successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation updated successfully"})
}

// DeleteReservation allows users to cancel a reservation
func DeleteReservation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] DeleteReservation: Received request")

	w.Header().Set("Content-Type", "application/json")

	var request struct {
		ReservationId int `json:"reservation_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("[ERROR] Invalid request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] Deleting reservation ID: %d\n", request.ReservationId)
	query := `DELETE FROM Reservation WHERE ReservationId = ?`
	_, err := db.Exec(query, request.ReservationId)
	if err != nil {
		log.Println("[ERROR] Error deleting reservation:", err)
		http.Error(w, "Failed to delete reservation", http.StatusInternalServerError)
		return
	}

	log.Println("[DEBUG] Reservation deleted successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation deleted successfully"})
}

// RegisterRoutes sets up the reservation routes with JWT validation
func RegisterRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/get-reservations", middleware.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		GetReservations(db, w, r)
	})).Methods("GET")

	router.HandleFunc("/make-reservations", middleware.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		MakeReservation(db, w, r)
	})).Methods("POST")

	router.HandleFunc("/update-reservations", middleware.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		UpdateReservation(db, w, r)
	})).Methods("PUT")

	router.HandleFunc("/delete-reservations", middleware.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		DeleteReservation(db, w, r)
	})).Methods("DELETE")
}
