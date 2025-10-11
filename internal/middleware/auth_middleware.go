package middleware

import (
    "net/http"
    "strings"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        userID, role, err := auth.ParseToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add user info to context if needed
        ctx := r.Context()
        ctx = context.WithValue(ctx, "userID", userID)
        ctx = context.WithValue(ctx, "role", role)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
