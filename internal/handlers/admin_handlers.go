package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"delivery_management_system/internal/db"
	"delivery_management_system/internal/models"
	"delivery_management_system/internal/redis"
	"delivery_management_system/internal/auth"
	"github.com/gorilla/mux"
)

// AdminHandler is a simple protected handler for admin routes
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow admins
	userID, err := auth.AuthenticateJWT(r)
	if err != nil || !auth.IsAdmin(userID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// You can extend this to show admin functionality (list orders, etc.)
	w.Write([]byte("Admin Dashboard"))
}

// UpdateOrderStatus updates the status of an order (for admin)
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	var statusUpdate struct {
		Status string `json:"status"`
	}

	// Decode the status update from the body
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the order status in the database
	err := db.UpdateOrderStatus(orderID, statusUpdate.Status)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		log.Println("Error updating order status:", err)
		return
	}

	// Update the order status in Redis
	redis.SetOrderStatus(orderID, statusUpdate.Status)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"order_id": orderID,
		"status":   statusUpdate.Status,
	})
}
