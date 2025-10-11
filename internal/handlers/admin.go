package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-redis/redis/v8"
    "github.com/jmoiron/sqlx"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
    "context"
)

func RegisterAdminRoutes(r chi.Router, db *sqlx.DB, rdb *redis.Client) {
    r.Get("/admin/orders", auth.RequireAdmin(adminListOrders(db, rdb)))
}

func adminListOrders(db *sqlx.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Admin authentication already handled by RequireAdmin
        keys, _ := rdb.Keys(context.Background(), "*").Result()
        orders := make([]map[string]string, 0)
        for _, k := range keys {
            status, _ := rdb.Get(context.Background(), k).Result()
            orders = append(orders, map[string]string{"id": k, "status": status})
        }
        json.NewEncoder(w).Encode(orders)
    }
}
