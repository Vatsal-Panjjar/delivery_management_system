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

type Order struct {
    ID     string `json:"id"`
    Status string `json:"status"`
    Item   string `json:"item"`
}

// Register order routes
func RegisterOrderRoutes(r chi.Router, db *sqlx.DB, rdb *redis.Client) {
    r.Post("/orders", createOrder(db, rdb))
    r.Get("/orders/{id}", trackOrder(db, rdb))
    r.Patch("/orders/{id}/cancel", cancelOrder(db, rdb))
}

func createOrder(db *sqlx.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, err := auth.VerifyToken(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var o Order
        if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        o.Status = "created"

        // Save to Redis for demo
        rdb.Set(context.Background(), o.ID, o.Status, 0)

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(o)
    }
}

func trackOrder(db *sqlx.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, err := auth.VerifyToken(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        id := chi.URLParam(r, "id")
        status, err := rdb.Get(context.Background(), id).Result()
        if err != nil {
            http.Error(w, "Order not found", http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(Order{ID: id, Status: status})
    }
}

func cancelOrder(db *sqlx.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, err := auth.VerifyToken(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        id := chi.URLParam(r, "id")
        rdb.Set(context.Background(), id, "cancelled", 0)
        json.NewEncoder(w).Encode(Order{ID: id, Status: "cancelled"})
    }
}
