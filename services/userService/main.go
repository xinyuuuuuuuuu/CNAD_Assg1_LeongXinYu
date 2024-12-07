package main

import (
	"database/sql"
	"cnad_assg1_leongxinyu/services/userService/controller"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/cnad_assg1")
	
	// handle error
	if err != nil {
		panic(err.Error())
	} 
	
	// database operation
	controller.Signup(db)
	controller.Login(db)
	defer db.Close() 
}