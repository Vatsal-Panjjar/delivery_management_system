package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"delivery_management_system/internal/redis"
	"delivery_management_system/internal/db"
	"github.com/gorilla/mux"
)

// TrackOrder retrieves the status of an order
func TrackOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// Check if the status exists in Redis first
	status, err := redis.GetOrderStatus(orderID)
	if err == nil {
		// If found in Redis, return it
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"order_id": orderID,
			"status":   status,
		})
		return
	}

	// If not found in Redis, retrieve from the database
	status, err = db.GetOrderStatus(orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		log.Println("Error fetching order status:", err)
		return
	}

	// Store the status in Redis for future use
	redis.SetOrderStatus(orderID, status)

	// Respond with the status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"order_id": orderID,
		"status":   status,
	})
}
