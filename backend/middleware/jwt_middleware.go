package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/nikhilsshet/wheelz-on-rent/backend/utils"
)

type contextKey string

const (
	ClaimsContextKey contextKey = "jwtClaims"
	UserIDKey        contextKey = "userID"
	UserRoleKey      contextKey = "userRole"
)

var jwtSecret = []byte("your_secret_key")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// You can attach claims to context here (optional)
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AuthMiddleware invoked")

		// Step 1: Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Step 2: Validate the JWT
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Step 3: Extract user info from claims
		userID, _ := claims["id"].(float64)   // JWT lib returns numbers as float64
		userRole, _ := claims["role"].(string)

		// Step 4: Store in context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		ctx = context.WithValue(ctx, UserIDKey, int(userID))
		ctx = context.WithValue(ctx, UserRoleKey, userRole)

		// Step 5: Proceed
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

