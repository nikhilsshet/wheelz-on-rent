package middleware

import (
	"context"
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
		// Step 1: Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Step 2: Extract token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Step 3: Validate token
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Step 4: Extract user details from claims
		userIDFloat, okID := claims["id"].(float64) // numbers come as float64 in JWT
		userRole, okRole := claims["role"].(string)

		if !okID || !okRole {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID := int(userIDFloat)

		// Step 5: Add to request context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, userRole)

		// Step 6: Continue
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

