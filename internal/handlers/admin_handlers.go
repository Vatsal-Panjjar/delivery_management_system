package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/jmoiron/sqlx"
)

// Admin middleware ensures only admin users can access these routes
func AdminMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims := r.Context().Value("user").(map[string]interface{})
        role := claims["role"].(string)
        if role != "admin" {
            http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// RegisterAdminRoutes sets up all admin routes
func RegisterAdminRoutes(r *chi.Mux, db *sqlx.DB) {
    r.Route("/admin", func(r chi.Router) {
        r.Use(AdminMiddleware)

        // Get all orders
        r.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
            var orders []struct {
                ID          int    `db:"id" json:"id"`
                Username    string `db:"username" json:"username"`
                Description string `db:"description" json:"description"`
                Status      string `db:"status" json:"status"`
            }
            err := db.Select(&orders, "SELECT * FROM orders")
            if err != nil {
                http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
                return
            }
            json.NewEncoder(w).Encode(orders)
        })

        // Update order status
        r.Post("/orders/{id}/update", func(w http.ResponseWriter, r *http.Request) {
            orderID := chi.URLParam(r, "id")
            var body struct {
                Status string `json:"status"`
            }
            if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
                http.Error(w, "Invalid request", http.StatusBadRequest)
                return
            }

            id, err := strconv.Atoi(orderID)
            if err != nil {
                http.Error(w, "Invalid order ID", http.StatusBadRequest)
                return
            }

            _, err = db.Exec("UPDATE orders SET status=$1 WHERE id=$2", body.Status, id)
            if err != nil {
                http.Error(w, "Failed to update order", http.StatusInternalServerError)
                return
            }

            w.WriteHeader(http.StatusOK)
        })

        // Get all users
        r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
            var users []struct {
                ID       int    `db:"id" json:"id"`
                Username string `db:"username" json:"username"`
                Role     string `db:"role" json:"role"`
            }
            err := db.Select(&users, "SELECT id, username, role FROM users")
            if err != nil {
                http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
                return
            }
            json.NewEncoder(w).Encode(users)
        })

        // Change user role
        r.Post("/users/{id}/role", func(w http.ResponseWriter, r *http.Request) {
            userID := chi.URLParam(r, "id")
            var body struct {
                Role string `json:"role"`
            }
            if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
                http.Error(w, "Invalid request", http.StatusBadRequest)
                return
            }

            id, err := strconv.Atoi(userID)
            if err != nil {
                http.Error(w, "Invalid user ID", http.StatusBadRequest)
                return
            }

            _, err = db.Exec("UPDATE users SET role=$1 WHERE id=$2", body.Role, id)
            if err != nil {
                http.Error(w, "Failed to update user role", http.StatusInternalServerError)
                return
            }

            w.WriteHeader(http.StatusOK)
        })
    })
}
