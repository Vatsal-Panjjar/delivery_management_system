package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserCtxKey).(*auth.Claims)

	var req struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	order := models.Order{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(claims.UserID),
		Source:      req.Source,
		Destination: req.Destination,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.DB.Exec(`INSERT INTO orders (id, user_id, source, destination, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		order.ID, order.UserID, order.Source, order.Destination, order.Status, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (s *Server) CancelOrder(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserCtxKey).(*auth.Claims)
	orderID := chi.URLParam(r, "id")

	_, err := s.DB.Exec(`UPDATE orders SET status=$1, updated_at=$2
		WHERE id=$3 AND user_id=$4 AND status != 'cancelled'`,
		"cancelled", time.Now(), orderID, claims.UserID)
	if err != nil {
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
}

func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	var order models.Order
	err := s.DB.Get(&order, "SELECT * FROM orders WHERE id=$1", orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func (s *Server) ListOrders(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserCtxKey).(*auth.Claims)
	var orders []models.Order
	err := s.DB.Select(&orders, "SELECT * FROM orders WHERE user_id=$1 ORDER BY created_at DESC", claims.UserID)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}
