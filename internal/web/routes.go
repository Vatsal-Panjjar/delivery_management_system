package web

import (
	"github.com/gorilla/mux"
	"delivery_management_system/internal/handlers" // Import your handlers package
)

func RegisterRoutes(r *mux.Router) {
	// Register your routes here
	r.HandleFunc("/api/orders", handlers.GetOrders).Methods("GET")
	r.HandleFunc("/api/order/{id}", handlers.GetOrder).Methods("GET")
	// Add more routes as necessary
}
