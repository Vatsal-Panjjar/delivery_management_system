package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"delivery_management_system/internal/db"
	"delivery_management_system/internal/service"
)

// CreateOrder handles the creation of an order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int `json:"user_id"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create order in DB
	orderID, err := db.CreateOrder(request.UserID, "pending")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create order: %v", err), http.StatusInternalServerError)
		return
	}

	// Start tracking the order asynchronously
	go service.StartTrackingOrder(orderID)

	// Respond with order ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"order_id": orderID})
}
