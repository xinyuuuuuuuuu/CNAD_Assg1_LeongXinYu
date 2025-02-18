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
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println(".env file not found, using system environment variables")
	} else {
		jwtSecret = []byte(os.Getenv("JWT_SECRET")) // store secret key for reuse
		envLoaded = true
	}
}

func GetJWTSecret() []byte {
	if !envLoaded {
		LoadEnv() // ensure env is loaded before accessing
	}
	return jwtSecret
}
