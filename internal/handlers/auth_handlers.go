package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-redis/redis/v8"
    "github.com/jmoiron/sqlx"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

func RegisterAuthRoutes(r *chi.Mux, db *sqlx.DB, rdb *redis.Client) {
    r.Route("/auth", func(r chi.Router) {
        r.Post("/signup", SignUpHandler(db))
        r.Post("/login", LoginHandler(db, rdb))
    })
}

func SignUpHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // TODO: parse JSON, hash password, insert into users table
        w.Write([]byte("Sign up endpoint"))
    }
}

func LoginHandler(db *sqlx.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // TODO: parse JSON, validate password, generate JWT, store session in Redis
        w.Write([]byte("Login endpoint"))
    }
}
