package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

// AdminPage shows all orders for admin dashboard.
func (s *Server) AdminPage(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserCtxKey).(*auth.Claims)
	if claims.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var orders []models.Order
	err := s.DB.Select(&orders, "SELECT * FROM orders ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

// UpdateOrderStatus allows admin to update any order's status.
func (s *Server) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserCtxKey).(*auth.Claims)
	if claims.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := s.DB.Exec(`UPDATE orders SET status=$1, updated_at=$2 WHERE id=$3`,
		req.Status, time.Now(), req.OrderID)
	if err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": req.Status})
}
