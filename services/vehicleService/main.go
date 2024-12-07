package vehicleService

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/cnad_assg1")
	
	// handle error
	if err != nil {
		panic(err.Error())
	} 
	
	// database operation
	
	defer db.Close() 
}