package billing

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"user-management/middleware"

	"github.com/gorilla/mux"
)

// Billing struct to represent the billing model
type Billing struct {
	BillingId     int       `json:"billing_id"`
	UserId        int       `json:"user_id"`
	ReservationId int       `json:"reservation_id"`
	BillingTotal  float64   `json:"billing_total"`
	PromoDiscount float64   `json:"promo_discount"`
	BillingFinal  float64   `json:"billing_final"`
	PaymentStatus string    `json:"payment_status"`
	PaymentMethod string    `json:"payment_method"`
	BillingDate   time.Time `json:"billing_date"`
}

// GetPaidBillings retrieves all paid billing records for the logged-in user
func GetPaidBillings(db *sql.DB, userId int) ([]Billing, error) {
	log.Printf("[DEBUG] Fetching paid billings for userId: %d\n", userId)

	query := `SELECT BillingId, UserId, ReservationId, BillingTotal, PromoDiscount, BillingFinal, PaymentStatus, PaymentMethod, BillingDate 
				FROM Billing WHERE UserId = ? AND PaymentStatus = 'Paid' ORDER BY BillingDate DESC`
	rows, err := db.Query(query, userId)
	if err != nil {
		log.Println("[ERROR] Database query error:", err)
		return nil, err
	}
	defer rows.Close()

	var billings []Billing
	for rows.Next() {
		var billing Billing
		var billingDate sql.NullString

		if err := rows.Scan(
			&billing.BillingId,
			&billing.UserId,
			&billing.ReservationId,
			&billing.BillingTotal,
			&billing.PromoDiscount,
			&billing.BillingFinal,
			&billing.PaymentStatus,
			&billing.PaymentMethod,
			&billingDate,
		); err != nil {
			log.Println("[ERROR] Error scanning billing row:", err)
			return nil, err
		}

		if billingDate.Valid {
			billing.BillingDate, _ = time.Parse("2006-01-02 15:04:05", billingDate.String)
		}

		billings = append(billings, billing)
	}

	return billings, nil
}

// GetUnpaidBillings retrieves all unpaid billing records for the logged-in user
func GetUnpaidBillings(db *sql.DB, userId int) ([]Billing, error) {
	log.Printf("[DEBUG] Fetching unpaid billings for userId: %d\n", userId)

	query := `SELECT BillingId, UserId, ReservationId, BillingTotal, PromoDiscount, BillingFinal, PaymentStatus, PaymentMethod, BillingDate 
				FROM Billing WHERE UserId = ? AND PaymentStatus = 'Pending' ORDER BY BillingDate DESC`
	rows, err := db.Query(query, userId)
	if err != nil {
		log.Println("[ERROR] Database query error:", err)
		return nil, err
	}
	defer rows.Close()

	var billings []Billing
	for rows.Next() {
		var billing Billing
		var billingDate sql.NullString

		if err := rows.Scan(
			&billing.BillingId,
			&billing.UserId,
			&billing.ReservationId,
			&billing.BillingTotal,
			&billing.PromoDiscount,
			&billing.BillingFinal,
			&billing.PaymentStatus,
			&billing.PaymentMethod,
			&billingDate,
		); err != nil {
			log.Println("[ERROR] Error scanning billing row:", err)
			return nil, err
		}

		if billingDate.Valid {
			billing.BillingDate, _ = time.Parse("2006-01-02 15:04:05", billingDate.String)
		}

		billings = append(billings, billing)
	}

	return billings, nil
}

// PayUnpaidBilling allows the user to pay for an unpaid billing record and activates the reservation
func PayUnpaidBilling(db *sql.DB, userId int, billingId string, paymentMethod string) (Billing, error) {
    log.Printf("[DEBUG] Processing payment for userId: %d, BillingId: %s\n", userId, billingId)

    query := `SELECT BillingId, UserId, ReservationId, PaymentStatus, BillingFinal FROM Billing WHERE BillingId = ? AND UserId = ?`
    row := db.QueryRow(query, billingId, userId)

    var billing Billing
    if err := row.Scan(&billing.BillingId, &billing.UserId, &billing.ReservationId, &billing.PaymentStatus, &billing.BillingFinal); err != nil {
        log.Println("[ERROR] Error retrieving billing record:", err)
        return billing, err
    }

    if billing.PaymentStatus == "Paid" {
        log.Println("[ERROR] Payment already made for this billing")
        return billing, nil
    }

    log.Println("[DEBUG] Payment processing successful")

    // ✅ Update Billing Status to "Paid" and Payment Method
    updateBillingQuery := `UPDATE Billing SET PaymentStatus = 'Paid', PaymentMethod = ? WHERE BillingId = ? AND UserId = ?`
    _, err := db.Exec(updateBillingQuery, paymentMethod, billing.BillingId, userId)
    if err != nil {
        log.Println("[ERROR] Failed to update payment status:", err)
        return billing, err
    }

    // ✅ Update Reservation Status to "Active"
    updateReservationQuery := `UPDATE Reservations SET ReserveStatus = 'Active' WHERE ReservationId = ? AND UserId = ?`
    _, err = db.Exec(updateReservationQuery, billing.ReservationId, userId)
    if err != nil {
        log.Println("[ERROR] Failed to update reservation status:", err)
        return billing, err
    }

    // ✅ Return updated details
    billing.PaymentStatus = "Paid"
    billing.PaymentMethod = paymentMethod

    return billing, nil
}


