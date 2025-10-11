package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

// AuthMiddleware ensures requests have a valid JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ParseToken now returns only 2 values: userID, role
		userID, role := auth.ParseToken(tokenString)

		// If userID is empty, consider token invalid
		if userID == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "role", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
