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
	VehicleMake        string     `json:"vehicle_make"`
	VehicleModel       string     `json:"vehicle_model"`
	VehicleType        string     `json:"vehicle_type"`
	LicensePlate       string     `json:"license_plate"`
}

// GetReservations retrieves active and pending reservations with vehicle details for logged-in user
func GetReservations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] GetReservations: Received request")

	// Extract userId from request context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		log.Println("[ERROR] Unauthorized - Missing userId in request context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("[DEBUG] Fetching active and pending reservations for userId: %d\n", userId)

	// Updated query to fetch reservations with vehicle details (make, model, type, license plate)
	query := `SELECT 
                  r.ReservationId, 
                  r.UserId, 
                  r.VehicleId, 
                  r.ReserveStatus, 
                  r.ReserveStartDate, 
                  r.ReserveEndDate, 
                  r.EstimatedTotalCost, 
                  r.CreatedDate, 
                  r.ModifiedDate, 
                  v.VehicleMake, 
                  v.VehicleModel, 
                  v.VehicleType, 
                  v.LicensePlate
              FROM 
                  Reservation r
              INNER JOIN 
                  Vehicle v ON r.VehicleId = v.VehicleId
              WHERE 
                  r.UserId = ? AND (r.ReserveStatus = 'Pending' OR r.ReserveStatus = 'Active')`

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
			&reservation.VehicleMake,  // Added VehicleMake
			&reservation.VehicleModel, // Added VehicleModel
			&reservation.VehicleType,  // Added VehicleType
			&reservation.LicensePlate, // Added LicensePlate
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
	if len(reservations) == 0 {
		// Return a message when no reservations are found
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "No reservations found"})
		return
	}

	// Return the reservations in the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
}

