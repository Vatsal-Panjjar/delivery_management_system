package handlers

import (
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
	UserID      int    `db:"user_id" json:"user_id"`
	Description string `db:"description" json:"description"`
	Status      string `db:"status" json:"status"`
	Username    string `db:"username" json:"username,omitempty"` // Optional for admin view
}

// RegisterOrderRoutes sets up the routes for orders
func RegisterOrderRoutes(r chi.Router, db *sqlx.DB, rdb *redis.Client) {
	r.Route("/orders", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(int)
			orders := []Order{}
			err := db.Select(&orders, "SELECT id, user_id, description, status FROM orders WHERE user_id=$1", userID)
			if err != nil {
				http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(orders)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(int)
			var order Order
			if err := json.NewDecoder(r.Body).Decode(&order); er
