module vehicle-service

go 1.23.2

require (
	github.com/go-sql-driver/mysql v1.9.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/rs/cors v1.11.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
)

require user-management v0.0.0

replace user-management => ../user-management
