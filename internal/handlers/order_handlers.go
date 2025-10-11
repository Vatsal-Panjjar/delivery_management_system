package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-redis/redis/v8"
    "github.com/jmoiron/sqlx"
)

type Order struct {
    ID          int    `db:"id" json:"id"`
    Username    string `db:"username" json:"username"`
    Description string `db:"description" json:"description"`
    Status      string `db:"status" json:"status"`
}

func RegisterOrderRoutes(r *chi.Mux, db *sqlx.DB, rdb *redis.Client) {
    r.Route("/orders", func(r chi.Router) {
        r.Use(AuthMiddleware)

        r.Get("/", func(w http.ResponseWriter, r *http.Request) {
            claims := r.Context().Value("user").(map[string]interface{})
            username := claims["username"].(string)
            role := claims["role"].(string)

            var orders []Order
            if role == "admin" {
                db.Select(&orders, "SELECT * FROM orders")
            } else {
                db.Select(&orders, "SELECT * FROM orders WHERE username=$1", username)
            }
            json.NewEncoder(w).Encode(orders)
        })

        r.Post("/", func(w http.ResponseWriter, r *http.Request) {
            var input struct{ Description string `json:"description"` }
            json.NewDecoder(r.Body).Decode(&input)
            claims := r.Context().Value("user").(map[string]interface{})
            username := claims["username"].(string)

            _, err := db.Exec("INSERT INTO orders (username, description, status) VALUES ($1,$2,'pending')",
                username, input.Description)
            if err != nil {
                http.Error(w, "Failed to create order", http.StatusInternalServerError)
                return
            }
            w.WriteHeader(http.StatusCreated)
        })

        r.Post("/{id}/cancel", func(w http.ResponseWriter, r *http.Request) {
            orderID := chi.URLParam(r, "id")
            _, err := db.Exec("UPDATE orders SET status='cancelled' WHERE id=$1", orderID)
            if err != nil {
                http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
                return
            }
            w.WriteHeader(http.StatusOK)
        })
    })
}
