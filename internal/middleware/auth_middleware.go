package middleware

import (
    "net/http"
    "strings"

    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        token = strings.TrimPrefix(token, "Bearer ")
        _, err := auth.ParseToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
