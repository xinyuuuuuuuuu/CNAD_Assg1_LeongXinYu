package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"user-management/config"

	"github.com/golang-jwt/jwt"
)

// middleware to validate jwt
func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Println("🟢 Received Authorization Header:", authHeader)
		if authHeader == "" {
			fmt.Println("🔴 Missing Token")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// remove bearer prefix from the token string
		// usual jwt format: Authorization: Bearer <JWT_TOKEN>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString) // remove trailing space

		// parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}

			// return the jwt secret key for validation
			return config.GetJWTSecret(), nil
		})

		// check if token is valid
		if err != nil || !token.Valid {
			fmt.Println("🔴 Invalid Token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// extract userId from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			fmt.Println("🔴 Invalid Token Claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// convert userId form float64 to int
		userIdFloat, ok := claims["userId"].(float64)
		if !ok {
			fmt.Println("🔴 Invalid User ID in Token")
			http.Error(w, "Invalid user id in token", http.StatusUnauthorized)
			return
		}
		userId := int(userIdFloat)

		fmt.Println("🟢 Token Valid. Extracted userID:", userId)

		// store userid in request context
		ctx := context.WithValue(r.Context(), "userId", userId)
		next(w, r.WithContext(ctx))
	}
}
