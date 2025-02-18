package routes

import (
	"database/sql"
	"net/http"
	"user-management/controller"

	"github.com/gorilla/mux"
)

// initialize user routes
func UserRoutes(router *mux.Router, db *sql.DB) {
	// user signup
	router.HandlerFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		controller.Signup(db, w, r)
	}).Methods("POST")

	// user login (generate JWT token)
	router.HandlerFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.Login(db, w, r)
	}).Methods("POST")

	// Routes that require JWT authentication
	// view user account details
	router.HandlerFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		controller.ViewAccountDetails(db, w, r)
	}).Methods("GET")

	// update user account details
	router.HandlerFunc("/account-update", func(w http.ResponseWriter, r *http.Request) {
		controller.UpdateAccountDetails(db, w, r)
	}).Methods("PUT")

}
