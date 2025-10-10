package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

func main() {
	// Connect to PostgreSQL
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=YOUR_PASSWORD dbname=delivery sslmode=disable")
	if err != nil {
		log.Fatal("PostgreSQL connection failed:", err)
	}

	// Initialize repository and cache
	deliveryRepo := repo.NewDeliveryRepo(db)
	redisCache := cache.NewRedisCache("localhost:6379")

	// Initialize handlers
	deliveryHandler := handlers.NewDeliveryHandler(deliveryRepo, redisCache)

	// Initialize router
	router := chi.NewRouter()

	// Register API routes
	handlers.RegisterRoutes(router, deliveryHandler)

	// Serve frontend files from web folder
	fs := http.FileServer(http.Dir("./web"))
	router.Handle("/*", fs)

	// Start server
	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