// MakeReservation allows users to create their reservation
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

	// Define structure to capture request body
	var reservation struct {
		VehicleId        int       `json:"vehicle_id"`
		ReserveStartDate time.Time `json:"reserve_start_date"`
		ReserveEndDate   time.Time `json:"reserve_end_date"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		log.Println("[ERROR] Invalid request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Debug received data
	log.Printf("[DEBUG] Reservation data - VehicleId: %d, StartDate: %s, EndDate: %s\n",
		reservation.VehicleId, reservation.ReserveStartDate, reservation.ReserveEndDate)

	// Validate vehicleId before using it
	if reservation.VehicleId <= 0 {
		log.Println("[ERROR] Invalid VehicleId")
		http.Error(w, "Invalid VehicleId", http.StatusBadRequest)
		return
	}

	// Ensure start and end dates are valid
	if reservation.ReserveStartDate.After(reservation.ReserveEndDate) {
		log.Println("[ERROR] Invalid reservation dates")
		http.Error(w, "Reserve start date must be before end date", http.StatusBadRequest)
		return
	}

	// Prevent duplicate reservations for the same user and vehicle
	var existingCount int
	checkQuery := `
		SELECT COUNT(*) FROM Reservation
		WHERE UserId = ? AND VehicleId = ? AND ReserveStatus IN ('Pending', 'Active')`
	err := db.QueryRow(checkQuery, userId, reservation.VehicleId).Scan(&existingCount)

	if err != nil {
		log.Println("[ERROR] Failed to check existing reservation:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingCount > 0 {
		log.Println("[ERROR] User already has a pending or active reservation for this vehicle")
		http.Error(w, "You already have a reservation for this vehicle", http.StatusConflict)
		return
	}

	// Fetch user's membership level from user_management database
	var membershipLevel string
	err = db.QueryRow("SELECT MembershipLevel FROM user_management.UserService WHERE UserId = ?", userId).Scan(&membershipLevel)
	if err != nil {
		log.Println("[ERROR] Failed to fetch user membership level:", err)
		http.Error(w, "Failed to fetch user membership level", http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG] User membership level: %s\n", membershipLevel)

	// Determine base price per hour based on membership level
	var basePrice float64
	switch membershipLevel {
	case "Premium":
		basePrice = 40.00
	case "VIP":
		basePrice = 30.00
	default:
		basePrice = 50.00
	}
	log.Printf("[DEBUG] Base Price per Hour: %.2f\n", basePrice)

	// Calculate number of hours booked
	hoursBooked := reservation.ReserveEndDate.Sub(reservation.ReserveStartDate).Hours()
	log.Printf("[DEBUG] Hours booked: %.2f\n", hoursBooked)

	// Calculate estimated total cost based on hours booked
	billingTotal := basePrice * hoursBooked
	log.Printf("[DEBUG] Calculated Billing Total: %.2f\n", billingTotal)

	// Always set the reservation status to "Pending" at creation
	reserveStatus := "Pending"

	// Insert the new reservation into the database
	insertQuery := `
		INSERT INTO Reservation 
			(UserId, VehicleId, ReserveStatus, ReserveStartDate, ReserveEndDate, EstimatedTotalCost, CreatedDate)
		VALUES 
			(?, ?, ?, ?, ?, ?, NOW())
	`
	formattedStartDate := reservation.ReserveStartDate.Format("2006-01-02 15:04:05")
	formattedEndDate := reservation.ReserveEndDate.Format("2006-01-02 15:04:05")

	result, err := db.Exec(insertQuery, userId, reservation.VehicleId, reserveStatus, formattedStartDate, formattedEndDate, billingTotal)
	if err != nil {
		log.Println("[ERROR] Error inserting reservation:", err)
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}
	reservationId, err := result.LastInsertId()
	if err != nil {
		log.Println("[ERROR] Error retrieving last insert id:", err)
	}

	log.Printf("[DEBUG] Reservation inserted with ID: %d and status: %s\n", reservationId, reserveStatus)

	// Fetch applicable promotions from billing_service
	var promoDiscount float64
	var promoType string
	err = db.QueryRow(`
		SELECT PromoDiscountPercentage, PromoType
		FROM billing_service.Promotion
		WHERE CURDATE() BETWEEN PromoStartDate AND PromoEndDate
		ORDER BY PromoStartDate DESC LIMIT 1
	`).Scan(&promoDiscount, &promoType)

	if err != nil {
		log.Println("[INFO] No active promotion found, proceeding without discount")
		promoDiscount = 0
	}

	// Calculate discount amount
	discountAmount := 0.00
	if promoType == "Percentage Discount" || promoType == "Seasonal Offer" {
		discountAmount = billingTotal * (promoDiscount / 100)
	} else if promoType == "Flat Discount" {
		discountAmount = promoDiscount
	}
	billingFinal := billingTotal - discountAmount

	log.Printf("[DEBUG] Promotion Applied: %s | Discount: %.2f | Final Billing: %.2f\n", promoType, discountAmount, billingFinal)

	// Insert a new billing record in the `billing_service` database
	billingInsertQuery := `
		INSERT INTO billing_service.Billing 
			(UserId, ReservationId, BillingTotal, PromoDiscount, BillingFinal, PaymentStatus, PaymentMethod, BillingDate)
		VALUES 
			(?, ?, ?, ?, ?, 'Pending', 'Nil', NOW())
	`
	_, err = db.Exec(billingInsertQuery, userId, reservationId, billingTotal, discountAmount, billingFinal)
	if err != nil {
		log.Println("[ERROR] Failed to create billing record in billing_service:", err)
		http.Error(w, "Failed to create billing record", http.StatusInternalServerError)
		return
	}

	log.Printf("[DEBUG] Billing record created in billing_service for Reservation ID: %d | Final Amount: %.2f\n", reservationId, billingFinal)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Reservation successfully created. Please proceed to billing.",
		"reservation_id": reservationId,
		"reserve_status": reserveStatus,
		"billing_total":  billingTotal,
		"promo_discount": discountAmount,
		"billing_final":  billingFinal,
	})
}

// UpdateReservation allows users to update their reservation details
func UpdateReservation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] UpdateReservation: Received request")

	// Extract userId from JWT middleware
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		log.Println("[ERROR] Unauthorized - Missing userId in request context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("[DEBUG] Updating reservation for userId: %d\n", userId)

	w.Header().Set("Content-Type", "application/json")

	// Define structure for update payload
	var updateData struct {
		ReservationId    int       `json:"reservation_id"`
		ReserveStartDate time.Time `json:"reserve_start_date"`
		ReserveEndDate   time.Time `json:"reserve_end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.Println("[ERROR] Invalid request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure start and end dates are valid
	if updateData.ReserveStartDate.After(updateData.ReserveEndDate) {
		log.Println("[ERROR] Invalid reservation dates")
		http.Error(w, "Reserve start date must be before end date", http.StatusBadRequest)
		return
	}

	// Fetch the existing reservation details (ReserveStatus and VehicleId)
	var existingReserveStatus string
	var currentVehicleId int
	err := db.QueryRow("SELECT ReserveStatus, VehicleId FROM Reservation WHERE ReservationId = ? AND UserId = ?",
		updateData.ReservationId, userId).Scan(&existingReserveStatus, &currentVehicleId)
	if err != nil {
		log.Println("[ERROR] Reservation not found or does not belong to user:", err)
		http.Error(w, "Reservation not found or access denied", http.StatusNotFound)
		return
	}
	log.Printf("[DEBUG] Existing reservation status: %s, VehicleId: %d\n", existingReserveStatus, currentVehicleId)

	// Fetch user's membership level to determine the price per hour
	var membershipLevel string
	err = db.QueryRow("SELECT MembershipLevel FROM user_management.UserService WHERE UserId = ?", userId).Scan(&membershipLevel)
	if err != nil {
		log.Println("[ERROR] Failed to fetch user membership level:", err)
		http.Error(w, "Failed to fetch user membership level", http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG] User membership level: %s\n", membershipLevel)

	// Assign hourly rate based on membership level
	var basePrice float64
	switch membershipLevel {
	case "VIP":
		basePrice = 30.00
	case "Premium":
		basePrice = 40.00
	default:
		basePrice = 50.00
	}
	log.Printf("[DEBUG] Base Price per Hour: %.2f\n", basePrice)

	// Calculate number of hours booked and new estimated cost
	hoursBooked := updateData.ReserveEndDate.Sub(updateData.ReserveStartDate).Hours()
	log.Printf("[DEBUG] Hours booked: %.2f\n", hoursBooked)

	estimatedTotalCost := basePrice * hoursBooked
	log.Printf("[DEBUG] Calculated Estimated Total Cost: %.2f\n", estimatedTotalCost)

	// Update the reservation in the database with new dates and recalculated cost
	query := `
		UPDATE Reservation 
		SET ReserveStartDate = ?, ReserveEndDate = ?, EstimatedTotalCost = ?, ModifiedDate = CURRENT_TIMESTAMP
		WHERE ReservationId = ? AND UserId = ?`
	_, err = db.Exec(query, updateData.ReserveStartDate, updateData.ReserveEndDate, estimatedTotalCost, updateData.ReservationId, userId)
	if err != nil {
		log.Println("[ERROR] Error updating reservation:", err)
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}
	log.Println("[DEBUG] Reservation updated successfully")

	// Update vehicle status to "Not Available" if necessary
	updateVehicleQuery := `UPDATE Vehicle SET VehicleStatus = 'Not Available' WHERE VehicleId = ?`
	_, err = db.Exec(updateVehicleQuery, currentVehicleId)
	if err != nil {
		log.Println("[ERROR] Failed to update vehicle status:", err)
		http.Error(w, "Failed to update vehicle status", http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG] Vehicle ID %d status updated to 'Not Available'\n", currentVehicleId)

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":              "Reservation updated successfully",
		"reservation_id":       updateData.ReservationId,
		"reserve_status":       "Active",
		"estimated_total_cost": estimatedTotalCost,
	})
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

	// Retrieve the vehicle ID associated with the reservation
	var vehicleId int
	err := db.QueryRow("SELECT VehicleId FROM Reservation WHERE ReservationId = ?", request.ReservationId).Scan(&vehicleId)
	if err != nil {
		log.Println("[ERROR] Failed to fetch vehicle ID:", err)
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}

	// Delete the reservation
	query := `DELETE FROM Reservation WHERE ReservationId = ?`
	_, err = db.Exec(query, request.ReservationId)
	if err != nil {
		log.Println("[ERROR] Error deleting reservation:", err)
		http.Error(w, "Failed to delete reservation", http.StatusInternalServerError)
		return
	}

	log.Println("[DEBUG] Reservation deleted successfully")

	// Update vehicle status to 'Available'
	updateVehicleQuery := `UPDATE Vehicle SET VehicleStatus = 'Available' WHERE VehicleId = ?`
	_, err = db.Exec(updateVehicleQuery, vehicleId)
	if err != nil {
		log.Println("[ERROR] Failed to update vehicle status:", err)
		http.Error(w, "Failed to update vehicle status", http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG] Vehicle ID %d status updated to 'Available'\n", vehicleId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation deleted successfully"})
}

