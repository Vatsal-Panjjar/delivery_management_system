package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/workers"
)

func main() {
	fmt.Println("🚀 Starting Delivery Management System Server...")

	// DB connection (hardcoded)
	pgURL := "postgres://postgres:rupupuru@01@localhost:5432/delivery_db?sslmode=disable"
	db, err := sqlx.Connect("postgres", pgURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect Postgres: %v", err)
	}
	defer db.Close()
	fmt.Println("✅ Connected to Postgres")

	// Redis
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect Redis: %v", err)
	}
	fmt.Println("✅ Connected to Redis")

	// Initialize router
	r := chi.NewRouter()

	// Handler setup
	h := handlers.NewHandler(db, rdb)

	// Auth routes
	r.Get("/register", h.ShowRegister)
	r.Post("/register", h.HandleRegister)
	r.Get("/login", h.ShowLogin)
	r.Post("/login", h.HandleLogin)
	r.Post("/logout", h.HandleLogout)

	// Protected user routes
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Get("/dashboard", h.Dashboard)
		r.Post("/orders", h.CreateOrder)
		r.Post("/orders/{id}/cancel", h.CancelOrder)
	})

	// Admin routes
	r.Group(func(r chi.Router) {
		r.Use(h.AdminMiddleware)
		r.Get("/admin", h.AdminDashboard)
		r.Post("/admin/orders/{id}/status", h.AdminUpdateStatus)
	})

	// Worker for async order updates
	tracker := workers.NewOrderTracker(db, rdb)
	go tracker.Run()

	fmt.Println("🌐 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
