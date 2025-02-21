package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once      sync.Once
	jwtSecret []byte
	envLoaded bool
)

// load the secret key
func LoadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println(".env file not found, using system environment variables")
	} else {
		jwtSecret = []byte(os.Getenv("JWT_SECRET")) // store secret key for reuse
		envLoaded = true
	}
}

// GetJWTSecret returns the JWT secret key
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "bWltMnpqNGw1ejV6cG45NWJrNGt6OTdybnpvMGhqZnQ0MnRyM2w4NGF1Z3pjMzZtdG1zcnY1OHpjczNkYXpldnpqZjk3anUwMDZoZGR3Z210a3ZxczI4dmZ1MmtsZ25uZTVxZDBjMGZyeTRkaDF2b2x4aHd1d2w0eWk2ZWxxMmg="
		// Default if env variable is not set
	}
	return []byte(secret)
}
