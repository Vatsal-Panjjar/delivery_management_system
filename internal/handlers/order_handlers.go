package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

type OrderHandler struct {
	DB  *sqlx.DB
	RDB *redis.Client
}

func RegisterUserRoutes(r *chi.Mux, db *sqlx.DB, rdb *redis.Client) {
	h := &OrderHandler{DB: db, RDB: rdb}

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/orders", h.CreateOrder)
		r.Get("/orders/{id}", h.GetOrder)
		r.Get("/orders", h.ListOrders)
	})
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Product string `json:"product"`
		Qty     int    `json:"qty"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	// Insert into DB (simplified)
	_, err := h.DB.Exec("INSERT INTO orders (product, quantity, status) VALUES ($1,$2,'pending')", req.Product, req.Qty)
	if err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var order struct {
		ID      int
		Product string
		Qty     int
		Status  string
	}
	h.DB.QueryRow("SELECT id, product, quantity, status FROM orders WHERE id=$1", id).Scan(&order.ID, &order.Product, &order.Qty, &order.Status)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	rows, _ := h.DB.Query("SELECT id, product, quantity, status FROM orders")
	defer rows.Close()
	orders := []map[string]any{}
	for rows.Next() {
		var id, qty int
		var product, status string
		rows.Scan(&id, &product, &qty, &status)
		orders = append(orders, map[string]any{"id": id, "product": product, "qty": qty, "status": status})
	}
	json.NewEncoder(w).Encode(orders)
}
