package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"user-management/config"
	"context"

	"github.com/golang-jwt/jwt"
)

// middleware to validate jwt
func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
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
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// extract userId from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// convert userId form float64 to int
		userIdFloat, ok := int(claims["userId"].float(64))
		if !ok {
			http.Error(w, "Invalid user id in token", http.StatusUnauthorized)
			return
		}
		userId := int(userIdFloat)

		// store userid in request context
		ctx := context.WithValue(r.Context(), "userId", userId)
		next(w, r.WithContext(ctx))

		// proceed to next handler if token is valid
		next(w, r)
	}
}
