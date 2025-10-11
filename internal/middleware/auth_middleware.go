package middleware

import (
    "net/http"

    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, err := auth.VerifyToken(r)
        if err != nil {
            http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