// GetRentalHistory retrieves past successful reservations for the logged-in user
func GetRentalHistory(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] GetRentalHistory: Received request")

	// Extract userId from request context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		log.Println("[ERROR] Unauthorized - Missing userId in request context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("[DEBUG] Fetching rental history for userId: %d\n", userId)

	// Updated query to fetch rental history along with vehicle details
	query := `
		SELECT 
			r.ReservationId, 
			r.UserId, 
			r.VehicleId, 
			r.ReserveStatus, 
			r.ReserveStartDate, 
			r.ReserveEndDate, 
			r.EstimatedTotalCost, 
			r.CreatedDate, 
			r.ModifiedDate,
			v.VehicleMake, 
			v.VehicleModel, 
			v.VehicleType, 
			v.LicensePlate
		FROM 
			Reservation r
		INNER JOIN 
			Vehicle v ON r.VehicleId = v.VehicleId
		WHERE 
			r.UserId = ? AND r.ReserveStatus = 'Completed'
		ORDER BY 
			r.ReserveEndDate DESC`

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

		// Scan all reservation details including vehicle info
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
			&reservation.VehicleMake,  // Added VehicleMake
			&reservation.VehicleModel, // Added VehicleModel
			&reservation.VehicleType,  // Added VehicleType
			&reservation.LicensePlate, // Added LicensePlate
		); err != nil {
			log.Println("[ERROR] Error scanning rental history row:", err)
			http.Error(w, "Error retrieving rental history", http.StatusInternalServerError)
			return
		}

		// Convert date strings to time.Time format
		const timeLayout = "2006-01-02 15:04:05"

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

		log.Printf("[DEBUG] Retrieved rental history: %+v\n", reservation)
		reservations = append(reservations, reservation)
	}

	if len(reservations) == 0 {
		log.Println("[DEBUG] No rental history found")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "No rental history found"})
		return
	}

	log.Printf("[DEBUG] Total rental history records found: %d\n", len(reservations))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
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

	router.HandleFunc("/get-rental-history", middleware.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		GetRentalHistory(db, w, r)
	})).Methods("GET")
}
