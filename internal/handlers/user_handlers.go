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

// CreateOrder creates a new order for a user
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	var user models.User

	// Parse incoming request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create the order in the database
	order.UserID = user.ID
	order.Status = "Created"
	orderID, err := db.CreateOrder(&order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		log.Println("Error creating order:", err)
		return
	}

	// Store order in Redis for fast access
	redis.SetOrderStatus(orderID, order.Status)

	// Respond with the new order ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": orderID,
		"status":   order.Status,
	})
}
