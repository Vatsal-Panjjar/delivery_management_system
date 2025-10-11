package handlers

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-redis/redis/v8"
    "github.com/jmoiron/sqlx"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

func RegisterOrderRoutes(r *chi.Mux, db *sqlx.DB, rdb *redis.Client) {
    r.Route("/orders", func(r chi.Router) {
        r.Post("/", middleware.AuthMiddleware(CreateOrderHandler(db, rdb)))
        r.Get("/{id}", middleware.AuthMiddleware(GetOrderHandler(db, rdb)))
        r.Put("/{id}/cancel", middleware.AuthMiddleware(CancelOrderHandler(db, rdb)))
    })
}

func CreateOrderHandler(db *sqlx.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Create Order"))
    }
}

func GetOrderHandler(db *sqlx.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Get Order"))
    }
}

func CancelOrderHandler(db *sqlx.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Cancel Order"))
    }
}
