package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

type key string

const (
    CtxUserIDKey key = "userID"
    CtxRoleKey   key = "role"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }
        token = strings.TrimPrefix(token, "Bearer ")
        claims, err := auth.ParseToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), CtxUserIDKey, (*claims)["user_id"])
        ctx = context.WithValue(ctx, CtxRoleKey, (*claims)["role"])
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func CtxUserID(r *http.Request) string {
    if val := r.Context().Value(CtxUserIDKey); val != nil {
        return val.(string)
    }
    return ""
}

func CtxRole(r *http.Request) string {
    if val := r.Context().Value(CtxRoleKey); val != nil {
        return val.(string)
    }
    return ""
}
