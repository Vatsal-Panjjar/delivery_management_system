package web

import (
	"delivery_management_system/internal/handlers"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/admin", handlers.AdminHandler)
	http.HandleFunc("/order", handlers.OrderHandler)
	http.HandleFunc("/tracking", handlers.TrackingHandler)
}
