package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

// Order represents the order structure
type Order struct {
	ID          int    `db:"id" json:"id"`
	UserID      string `db:"user_id" json:"user_id"`
	Description string `db:"description" json:"description"`
	Status      string `db:"status" json:"status"`
	Username    string `db:"username" json:"username,omitempty"` // Optional for admin view
}

// RegisterOrderRoutes sets up the routes for orders
func RegisterOrderRoutes(r chi.Router, db *sqlx.DB, rdb *redis.Client) {
	r.Route("/orders", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Get all orders for the logged-in user
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(string)
			orders := []Order{}
			err := db.Select(&orders, "SELECT id, user_id, description, status FROM orders WHERE user_id=$1", userID)
			if err != nil {
				http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(orders)
		})

		// Create a new order
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(string)
			var order Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			// Insert order into Postgres
			var orderID int
			err := db.QueryRow(
				"INSERT INTO orders (user_id, description, status) VALUES ($1, $2, 'pending') RETURNING id",
				userID, order.Description).Scan(&orderID)
			if err != nil {
				http.Error(w, "Failed to create order", http.StatusInternalServerError)
				return
			}

			order.ID = orderID
			order.UserID = userID
			order.Status = "pending"

			// Optional: cache order in Redis
			rdb.Set(context.Background(), "order:"+strconv.Itoa(orderID), order.Description, 0)

			json.NewEncoder(w).Encode(order)
		})

		// Cancel an existing order
		r.Post("/{orderID}/cancel", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(string)
			orderIDStr := chi.URLParam(r, "orderID")
			orderID, _ := strconv.Atoi(orderIDStr)

			// Update status to cancelled
			_, err := db.Exec("UPDATE orders SET status='cancelled' WHERE id=$1 AND user_id=$2", orderID, userID)
			if err != nil {
				http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
				return
			}

			// Remove from Redis cache if exists
			rdb.Del(context.Background(), "order:"+orderIDStr)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Order cancelled"})
		})
	})
}
