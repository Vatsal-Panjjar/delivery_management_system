package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/workers"
)

type OrderHandler struct {
	store   *db.Store
	cache   *cache.Cache
	tracker *workers.OrderTracker
}

func RegisterOrderHandlers(r chi.Router, store *db.Store, cacheClient *cache.Cache, tracker *workers.OrderTracker) {
	h := &OrderHandler{store: store, cache: cacheClient, tracker: tracker}
	r.Group(func(gr chi.Router) {
		gr.Use(middleware.AuthMiddleware)
		gr.Post("/orders", h.CreateOrder)
		gr.Get("/orders", h.ListOrders)
		gr.Post("/orders/{id}/cancel", h.CancelOrder)
	})
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct{ Description string `json:"description"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	uid := r.Context().Value(middleware.CtxUserID).(int)

	var id int
	if err := h.store.DB.QueryRow("INSERT INTO orders (user_id, description, status, created_at, updated_at) VALUES ($1,$2,'pending',now(),now()) RETURNING id",
		uid, req.Description).Scan(&id); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	h.tracker.Enqueue(id)
	_ = h.cache.Set(context.Background(), "order:"+strconv.Itoa(id)+":status", "pending", 0)

	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "pending"})
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.CtxUserID).(int)
	var orders []models.Order
	if err := h.store.DB.Select(&orders, "SELECT * FROM orders WHERE user_id=$1 ORDER BY created_at DESC", uid); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.CtxUserID).(int)
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	tx := h.store.DB.MustBegin()
	var cur string
	if err := tx.Get(&cur, "SELECT status FROM orders WHERE id=$1 FOR UPDATE", id); err != nil {
		tx.Rollback()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if cur == "delivered" {
		tx.Rollback()
		http.Error(w, "cannot cancel delivered order", http.StatusBadRequest)
		return
	}
	if _, err := tx.Exec("UPDATE orders SET status='cancelled', updated_at=now() WHERE id=$1 AND user_id=$2", id, uid); err != nil {
		tx.Rollback()
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	_ = h.cache.Del(context.Background(), "order:"+idStr+":status")
	w.WriteHeader(http.StatusOK)
}
