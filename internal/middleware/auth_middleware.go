package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

type ctxKey string

const CtxUser = ctxKey("user")
const CtxRole = ctxKey("role")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        header := r.Header.Get("Authorization")
        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        token := strings.TrimPrefix(header, "Bearer ")
        username, role := auth.ParseToken(token) // Returns username, role
        ctx := context.WithValue(r.Context(), CtxUser, username)
        ctx = context.WithValue(ctx, CtxRole, role)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
