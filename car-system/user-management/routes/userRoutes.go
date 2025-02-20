package routes

import (
	"database/sql"
	"net/http"
	"user-management/controller"
	"user-management/middleware"

	"github.com/gorilla/mux"
)

// initialize user routes
func UserRoutes(router *mux.Router, db *sql.DB) {
	// user signup
	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		controller.Signup(db, w, r)
	}).Methods("POST")

	// user login (generate JWT token)
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.Login(db, w, r)
	}).Methods("POST")

	// Routes that require JWT authentication
	// view user account details
	router.HandleFunc("/account", middleware.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		controller.ViewAccountDetails(db, w, r)
	})).Methods("GET")

	// update user account details
	router.HandleFunc("/account-update", middleware.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		controller.UpdateAccountDetails(db, w, r)
	})).Methods("PUT")

	// update user membership level
	router.HandleFunc("/update-membership", middleware.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		controller.UpdateMembershipLevel(db, w, r)
	})).Methods("PUT")

}