// GetBillingDetails retrieves billing details for a specific reservation
func GetBillingDetails(db *sql.DB, userId int, reservationId string) (Billing, error) {
	log.Printf("[DEBUG] Fetching billing details for userId: %d, ReservationId: %s\n", userId, reservationId)

	// ✅ Retrieve MySQL Timestamp and Convert it to Singapore Time
	query := `SELECT BillingId, UserId, ReservationId, BillingTotal, PromoDiscount, BillingFinal, PaymentStatus, PaymentMethod, 
						BillingDate FROM Billing WHERE UserId = ? AND ReservationId = ?`
	row := db.QueryRow(query, userId, reservationId)

	var billing Billing
	var rawBillingDate sql.NullTime // Store raw timestamp

	// ✅ Debug: Log Raw Database Query Result
	if err := row.Scan(
		&billing.BillingId,
		&billing.UserId,
		&billing.ReservationId,
		&billing.BillingTotal,
		&billing.PromoDiscount,
		&billing.BillingFinal,
		&billing.PaymentStatus,
		&billing.PaymentMethod,
		&rawBillingDate,
	); err != nil {
		log.Println("[ERROR] Billing record not found:", err)
		return billing, err
	}

	// ✅ Debugging: Check Raw Timestamp from MySQL
	if rawBillingDate.Valid {
		log.Printf("[DEBUG] Raw BillingDate (MySQL UTC): %s\n", rawBillingDate.Time.Format("2006-01-02 15:04:05 MST"))
	} else {
		log.Println("[WARN] Billing date is NULL, using current system time.")
		rawBillingDate.Time = time.Now().UTC() // Default to UTC
	}

	// ✅ Convert to Singapore Time (UTC+8)
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		log.Println("[ERROR] Failed to load timezone, using UTC as fallback:", err)
		billing.BillingDate = rawBillingDate.Time.UTC() // Use UTC if conversion fails
	} else {
		billing.BillingDate = rawBillingDate.Time.In(loc) // Convert to SGT
	}

	// ✅ Debugging: Log Converted Timestamp
	log.Printf("[DEBUG] Converted BillingDate (SGT): %s\n", billing.BillingDate.Format("2006-01-02 15:04:05 MST"))

	// ✅ Return Struct to API
	return billing, nil
}

// Handlers
func GetPaidBillingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value("userId").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		billings, err := GetPaidBillings(db, userId)
		if err != nil {
			http.Error(w, "Error fetching paid billings", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(billings)
	}
}

func GetUnpaidBillingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value("userId").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		billings, err := GetUnpaidBillings(db, userId)
		if err != nil {
			http.Error(w, "Error fetching unpaid billings", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(billings)
	}
}

func PayUnpaidBillingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value("userId").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		params := mux.Vars(r)
		billingId := params["billing_id"]

		var request struct {
			PaymentMethod string `json:"payment_method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		billing, err := PayUnpaidBilling(db, userId, billingId, request.PaymentMethod)
		if err != nil {
			http.Error(w, "Payment processing failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":        "Payment successful",
			"billing_id":     billing.BillingId,
			"payment_status": billing.PaymentStatus,
			"payment_method": billing.PaymentMethod,
		})
	}
}

// Handler function to fetch billing details for a reservation
func GetBillingDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value("userId").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
		reservationId := params["reservation_id"]

		// ✅ Fetch Billing Details
		billing, err := GetBillingDetails(db, userId, reservationId)
		if err != nil {
			http.Error(w, "Billing not found", http.StatusNotFound)
			return
		}

		// ✅ Debugging: Log API Response Before Sending
		log.Printf("[DEBUG] Sending API Response: BillingDate (SGT) → %s\n",
			billing.BillingDate.Format("2006-01-02 15:04:05 MST"))

		// ✅ Send API Response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(billing)
	}
}

// RegisterBillingRoutes sets up the billing routes
func RegisterBillingRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/get-paid-billings", middleware.ValidateJWT(GetPaidBillingsHandler(db))).Methods("GET")
	router.HandleFunc("/get-unpaid-billings", middleware.ValidateJWT(GetUnpaidBillingsHandler(db))).Methods("GET")
	router.HandleFunc("/pay-unpaid-billing/{billing_id}", middleware.ValidateJWT(PayUnpaidBillingHandler(db))).Methods("POST")
	router.HandleFunc("/get-billing/{reservation_id}", middleware.ValidateJWT(GetBillingDetailsHandler(db))).Methods("GET")
}
