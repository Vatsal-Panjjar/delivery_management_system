package handlers

import (
    "net/http"
    "encoding/json"
    "github.com/go-chi/chi/v5"
    "github.com/go-redis/redis/v8"
    "context"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

type Order struct {
    ID     int    `db:"id"`
    UserID string `db:"user_id"`
    Item   string `db:"item"`
    Status string `db:"status"`
}

func RegisterOrderRoutes(r *chi.Mux, rdb *redis.Client) {
    r.Post("/order", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var order Order
        json.NewDecoder(r.Body).Decode(&order)
        order.UserID = middleware.CtxUserID(r)
        order.Status = "created"
        _, err := db.DB.Exec("INSERT INTO orders (user_id, item, status) VALUES ($1, $2, $3)", order.UserID, order.Item, order.Status)
        if err != nil {
            http.Error(w, "Failed to create order", http.StatusInternalServerError)
            return
        }
        w.Write([]byte("Order created"))
    })))
}
