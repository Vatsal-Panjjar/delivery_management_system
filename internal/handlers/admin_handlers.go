package handlers

import (
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

type AdminHandler struct {
	store *db.Store
	cache *cache.Cache
	tracker *workers.OrderTracker
}

func RegisterAdminHandlers(r chi.Router, store *db.Store, cacheClient *cache.Cache, tracker *workers.OrderTracker) {
	h := &AdminHandler{store: store, cache: cacheClient, tracker: tracker}
	r.Group(func(gr chi.Router) {
		gr.Use(middleware.AuthMiddleware)
		gr.Get("/admin/orders", h.ListAllOrders)
		gr.Post("/admin/orders/{id}/status", h.UpdateOrderStatus)
	})
}

func (h *AdminHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.CtxRole).(string)
	if role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var orders []models.Order
	if err := h.store.DB.Select(&orders, "SELECT * FROM orders ORDER BY created_at DESC"); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func (h *AdminHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.CtxRole).(string)
	if role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	var req struct{ Status string `json:"status"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if _, err := h.store.DB.Exec("UPDATE orders SET status=$1, updated_at=now() WHERE id=$2", req.Status, id); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
