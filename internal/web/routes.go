package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

func RegisterRoutes(r *chi.Mux, store *db.Store, rdb *redis.Client) {
	userHandler := handlers.NewUserHandler(store)
	orderHandler := handlers.NewOrderHandler(store)
	adminHandler := handlers.NewAdminHandler(store)

	// Public
	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)

	// Protected
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.AuthMiddleware)
		pr.Post("/order", orderHandler.CreateOrder)
	})

	// Admin
	r.Group(func(ar chi.Router) {
		ar.Use(middleware.AuthMiddleware)
		ar.Get("/admin/orders", adminHandler.GetAllOrders)
	})
}
