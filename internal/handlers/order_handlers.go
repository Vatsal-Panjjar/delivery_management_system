package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/go-redis/redis/v8"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

type OrderHandler struct {
	DB  *sqlx.DB
	RDB *redis.Client
}

func RegisterOrderRoutes(r chi.Router, db *sqlx.DB, rdb *redis.Client) {
	h := &OrderHandler{DB: db, RDB: rdb}
	r.Group(func(gr chi.Router) {
		gr.Use(middleware.AuthMiddleware)
		gr.Post("/orders", h.CreateOrder)
		gr.Get("/orders", h.GetOrders)
	})
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Description string `json:"description"`
	}
	var req Req
	json.NewDecoder(r.Body).Decode(&req)

	_, err := h.DB.Exec(`INSERT INTO orders (description, status) VALUES ($1, 'pending')`, req.Description)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order created successfully"))
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, description, status FROM orders`)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Order struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	var orders []Order
	for rows.Next() {
		var o Order
		rows.Scan(&o.ID, &o.Description, &o.Status)
		orders = append(orders, o)
	}
	json.NewEncoder(w).Encode(orders)
}
